package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	eventCh chan fsnotify.Event
}

// Create a new Watcher
func NewWatcher() (Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return Watcher{}, nil
	}
	ch := make(chan fsnotify.Event)
	return Watcher{watcher: watcher, eventCh: ch}, err
}

// TODO: Add paths recursively
func (w *Watcher) Add(path string) {
	w.watcher.Add(path)
}

// Watcher watches a set of paths, delivering events on a channel.
func (w *Watcher) Watch() {
	log.Println("Watching filesystem for changes...")
	for event := range w.watcher.Events {
		if event.Has(fsnotify.Write) {
			w.eventCh <- event
		}
	}
}
