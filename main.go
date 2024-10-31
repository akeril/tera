package main

func main() {

	cfg := NewConfig()

	runner := NewRunnerWithConfig(&cfg)

	runner.Run()
}
