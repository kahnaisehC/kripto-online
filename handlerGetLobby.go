package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (cfg *config) handlerGetLobby(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	lobbyIDString := r.PathValue("gameID")
	lobbyIDInt, err := strconv.Atoi(lobbyIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	lobbyID := ID(lobbyIDInt)
	lobby, ok := cfg.lobbies[lobbyID]
	if !ok {
		respondWithError(w, http.StatusNotFound, err)
		return
	}

	data, err := json.Marshal(lobby)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	respondWithJSON(w, http.StatusOK, data)
}
