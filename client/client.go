package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"path"

	"net/http"
	"os"

	"http3test/util"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

var ErrorClientVersion = fmt.Errorf("invalid client version")
var clientArray = []func(){client0, client1, client2}
var addr = "localhost:8443"
var File *os.File

func RunClient(clientVersion int) error {
	if clientVersion >= len(clientArray) {
		return ErrorClientVersion
	}
	if clientVersion > -1 {
		var err error
		// if you hard-code this file name, and its your home dir,
		// your code leaks information about machine configureation.
		// so, we use os.UserHomeDir
		fn, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		fn = path.Join(fn, ".ssl-key.log")

		//fn := "ssl-key.log"
		File, err = os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)

		if err != nil {
			panic(err)
		}
		log.Printf("Log file %s created\n", fn)
		defer File.Close()
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
			RootCAs:      getRootCA(util.GetCertFilePath("certs/cert.pem")),
			KeyLogWriter: File,
			//InsecureSkipVerify: true,
		},
	}
	//defer roundTripper.Close()

	client := &http.Client{
		Transport: roundTripper,
	}
	_ = currentPath

	client0Get(client, fmt.Sprintf("https://%s", addr))
	client0Get(client, fmt.Sprintf("https://%s/endpoint-one", addr))
}

func client0Get(client *http.Client, endpoint string) {
	rsp, err := client.Get(endpoint)
	if err != nil {
		panic(err)
	}

	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		panic(err)
	}
	rsp.Body.Close()

	log.Println("Endpoint: ", endpoint)
	log.Printf("Body length: %d bytes \n", body.Len())
	log.Printf("Response body %s \n", body.Bytes())
}

// client 1 uses a https client, but with an http3 transport
func client1() {
	tlsConf := &tls.Config{
		RootCAs:      getRootCA(util.GetCertFilePath("certs/cert.pem")),
		NextProtos:   []string{"quic-echo-example"},
		KeyLogWriter: File,
	}

	roundTripper := &http3.RoundTripper{
		TLSClientConfig: tlsConf,
		QuicConfig:      &quic.Config{},
	}
	client := &http.Client{
		Transport: roundTripper,
	}
	endpoint := fmt.Sprintf("https://%s", addr)
	data := clientOneGet(client, endpoint)
	fmt.Printf("Client: %s Got '%s'\n", endpoint, data)

	// next!
	endpoint = fmt.Sprintf("https://%s/endpoint-one", addr)
	data = clientOneGet(client, endpoint)

	fmt.Printf("Client: %s Got '%s'\n", endpoint, data)

}

func clientOneGet(client *http.Client, endpoint string) string {
	resp, err := client.Get(endpoint)
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	fmt.Printf("Client: Got '%s'\n", data)
	return string(data)

}

func client2() {
	var message = "hello"
	tlsConf := &tls.Config{
		RootCAs:      getRootCA(util.GetCertFilePath("certs/cert.pem")),
		KeyLogWriter: File,
		//InsecureSkipVerify: true,
		NextProtos: []string{"quic-echo-example"},
	}
	conn, err := quic.DialAddr(context.Background(), "localhost:8443", tlsConf, nil)
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
	caCertRaw, err := os.ReadFile(certPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(caCertRaw)

	if !ok {
		panic("failed to parse root certificate")
	}

	return certPool
}
