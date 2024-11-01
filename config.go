package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Port     int
	WatchDir string
}

func ParseConfig() (Config, error) {
	port := flag.Int("port", 5199, "Specify the port number")
	watch := flag.String("watch", ".", "Specify the directory to be watched")

	flag.Usage = func() {
		fmt.Println("Render file changes in the browser")
		fmt.Println("\nUsage:")
		fmt.Println("  tera [OPTS] [FILE]")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	cfg := Config{
		Port:     *port,
		WatchDir: *watch,
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if c.Port < 1024 {
		return errors.New("Configuration Error: Invalid port")
	}
	if _, err := os.Stat(c.WatchDir); err != nil {
		return errors.New("Configuration error: Directory not found")
	}
	return nil
}
