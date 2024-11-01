package main

import "log"

func main() {

	cfg, err := ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	runner := NewRunnerWithConfig(cfg)

	runner.Run()
}
