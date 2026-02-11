package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func (cfg *config) handlerGetGame(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	gameIdString := r.PathValue("gameID")
	if gameIdString == "" {
		respondWithError(w, http.StatusBadRequest, errors.New("game id not present in url"))
		return
	}
	gameId, err := strconv.Atoi(gameIdString)
	if err != nil {

		println(gameIdString)
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	game, ok := cfg.games[gameId]
	if !ok {
		respondWithError(w, http.StatusBadRequest, errors.New("game id not found"))
		return
	}

	gameData, err := json.Marshal(game)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		respondWithError(w, http.StatusInternalServerError, err)
	}

	respondWithJSON(w, http.StatusOK, gameData)
}
