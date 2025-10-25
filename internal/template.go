package internal

import (
	"html/template"
	"io/fs"
	"net/http"
)

type Template struct {
	templates *template.Template
}

// NewTemplate takes a html/template instance with parsed files.
func NewTemplate(templates *template.Template) *Template {
	return &Template{templates}
}

// NewTemplateFS parses templates from any FS and returns a Template instance.
func NewTemplateFS(fs fs.FS, pattern string) (*Template, error) {
	tmpl, err := template.ParseFS(fs, pattern)
	if err != nil {
		return nil, err
	}

	return &Template{
		templates: tmpl,
	}, nil
}

func (t *Template) Render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := t.templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
