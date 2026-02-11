package main

import (
	"net/http"
)

func middlewareLogParty(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return f
}
