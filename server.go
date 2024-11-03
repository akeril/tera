package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

type Server struct {
	server  *http.Server
	port    int
	clients map[*websocket.Conn]bool
	mu      *sync.RWMutex
}

func NewServer(port int, watchDir string) Server {
	// define http handlers
	server := Server{
		server: &http.Server{
			Addr: ":" + strconv.Itoa(port),
		},
		port:    port,
		mu:      &sync.RWMutex{},
		clients: make(map[*websocket.Conn]bool),
	}
	http.Handle("GET /", http.FileServer(http.Dir(watchDir)))
	http.HandleFunc("GET /tera", server.handleDefault)
	http.HandleFunc("GET /__internal/ws", server.handleWS)

	return server

}

//go:embed templates/*
var fs embed.FS

// handle default screen for tera
func (s Server) handleDefault(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("templates/templ.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	templ.Execute(w, nil)
}

// handles incoming websocket requests
func (s Server) handleWS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// close socket connection
	defer func() {
		log.Printf("Closing websocket connection from %v\n", conn.RemoteAddr())
		conn.Close()
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
	}()

	log.Printf("New websocket connection from %v\n", conn.RemoteAddr())
	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	// return if connection is terminated
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func (s Server) Serve() {
	log.Println("Listening on port", s.port)
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (s Server) BroadcastEvents(ch chan fsnotify.Event) {
	for event := range ch {
		log.Println(event)
		data, _ := json.Marshal(event)

		s.mu.Lock()
		for conn := range s.clients {
			log.Printf("Broadcasting event to %v\n", conn.RemoteAddr())
			conn.WriteMessage(websocket.TextMessage, data)
		}
		s.mu.Unlock()
	}
}
