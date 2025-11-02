package internal

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

type Template struct {
	templates *template.Template
}

// NewTemplate takes a html/template instance with parsed files.
func NewTemplate(templates *template.Template) *Template {
	return &Template{templates}
}

// Render a template to a ResponseWriter. If an error is returned, no
// data has been written to the response writer and the error needs handling.
func (t *Template) Render(w http.ResponseWriter, name string, data any) error {
	buf := &bytes.Buffer{}
	if err := t.templates.ExecuteTemplate(buf, name, data); err != nil {
		return fmt.Errorf("template error: %w", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, err := io.Copy(w, buf)
	return err
}
