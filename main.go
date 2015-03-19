package main

import (
	"flag"
	"github.com/mauidude/go-chat/client"
	"github.com/mauidude/go-chat/server"
	"log"
	"os"
)

func main() {
	var runAsClient bool
	var addr string
	flag.BoolVar(&runAsClient, "client", false, "run as a client")
	flag.StringVar(&addr, "addr", ":9090", "the address of the server")
	flag.Parse()

	if runAsClient {

		c, err := client.New(addr, os.Stdin, os.Stdout)
		if err != nil {
			log.Fatalf("unable to connect to server at %s", addr)
		}

		log.Fatal(c.Run())
	} else {

		s := &server.Server{}

		log.Println("starting server")
		log.Fatal(s.ListenAndServe(":9090"))
	}
}
