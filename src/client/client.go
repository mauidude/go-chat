package client

import (
	"bufio"
	"io"
	"net"
)

type Client struct {
	conn   net.Conn
	input  *bufio.Reader
	output io.Writer
	reader *bufio.Reader
}

func (c *Client) send(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Client) receive(b []byte) error {
	_, err := c.output.Write(b)
	return err
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Run() error {
	sendChannel := make(chan []byte)
	receiveChannel := make(chan []byte)
	errChannel := make(chan error)

	// listen for input from user
	go func() {
		for {
			line, _, err := c.input.ReadLine()
			if err != nil {
				errChannel <- err
			}
			sendChannel <- line
		}
	}()

	// listen for server messages
	go func() {
		for {
			line, _, err := c.reader.ReadLine()
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
