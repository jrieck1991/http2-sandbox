package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/http2"
)

const (
	addr string = "https://localhost:8080/stream"
)

func main() {

	// init client with http2, allow self signed certs
	c := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// form request with headers to enable streaming
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "text/event-stream")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Access-Control-Allow-Origin", "*")

	// send request
	rs, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	log.Printf("connected: %+v", rs)

	// copy data streamed to std out
	n, err := io.Copy(os.Stdout, rs.Body)
	if err != nil {
		log.Printf("%d bytes written", n)
		if err == io.EOF {
			// no more data to read, exit without error
			log.Println(err)
		} else {
			panic(err)
		}
	}
}
