package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
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

	templ, err := template.ParseFS(fs, "templates/index.templ.js")
	if err != nil {
		return nil, err
	}

	err = templ.Execute(wr, cfg)
	if err != nil {
		return nil, err
	}
	return wr, nil
}

func generateIndex(cfg TemplConfig) ([]byte, error) {
	if strings.HasSuffix(cfg.Entrypoint, ".html") {
		html, err := os.ReadFile(cfg.Entrypoint)
		if err != nil {
			return nil, err
		}
		html = injectLink(html)
		return html, nil
	}
	return nil, nil
}

func injectLink(html []byte) []byte {
	script := "<script src='/tera'></script>"
	inject := fmt.Sprintf("<head>%v</head>", script)
	return []byte(inject + string(html))
}
