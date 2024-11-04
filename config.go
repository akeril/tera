package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port       int
	WatchDir   string
	Entrypoint string
	Exts       []string
}

func ParseConfig() (Config, error) {
	port := flag.Int("port", 5199, "Specify the port number")
	watch := flag.String("watch", ".", "Specify the directory to be watched")
	exts := flag.String("exts", "", "Specify file extensions to be filtered to")

	flag.Usage = func() {
		fmt.Println("Render file changes in the browser")
		fmt.Println("\nUsage:")
		fmt.Println("  tera [OPTS] [FILE]")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	entrypoint := flag.Arg(0)

	cfg := Config{
		Port:       *port,
		WatchDir:   *watch,
		Entrypoint: entrypoint,
		Exts:       parseExts(*exts),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func parseExts(exts string) []string {
	return strings.Split(exts, ",")
}

func (c Config) Validate() error {
	if c.Port < 1024 {
		return errors.New("Configuration Error: Invalid port")
	}
	if _, err := os.Stat(c.WatchDir); err != nil {
		return errors.New("Configuration error: Directory not found")
	}
	if c.Entrypoint == "" {
		return errors.New("Configuration error: Entrypoint positional field required!")
	}
	if _, err := os.Stat(c.Entrypoint); err != nil {
		return errors.New("Configuration error: Entrypoint not found")
	}
	return nil
}
