package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"

	"net/http"
	"os"
	"path"

	"github.com/quic-go/quic-go"
)

var ErrorClientVersion = fmt.Errorf("invalid client version")
var clientArray = []func(){client0, client1}
var addr = "https://localhost:443"

func RunClient(clientVersion int) error {
	if clientVersion >= len(clientArray) {
		return ErrorClientVersion
	}
	if clientVersion > -1 {
		clientArray[clientVersion]()
		return nil
	}
	return nil
}

// Client 0 is a basic HTTPS client. It connects to a server using HTTPS.
func client0() {
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	roundTripper := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: getRootCA(currentPath),
			//InsecureSkipVerify: true,
		},
	}
	//defer roundTripper.Close()

	client := &http.Client{
		Transport: roundTripper,
	}
	_ = currentPath
	addr := "https://localhost:443"
	rsp, err := client.Get(addr)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()

	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		panic(err)
	}

	log.Printf("Body length: %d bytes \n", body.Len())
	log.Printf("Response body %s \n", body.Bytes())
}
func client1() {
	var message = "hello"
	addr := "localhost:8443"
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	tlsConf := &tls.Config{
		RootCAs: getRootCA(currentPath),
		//InsecureSkipVerify: true,
		NextProtos: []string{"quic-echo-example"},
	}
	conn, err := quic.DialAddr(context.Background(), addr, tlsConf, nil)
	if err != nil {
		panic(err)
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Client: Sending '%s'\n", message)
	_, err = stream.Write([]byte(message))
	if err != nil {
		panic(err)
	}

	buf := make([]byte, len(message))
	_, err = io.ReadFull(stream, buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Client: Got '%s'\n", buf)
}
func getRootCA(certPath string) *x509.CertPool {
	caCertPath := path.Join(certPath, "certs/cert.pem")
	caCertRaw, err := os.ReadFile(caCertPath)
	if err != nil {
		panic(err)
	}

	//p, _ := pem.Decode(caCertRaw)
	/*if p.Type != "CERTIFICATE" {
		panic("expected a certificate")
	}*/

	/*	caCert, err := x509.ParseCertificate(p.Bytes)
		if err != nil {
			panic(err)
		} */

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(caCertRaw)

	if !ok {
		panic("failed to parse root certificate")
	}

	return certPool
}
