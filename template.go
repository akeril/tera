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
	Script     []byte
}

func generateTemplate(cfg TemplConfig) (io.ReadWriter, error) {
	wr := new(bytes.Buffer)

	templ, err := template.ParseFiles("templates/index.templ.js")
	if err != nil {
		return nil, err
	}

	err = templ.Execute(wr, cfg)
	if err != nil {
		return nil, err
	}
	return wr, nil
}
