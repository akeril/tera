package main

import "log"

func main() {

	cfg, err := ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	runner, err := NewRunnerWithConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	runner.Run()
}
