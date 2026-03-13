package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func (cfg *config) handlerGetLobby(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	lobbyIDString := r.PathValue("gameID")
	lobbyIDInt, err := strconv.Atoi(lobbyIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err, []byte("invalid id"))
		return
	}
	lobbyID := ID(lobbyIDInt)

	lobby, ok := cfg.lobbies[lobbyID]
	if !ok {
		respondWithError(w, http.StatusNotFound, errors.New("couldn't find lobby with that id"), []byte("there isno lobyy with the given ID"))
		return
	}

	data, err := json.Marshal(lobby)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, []byte("internal server error"))
		return
	}
	respondWithJSON(w, http.StatusOK, data)
}
