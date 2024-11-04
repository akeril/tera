package main

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	eventCh chan fsnotify.Event
	exts    []string
}

// Create a new Watcher
func NewWatcher(exts []string) (Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return Watcher{}, nil
	}
	ch := make(chan fsnotify.Event)
	return Watcher{watcher: watcher, eventCh: ch, exts: exts}, err
}

// TODO: Add paths recursively
func (w *Watcher) Add(path string) {
	w.watcher.Add(path)
}

// Watcher watches a set of paths, delivering events on a channel.
func (w *Watcher) Watch() {
	log.Println("Watching filesystem for changes...")
	for event := range w.watcher.Events {
		if event.Has(fsnotify.Write) && (len(w.exts) == 0 || hasSuffix(w.exts, event.Name)) {
			w.eventCh <- event
		}
	}
}

func hasSuffix(list []string, s string) bool {
	for _, ext := range list {
		if strings.HasSuffix(s, ext) {
			return true
		}
	}
	return false
}
