package server

import (
	"fmt"
	"log"
	"net/http"
)

// A Broker holds open client connections,
// listens for incoming events on its Notifier channel
// and broadcast event data to all registered connections
type Broker struct {
	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan []byte
	// New client connections
	newClients chan chan []byte
	// Closed client connections
	closingClients chan chan []byte
	// Client connections registry
	clients map[chan []byte]bool
}

func New() *Broker {

	b := &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	go b.listen()

	return b
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// check client supports flusher
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Printf("server flush not supported for request: %+v", r)
		w.Write([]byte("server flush not supported\n"))
	}

	// set headers to support keepalive http connections
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// each connection registers its own message channel with brokers conn registry
	messageChan := make(chan []byte)

	// signal the broker that we have a new connection
	b.newClients <- messageChan

	// when handler exits, notify broker
	defer func() {
		b.closingClients <- messageChan
	}()

	// on connection close, un register message channel
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		b.closingClients <- messageChan
	}()

	// block waiting for messages
	for {

		// server sent events compatible
		fmt.Fprintf(w, "%s\n\n", <-messageChan)

		// flush data immediately instead of buffering
		flusher.Flush()
	}
}

func (b *Broker) listen() {

	for {
		select {

		case s := <-b.newClients:
			// a new client has connected, register their message channel
			b.clients[s] = true
			log.Printf("Client added. %d registered clients", len(b.clients))

		case s := <-b.closingClients:
			// a client is dettached and we want
			// to stop sending them messages
			delete(b.clients, s)
			log.Printf("Removed client. %d registered clients", len(b.clients))

		case event := <-b.Notifier:

			// new event, send it to all connected clients
			for c := range b.clients {
				c <- event
			}
		}

	}
}
