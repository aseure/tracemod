package main

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"
)

//go:embed index.html.tmpl
var htmlTemplate string

func generateHTML(diagram string) (string, error) {
	tmpl, err := template.New("index.html.tmpl").Parse(htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("could not parse template: %w", err)
	}

	outputHtml, err := os.CreateTemp("", "tracemod_*.html")
	if err != nil {
		return "", fmt.Errorf("could not create HTML output file: %w", err)
	}
	defer outputHtml.Close()

	err = tmpl.Execute(outputHtml, diagram)
	if err != nil {
		return "", fmt.Errorf("could not generate HTML output file: %w", err)
	}

	return outputHtml.Name(), nil
}
