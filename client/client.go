package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"

	"net/http"
	"os"

	"http3test/util"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

var ErrorClientVersion = fmt.Errorf("invalid client version")
var clientArray = []func(){client0, client1}
var addr = "localhost:8443"

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
			RootCAs: getRootCA(util.GetCertFilePath("certs/cert.pem")),
			//InsecureSkipVerify: true,
		},
	}
	//defer roundTripper.Close()

	client := &http.Client{
		Transport: roundTripper,
	}
	_ = currentPath

	rsp, err := client.Get(fmt.Sprintf("https://%s", addr))
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
	tlsConf := &tls.Config{
		RootCAs:    getRootCA(util.GetCertFilePath("certs/cert.pem")),
		NextProtos: []string{"quic-echo-example"},
	}

	roundTripper := &http3.RoundTripper{
		TLSClientConfig: tlsConf,
		QuicConfig:      &quic.Config{},
	}
	client := &http.Client{
		Transport: roundTripper,
	}
	resp, err := client.Get(fmt.Sprintf("https://%s", addr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Client: Got '%s'\n", data)

}

/*func client1() {
	var message = "hello"
	tlsConf := &tls.Config{
		RootCAs: getRootCA(util.GetCertFilePath("certs/cert.pem")),
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
}*/

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
