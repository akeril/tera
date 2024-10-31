package main

type Config struct {
	Port     string
	WatchDir string
}

func NewConfig() Config {
	return Config{
		Port:     "8080",
		WatchDir: ".",
	}
}
