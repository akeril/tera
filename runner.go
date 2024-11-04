package main

type Runner struct {
	cfg     Config
	watcher Watcher
	server  Server
}

func NewRunnerWithConfig(cfg Config) (*Runner, error) {
	watcher, err := NewWatcher(cfg.Exts)
	if err != nil {
		return nil, err
	}
	watcher.Add(cfg.WatchDir)

	server := NewServer(cfg.Port, cfg.WatchDir, cfg.Entrypoint)
	return &Runner{
		cfg:     cfg,
		watcher: watcher,
		server:  server,
	}, nil
}

func (r *Runner) Run() {
	// watch fs for changes
	go r.watcher.Watch()

	go r.server.BroadcastEvents(r.watcher.eventCh)

	// serve http requests
	r.server.Serve()
}
