package main

import (
	"log"
	"net/http"
	"strconv"
)

type Runner struct {
	cfg     Config
	watcher Watcher
	server  *http.Server
}

func NewRunnerWithConfig(cfg Config) (*Runner, error) {
	server := http.Server{
		Addr: ":" + strconv.Itoa(cfg.Port),
	}
	watcher, err := NewWatcher()
	if err != nil {
		return nil, err
	}
	watcher.Add(cfg.WatchDir)

	return &Runner{
		cfg:     cfg,
		watcher: watcher,
		server:  &server,
	}, nil
}

func (r *Runner) Run() {
	log.Println("Listening on port", r.cfg.Port)
	http.Handle("GET /", http.FileServer(http.Dir(r.cfg.WatchDir)))

	go r.watcher.Watch()
	if err := r.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
