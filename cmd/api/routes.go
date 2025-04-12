package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) setup() chi.Router {
	r := chi.NewRouter()
	// route to check the health
	// should the used to responds to pings for healthCheck from monitoring tools
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	// all the routes should be prefixed with api
	r.Route("/api", func(r chi.Router) {
		// register the uesr
		r.Post("/register", app.register)
		// signIn the user
		r.Post("/login", app.login)
		// verify otp for mfa
		r.Get("/verify-otp/{otp}", app.verifyOtp)

		// group all the protected routes
		r.Group(func(protected chi.Router) {
			// configure middleware to check the auth
			// TODO: add more middleware to to have granular controle over roles and permissions
			protected.Use(app.sessionAuth)
			protected.Get("/users", app.getAllUsers)
		})

	})

	return r
}
