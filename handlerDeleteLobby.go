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
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	gameID := ID(gameIDInt)
	if _, ok := cfg.lobbies[ID(gameID)]; !ok {
		respondWithError(w, http.StatusNotFound, errors.New("could'n fiind gameID "))
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
	if userID != cfg.lobbies[gameID].AdminID {
		respondWithError(w, http.StatusUnauthorized, fmt.Errorf("%v is not the admin of %v, %v is", userID, gameID, cfg.lobbies[gameID].AdminID))
		return
	}
	delete(cfg.lobbies, gameID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully deleted the lobby"))
}
