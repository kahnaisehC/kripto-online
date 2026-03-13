package main

import (
	"encoding/json"
	"net/http"
)

const (
	MinLobbyPageSize     = 1
	DefaultLobbyPageSize = 10
	MaxLobbyPageSize     = 64
)

func (cfg *config) handlerGetLobbies(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data, err := json.Marshal(cfg.lobbies)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, []byte("intrnal server error. Contact the admin"))
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}
