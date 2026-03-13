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
		respondWithError(w, http.StatusBadRequest, err, []byte("given game ID is not an integer"))
		return
	}
	gameID := ID(gameIDInt)
	if _, ok := cfg.lobbies[gameID]; !ok {
		respondWithError(w, http.StatusNotFound, errors.New("could'n fiind lobby with the given ID"), []byte("couldnt find lobby with the given ID"))
		return
	}

	userIDCookie, err := r.Cookie("userID")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, []byte("couldn't find user ID"))
		return
	}
	userIDInt, err := strconv.Atoi(userIDCookie.Value)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, []byte("cookie login is not valid. login again"))
		return
	}
	userID := ID(userIDInt)
	_, ok := cfg.playerIDtoUsername[userID]
	if !ok {
		respondWithError(w, http.StatusNotFound, errors.New("couldn't find user ID"), []byte("there is no player with the given ID"))
		return
	}

	err = cfg.lobbies[gameID].Join(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, []byte("couldn't join the lobby"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully joined the lobby"))
}
