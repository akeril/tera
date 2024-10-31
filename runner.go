package main

import (
	"log"
	"net/http"
)

type Runner struct {
	cfg    *Config
	server *http.Server
}

func NewRunner() *Runner {
	cfg := NewConfig()
	return NewRunnerWithConfig(&cfg)
}

func NewRunnerWithConfig(cfg *Config) *Runner {
	server := http.Server{
		Addr: ":" + cfg.Port,
	}
	return &Runner{
		cfg:    cfg,
		server: &server,
	}
}

func (r *Runner) Run() {
	http.Handle("GET /", http.FileServer(http.Dir(r.cfg.WatchDir)))

	if err := r.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
