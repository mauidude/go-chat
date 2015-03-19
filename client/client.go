package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type Client struct {
	conn   net.Conn
	input  *bufio.Reader
	output io.Writer
	reader *bufio.Reader
}

// Send to Server
func (c *Client) send(s string) error {
	_, err := io.WriteString(c.conn, s)
	return err
}

// Write to output
func (c *Client) receive(s string) error {
	_, err := io.WriteString(c.output, s)
	return err
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Run() error {
	sendChannel := make(chan string)
	receiveChannel := make(chan string)
	errChannel := make(chan error)

	// listen for input from user
	go func() {
		for {
			line, err := c.input.ReadString('\n')
			if err != nil {
				errChannel <- err
			}
			sendChannel <- line
		}
	}()

	// listen for server messages
	go func() {
		for {
			line, err := c.reader.ReadString('\n')
			if err != nil {
				errChannel <- err
			}
			receiveChannel <- line
		}
	}()

	go func() {
		for {
			select {
			case line := <-sendChannel:
				c.send(line)
			case line := <-receiveChannel:
				c.receive(line)
			}
		}
	}()

	select {
	case err := <-errChannel:
		fmt.Println("Error: ", err)
		return err
	}

	return nil
}

func New(addr string, input io.Reader, output io.Writer) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, input: bufio.NewReader(input), output: output, reader: bufio.NewReader(conn)}, nil
}
