package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	addr string = "localhost:8080"
	crt  string = "server/server.crt"
	key  string = "server/server.key"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/stream", stream)
	s := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// listen for client connections
	if err := http.ListenAndServeTLS(addr, crt, key, s.Handler); err != nil {
		panic(err)
	}
}

func stream(w http.ResponseWriter, r *http.Request) {

	// check client can support streaming
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Println("flush not supported by client")
		w.Write([]byte("flush not supported by client"))
	}

	// send data
	for i := 0; i < 100; i++ {

		time.Sleep(1 * time.Second)

		fmt.Fprintf(w, "%d\n", i)
		flusher.Flush()
	}

	return
}
