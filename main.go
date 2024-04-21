package main

import (
	"flag"
	"http3test/client"
	"http3test/server"
	"log"
)

// Main starts the server
func main() {
	var bClient bool
	var curServer int

	flag.BoolVar(&bClient, "c", false, "run as a client")
	flag.IntVar(&curServer, "s", 0, "server type")
	flag.Parse()

	if bClient {
		client.Client()
	} else {

		err := server.RunServer(curServer)

		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		select {}
	}
}
