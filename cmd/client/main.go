package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

const (
	addr string = "https://localhost:8080"
)

func main() {

	c := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	rs, err := c.Get(fmt.Sprintf("%s%s", addr, "/foo"))
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		panic(err)
	}

	log.Printf("\ncode=%d\nbody=%s\nheaders=%v", rs.StatusCode, body, rs.Header)
}
