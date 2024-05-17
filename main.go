package main

import (
	"flag"
	"http3test/client"
	"http3test/server"
	"log"
	"time"
)

// Main starts the server
func main() {
	var clientVersion, curServer, repeats int

	var err error

	flag.IntVar(&clientVersion, "c", -1, "run as a client")
	flag.IntVar(&curServer, "s", 0, "server type")
	flag.IntVar(&repeats, "n", 1, "number of times to run the client")
	flag.Parse()

	// if the client flag is set to anything other than -1, run the client

	err = server.RunServer(curServer)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	//lazy wait for server to come up
	time.Sleep(1 * time.Second)

	for range repeats {
		err = client.RunClient(clientVersion)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	// block forever
	if clientVersion == -1 {
		select {}
	}

}
