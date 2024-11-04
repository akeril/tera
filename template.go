package main

import (
	"bytes"
	"embed"
	"html/template"
	"io"
)

//go:embed templates/*
var fs embed.FS

type TemplConfig struct {
	Uri        string
	Entrypoint string
}

func generateTemplate(cfg TemplConfig) (io.ReadWriter, error) {
	templ, err := template.ParseFS(fs, "templates/templ.html")
	if err != nil {
		return nil, err
	}

	wr := new(bytes.Buffer)
	err = templ.Execute(wr, cfg)
	if err != nil {
		return nil, err
	}
	return wr, nil
}
