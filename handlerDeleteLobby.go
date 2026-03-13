package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (cfg *config) handlerDeleteLobby(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	gameIDString := r.PathValue("gameID")
	gameIDInt, err := strconv.Atoi(gameIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err, []byte("game ID must be an integer"))
		return
	}
	gameID := ID(gameIDInt)
	if _, ok := cfg.lobbies[ID(gameID)]; !ok {
		respondWithError(w, http.StatusNotFound, errors.New("could'n fiind gameID "), []byte("couldn't find lobby with the given id"))
		return
	}

	userIDCookie, err := r.Cookie("userID")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, []byte("must login first"))
		return
	}
	userIDInt, err := strconv.Atoi(userIDCookie.Value)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, []byte("invalid userID, login again"))
		return
	}
	userID := ID(userIDInt)
	if userID != cfg.lobbies[gameID].AdminID {
		respondWithError(w, http.StatusUnauthorized, fmt.Errorf("%v is not the admin of %v, %v is", userID, gameID, cfg.lobbies[gameID].AdminID), []byte("you should be the admin to delete the lobby"))
		return
	}
	cfg.lobbies[gameID].Close()
	delete(cfg.lobbies, gameID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully deleted the lobby"))
}
