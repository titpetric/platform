package user

import (
	"html/template"
	"io/fs"
	"net/http"
	"path"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/module/user/storage"
)

type Handler struct {
	UserStorage    *storage.UserStorage
	SessionStorage *storage.SessionStorage

	Templates map[string]*internal.Template
}

func NewHandler(themeFS fs.FS) (*Handler, error) {
	files, err := fs.Glob(TemplateFS, "templates/*.tpl")
	if err != nil {
		return nil, err
	}

	templates := make(map[string]*internal.Template)

	for _, f := range files {
		file := path.Base(f)
		contents, _ := fs.ReadFile(TemplateFS, f)

		tmpl := template.Must(template.ParseFS(themeFS, "templates/*.tpl"))
		tmpl = template.Must(tmpl.New(file).Parse(string(contents)))
		tmpl = template.Must(tmpl.New("wrapper").Parse(`
			{{define "content"}}{{template "` + file + `" .}}{{end}}
			{{template "base.tpl" .}}
		`))

		templates[file] = internal.NewTemplate(tmpl)
	}

	db, err := platform.Database.Connect()
	if err != nil {
		return nil, err
	}

	userStorage := storage.NewUserStorage(db)
	sessionStorage := storage.NewSessionStorage(db)

	return &Handler{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
		Templates:      templates,
	}, nil
}

// Name returns the name of the containing package.
func (h *Handler) Name() string {
	return "user"
}

// View is a helper to add modularity to templates. It renders a view with the theme base.tpl.
// The intent is to override the "content" block in the base.tpl with a view.
func (h *Handler) View(w http.ResponseWriter, name string, data any) {
	tmpl, ok := h.Templates[name]
	if ok {
		tmpl.Render(w, "wrapper", data)
	}
}

// Mount registers login, logout, and register routes.
func (h *Handler) Mount(r platform.Router) {
	r.Get("/login", h.LoginView)
	r.Post("/login", h.Login)
	r.Post("/logout", h.Logout)

	r.Get("/register", h.RegisterView)
	r.Post("/register", h.Register)
}

// Close implements a closer for graceful shutdown.
func (h *Handler) Close() {
	// Nothing to flush.
}
