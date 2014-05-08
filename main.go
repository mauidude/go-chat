package main

import (
	"client"
	"flag"
	"log"
	"os"
	"server"
)

func main() {
	var runAsClient bool
	flag.BoolVar(&runAsClient, "client", false, "run as a client")
	flag.Parse()

	if runAsClient {
		var addr string
		flag.StringVar(&addr, "addr", ":9090", "the address of the server")

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
