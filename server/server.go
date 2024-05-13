package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"http3test/util"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

const certPath = "certs/cert.pem"
const keyPath = "certs/private.key"

var addr = `0.0.0.0:8443`

var ErrInvalidServerType = fmt.Errorf("invalid server type")

func RunServer(serverType int) error {
	if serverType < 0 || serverType >= len(Servers) {
		return ErrInvalidServerType
	}
	go Servers[serverType]()
	return nil
}

var Servers = []func(){server0, server1, server2, server3}

// Server 0 is a basic HTTPS server. It listens on TCP and serves HTTPS.
// Server 0 does not use QUIC.
func server0() {
	serverName := "server 0"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v from %s", r.URL.Path, serverName)
	})
	/*if len(os.Args) > 1 {
		addr = os.Args[1]
	}*/
	// get the current working director

	log.Println("server 0 listens and servers HTTPS")
	log.Fatal(http.ListenAndServeTLS(addr, util.GetCertFilePath(certPath), util.GetCertFilePath(keyPath), nil))
}

// Server 1 is the most polite QUIC server. It listens on TCP and politely
// informs the client that it is using HTTP/3. Most modern browsers will handle this
// server well.
func server1() {
	serverName := "server 1"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v from %s", r.URL.Path, serverName)
	})

	/*if len(os.Args) > 1 {
		addr = os.Args[1]
	}*/
	log.Println("server 1 listens and servers HTTP/3")
	log.Fatal(http3.ListenAndServe(addr, util.GetCertFilePath(certPath), util.GetCertFilePath(keyPath), nil))

}

// Server 2 uses the listenand serve QUIC function. Without the additional
// header provided at the TCP layer, some browsers will struggle with this
// server.
func server2() {
	serverName := "server 2"
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//r.Header.Add("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"message": "Hello, %v from %s"}`, r.URL.Path, serverName)))
	})

	// ... add HTTP handlers to mux ...
	// If mux is nil, the http.DefaultServeMux is used.

	log.Println("Listening on", addr)
	log.Println("server 2 listens and servers QUIC")
	log.Fatal(http3.ListenAndServeQUIC(addr, util.GetCertFilePath(certPath), util.GetCertFilePath(keyPath), mux))
}

// server3 is a QUIC server that echos all data on the first stream opened by the client
// server 3 is NOT an HTTP server, so HTTP clients will have a hard time with it.
func server3() {
	addr := "localhost:8443"
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	listener, err := quic.ListenAddr(addr, &tls.Config{
		Certificates: []tls.Certificate{getCertificate(currDir)},
		NextProtos:   []string{"quic-echo-example"},
	}, nil)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer listener.Close()

	conn, err := listener.Accept(context.Background())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer stream.Close()

	// Echo through the loggingWriter
	_, err = io.Copy(loggingWriter{stream}, stream)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Println("Closed")
}

type loggingWriter struct {
	io.Writer
}

func (w loggingWriter) Write(p []byte) (n int, err error) {
	log.Printf("Got message: %s", p)
	return w.Writer.Write(p)
}
func getCertificate(certPath string) tls.Certificate {
	caCertPath := path.Join(certPath, "certs/cert.pem")

	cert, err := tls.LoadX509KeyPair(caCertPath, path.Join(certPath, "certs/private.key"))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return cert
}
