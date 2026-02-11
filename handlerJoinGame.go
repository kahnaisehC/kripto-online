package main

import "net/http"

func (cfg *config) handlerJoinGame(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
}
