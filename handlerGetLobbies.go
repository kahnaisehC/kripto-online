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
	lobbies := make([]Lobby, 0, DefaultLobbyPageSize)

	for _, lobby := range cfg.lobbies {
		lobbies = append(lobbies, *lobby)
		if len(lobbies) == DefaultLobbyPageSize {
			break
		}
	}
	data, err := json.Marshal(lobbies)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}
