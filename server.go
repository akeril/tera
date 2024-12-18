package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

type Server struct {
	server     *http.Server
	port       int
	clients    map[*websocket.Conn]bool
	mu         *sync.RWMutex
	entrypoint string
	watchDir   string
}

func NewServer(port int, watchDir string, entrypoint string) Server {
	// define http handlers
	s := Server{
		server: &http.Server{
			Addr: ":" + strconv.Itoa(port),
		},
		port:       port,
		mu:         &sync.RWMutex{},
		clients:    make(map[*websocket.Conn]bool),
		entrypoint: entrypoint,
		watchDir:   watchDir,
	}
	http.HandleFunc("GET /", s.Router)
	return s
}

func (s Server) Serve() {
	log.Printf("Listening on http://localhost:%d", s.port)
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// multiplex on websocket and http requests
func (s Server) Router(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "websocket" {
		s.handleWS(w, r)
		return
	}
	if r.URL.Path == "/" {
		s.handleIndex(w, r)
		return
	}
	if r.URL.Path == "/tera" {
		s.handleReloader(w, r)
		return
	}
	s.handleFS(w, r)
}

// returns index page for tera with scripts injected
func (s Server) handleIndex(w http.ResponseWriter, _ *http.Request) {
	data, err := ParseEntryPoint(TemplConfig{
		Uri:        fmt.Sprintf("ws://localhost:%v", s.port),
		Entrypoint: s.entrypoint,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// returns javascript file for client side live reloading
func (s Server) handleReloader(w http.ResponseWriter, _ *http.Request) {
	data, err := ParseTemplate("templates/tera.js", TemplConfig{
		Uri:        fmt.Sprintf("ws://localhost:%v", s.port),
		Entrypoint: s.entrypoint,
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/javascript; charset=utf-8")
	w.Write(data)
}

// file server handler
func (s Server) handleFS(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir(s.watchDir)).ServeHTTP(w, r)
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

// invoke on separate thread to broadcast events
func (s Server) BroadcastEvents(ch chan fsnotify.Event) {
	for event := range ch {
		log.Println(event)
		data, _ := json.Marshal(event)
		s.BroadcastEvent(ch, data)
	}
}

// separate function because defer is function scoped
func (s Server) BroadcastEvent(ch chan fsnotify.Event, data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for conn := range s.clients {
		log.Printf("Broadcasting event to %v\n", conn.RemoteAddr())
		conn.WriteMessage(websocket.TextMessage, data)
	}
}
