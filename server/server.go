package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/quic-go/quic-go/http3"
)

const cert = "./certs/cert.pem"
const key = "./certs/private.key"

var ErrInvalidServerType = fmt.Errorf("invalid server type")

func RunServer(serverType int) error {
	if serverType < 0 || serverType >= len(Servers) {
		return ErrInvalidServerType
	}
	go Servers[serverType]()
	return nil
}

var Servers = []func(){server0, server1, server2}

func server0() {
	serverName := "server 0"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v from %s", r.URL.Path, serverName)
	})
	addr := `0.0.0.0:443`
	/*if len(os.Args) > 1 {
		addr = os.Args[1]
	}*/
	log.Println("server 0 listens and servers HTTPS")
	log.Fatal(http.ListenAndServeTLS(addr, cert, key, nil))
}

func server1() {
	serverName := "server 1"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v from %s", r.URL.Path, serverName)
	})
	addr := `0.0.0.0:443`
	/*if len(os.Args) > 1 {
		addr = os.Args[1]
	}*/
	log.Println("server 1 listens and servers HTTP/3")
	log.Fatal(http3.ListenAndServe(addr, cert, key, nil))

}

func server2() {
	serverName := "server 2"
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//r.Header.Add("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"message": "Hello, %v from %s"}`, r.URL.Path, serverName)))
	})

	// ... add HTTP handlers to mux ...
	// If mux is nil, the http.DefaultServeMux is used.
	addr := `0.0.0.0:443`
	log.Println("Listening on", addr)
	log.Println("server 2 listens and servers QUIC")
	log.Fatal(http3.ListenAndServeQUIC(`0.0.0.0:443`, cert, key, mux))
}
