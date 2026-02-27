package main

import (
	"errors"
	"net/http"
	"strconv"
)

func (cfg *config) handlerPatchLobby(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	gameIDString := r.PathValue("gameID")
	gameIDInt, err := strconv.Atoi(gameIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	gameID := ID(gameIDInt)
	if _, ok := cfg.lobbies[ID(gameID)]; !ok {
		respondWithError(w, http.StatusNotFound, errors.New("could'n fiind gameID "))
		return
	}
	if cfg.lobbies[gameID].Size >= len(cfg.lobbies[gameID].Players) {
		respondWithError(w, http.StatusNotAcceptable, errors.New("lobby is full"))
		return
	}

	userIDCookie, err := r.Cookie("userID")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err)
		return
	}
	userIDInt, err := strconv.Atoi(userIDCookie.Value)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err)
		return
	}
	userID := ID(userIDInt)
	userName, ok := cfg.playerIDtoUsername[userID]
	if !ok {
		respondWithError(w, http.StatusNotFound, errors.New("couldn't find user ID"))
	}
	cfg.lobbies[gameID].Players[userID] = userName
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully joined the lobby"))
}
