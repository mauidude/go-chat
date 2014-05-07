package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	name string
	conn net.Conn
}

type Message struct {
	data   []byte
	source *Client
}

type Server struct {
	clients []*Client
}

func (s *Server) Serve(l net.Listener) error {
	s.clients = make([]*Client, 0)
	messageChan := make(chan *Message)
	counter := 0

	go func() {
		for {
			msg := <-messageChan

			for _, c := range s.clients {
				if c != msg.source {
					_, err := fmt.Fprintf(c.conn, "%s: %s", msg.source.name, msg.data)
					if err != nil {
						log.Println("unable to write to %s: %s", c.conn.RemoteAddr(), err.Error())
					}
				}
			}
		}
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn: conn,
			name: fmt.Sprintf("guest%d", counter),
		}

		counter++

		s.clients = append(s.clients, client)
		go s.handle(client, messageChan)
	}
}

func (s *Server) handle(c *Client, messageChan chan *Message) {
	log.Println("new connection from", c.conn.RemoteAddr())

	fmt.Fprintln(c.conn, "Welcome %s", c.name)
	for {
		msg, _ := bufio.NewReader(c.conn).ReadString('\n')
		messageChan <- &Message{
			data:   []byte(msg),
			source: c,
		}
	}
}

func main() {
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal("unable to listen", err)
	}

	s := &Server{}

	log.Println("starting server")
	log.Fatal(s.Serve(l))
}
