package service

import "github.com/titpetric/platform"

func (h *Service) Mount(r platform.Router) error {
	r.Get("/login", h.LoginView)
	r.Post("/login", h.Login)

	r.Get("/logout", h.LogoutView)
	r.Post("/logout", h.Logout)

	r.Get("/register", h.RegisterView)
	r.Post("/register", h.Register)

	return nil
}
