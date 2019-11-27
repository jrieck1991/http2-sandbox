package main

import (
	"fmt"
	"net/http"
	"sandbox/http2-sandbox/internal/server"
	"time"
)

const (
	addr string = "localhost:8080"
)

func main() {

	// init server
	b := server.New()

	// send events
	go func() {

		tick := time.Tick(5 * time.Second)
		for t := range tick {

			b.Notifier <- []byte(fmt.Sprintf("event %v", t))
		}

	}()

	// listen for client connections
	if err := http.ListenAndServeTLS(addr, "internal/server/server.crt", "internal/server/server.key", b); err != nil {
		panic(err)
	}

}
