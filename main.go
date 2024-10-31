package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	port := "8080"
	watchDir := "." // current directory

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("GET /", http.FileServer(http.Dir(watchDir)))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
