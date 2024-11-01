package main

import (
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	server *http.Server
	port   int
}

func NewServer(port int, watchDir string) Server {
	http.Handle("GET /", http.FileServer(http.Dir(watchDir)))
	return Server{
		server: &http.Server{
			Addr: ":" + strconv.Itoa(port),
		},
		port: port,
	}
}

func (s Server) Serve() {
	log.Println("Listening on port", s.port)
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
