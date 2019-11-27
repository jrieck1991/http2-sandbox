package main

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"

	"golang.org/x/net/http2"
)

const (
	addr string = "https://localhost:8080"
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

	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "text/event-stream")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Access-Control-Allow-Origin", "*")

	rs, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, rs.Body)
}
