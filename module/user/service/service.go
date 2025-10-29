package service

import (
	"html/template"
	"io/fs"
	"net/http"
	"path"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/module/user/storage"
)

// Service encapsulates what we need to get from the handler.
type Service struct {
	templates map[string]*internal.Template

	UserStorage    *storage.UserStorage
	SessionStorage *storage.SessionStorage
}

// NewService takes in required dependencies to support the MVC framework.
func NewService(options *Options) (*Service, error) {
	svc := &Service{
		UserStorage:    options.UserStorage,
		SessionStorage: options.SessionStorage,
	}

	// Set up the template views.
	if err := svc.initTemplates(options.ThemeFS, options.ModuleFS); err != nil {
		return nil, err
	}
	return svc, nil
}

// initTemplate goes through the moduleFS templates and creates
// a wrapper that replaces the content block from base.tpl (themeFS).
//
// This is an implementation detail of go templates, can't just
// load all the templates as each `{{define}}` is evaluated at
// parsing time, not runtime.
func (h *Service) initTemplates(themeFS, moduleFS fs.FS) error {
	files, err := fs.Glob(moduleFS, "template/*.tpl")
	if err != nil {
		return err
	}

	templates := make(map[string]*internal.Template)

	for _, f := range files {
		file := path.Base(f)
		contents, _ := fs.ReadFile(moduleFS, f)

		tmpl := template.Must(template.ParseFS(themeFS, "template/*.tpl"))
		tmpl = template.Must(tmpl.New(file).Parse(string(contents)))
		tmpl = template.Must(tmpl.New("wrapper").Parse(`
			{{define "content"}}{{template "` + file + `" .}}{{end}}
			{{template "base.tpl" .}}
		`))

		templates[file] = internal.NewTemplate(tmpl)
	}

	h.templates = templates

	return nil
}

// View is a helper to add modularity to templates. It renders a view with the theme base.tpl.
// The intent is to override the "content" block in the base.tpl with a view.
func (h *Service) View(w http.ResponseWriter, name string, data any) {
	tmpl, ok := h.templates[name]
	if ok {
		tmpl.Render(w, "wrapper", data)
	}
}

// Close just clears the loaded template map. After that, the *Service should not be used.
func (h *Service) Close() {
	h.templates = nil
}
