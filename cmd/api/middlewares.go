package main

import (
	"fmt"
	"net/http"
)

func (a *application) sessionAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")

		fmt.Printf("cookie %s\n", cookie)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user := a.config.sessionStore.GetSession(cookie.Value)

		fmt.Printf("user %s\n", user)

		if user == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		r.Header.Set("user", user)

		next.ServeHTTP(w, r)
	})
}
