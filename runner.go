package main

import (
	"log"
	"net/http"
	"strconv"
)

type Runner struct {
	cfg    Config
	server *http.Server
}

func NewRunnerWithConfig(cfg Config) *Runner {
	server := http.Server{
		Addr: ":" + strconv.Itoa(cfg.Port),
	}
	return &Runner{
		cfg:    cfg,
		server: &server,
	}
}

func (r *Runner) Run() {
	log.Println("Listening on port", r.cfg.Port)
	http.Handle("GET /", http.FileServer(http.Dir(r.cfg.WatchDir)))

	if err := r.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
