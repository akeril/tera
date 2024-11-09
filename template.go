package main

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"
)

//go:embed templates/*
var templFS embed.FS

type TemplConfig struct {
	Uri        string
	Entrypoint string
}

func ParseTemplate(file string, cfg TemplConfig) ([]byte, error) {
	var buf bytes.Buffer

	templ, err := template.ParseFS(templFS, file)
	if err != nil {
		return nil, err
	}

	err = templ.Execute(&buf, cfg)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ParseEntryPoint(cfg TemplConfig) ([]byte, error) {
	var data []byte
	var err error

	if strings.HasSuffix(cfg.Entrypoint, ".html") {
		data, err = os.ReadFile(cfg.Entrypoint)
	} else if strings.HasSuffix(cfg.Entrypoint, ".pdf") {
		data, err = ParseTemplate("templates/pdf.templ.html", cfg)
	} else {
		data, err = nil, errors.New("Unsupported file type!")
	}
	if err != nil {
		return nil, err
	}
	html := injectScript(data, "<script src='/tera'></script>")
	return html, nil
}

func injectScript(html []byte, script string) []byte {
	inject := fmt.Sprintf("<head>%v</head>", script)
	return []byte(inject + string(html))
}
