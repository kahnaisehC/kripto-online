package main

import (
	"net/http"
	"strconv"
)

func (cfg *config) middlewareLogin(f http.HandlerFunc) http.HandlerFunc {
	// check login

	return func(w http.ResponseWriter, r *http.Request) {

		idCookie, err := r.Cookie("ID")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		userNameCookie, err := r.Cookie("userName")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return

		}

		id, err := strconv.Atoi(idCookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		userName, ok := cfg.playerIDtoUsername[ID(id)]
		if !ok || userName != userNameCookie.Value {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		f(w, r)
	}
}
