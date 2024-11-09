package main

import (
	"crypto/md5"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	eventCh chan fsnotify.Event
	exts    []string
	cache   *FSCache
}

// Create a new Watcher
func NewWatcher(exts []string) (Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return Watcher{}, nil
	}
	ch := make(chan fsnotify.Event)
	return Watcher{
		watcher: watcher,
		eventCh: ch,
		exts:    exts,
		cache:   NewFSCache(),
	}, err
}

// TODO: Add paths recursively
func (w *Watcher) Add(path string) {
	w.watcher.Add(path)
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if d.IsDir() {
			w.watcher.Add(path)
		}
		return nil
	})
}

// Watcher watches a set of paths, delivering events on a channel.
func (w *Watcher) Watch() {
	log.Println("Watching filesystem for changes...")
	for event := range w.watcher.Events {
		if event.Has(fsnotify.Write) && (len(w.exts) == 0 || hasSuffix(w.exts, event.Name)) && w.cache.HasChanged(event.Name) {
			w.eventCh <- event
		}
		if f, err := os.Stat(event.Name); event.Has(fsnotify.Create) && err == nil && f.IsDir() {
			w.watcher.Add(event.Name)
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

// FS Cache for detecting file changes
type FSCache struct {
	hashMap map[string][16]byte
}

func NewFSCache() *FSCache {
	return &FSCache{
		hashMap: make(map[string][16]byte),
	}
}

func (c *FSCache) ComputeHash(path string) ([16]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return [16]byte{}, err
	}
	return md5.Sum(data), nil
}

func (c *FSCache) HasChanged(path string) bool {
	h, err := c.ComputeHash(path)
	if err != nil {
		return true
	}
	if _, ok := c.hashMap[path]; !ok || c.hashMap[path] != h {
		c.hashMap[path] = h
		return true
	}
	return false
}
