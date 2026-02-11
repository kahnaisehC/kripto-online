package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *config) handlerGetAllGames(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := json.Marshal(cfg.games)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}
