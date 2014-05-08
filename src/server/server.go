package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	name string
	conn net.Conn
	id   int
}

type Message struct {
	data   []byte
	source *Client
}

type IdClients map[int]*Client

type Server struct {
	clients IdClients
}

func (s *Server) ListenAndServe(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return s.Serve(l)
}

func (s *Server) Serve(l net.Listener) error {
	s.clients = make(IdClients, 0)
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
			id:   counter,
		}

		counter++

		// s.clients = append(s.clients, client)
		s.clients[client.id] = client
		go s.handle(client, messageChan)
	}
}

func (s *Server) handle(c *Client, messageChan chan *Message) {
	log.Println("new connection from", c.conn.RemoteAddr())

	fmt.Fprintf(c.conn, "Welcome %s\n", c.name)
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			fmt.Println("Err", err)
			delete(s.clients, c.id)
			break
		}
		messageChan <- &Message{
			data:   []byte(msg),
			source: c,
		}
	}
}
