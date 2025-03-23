package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) setup() chi.Router {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted)
		})

		r.Post("/register", app.register)
		r.Post("/login", app.login)

		r.Group(func(protected chi.Router) {
			protected.Use(app.sessionAuth)
			protected.Get("/users", app.getAllUsers)
		})

	})

	return r
}
