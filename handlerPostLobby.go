package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

var _lobbyID ID

func newLobbyID() ID {
	_lobbyID++
	return _lobbyID
}

func (cfg *config) handlerPostLobby(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	r.ParseForm()

	adminIDCookie, err := r.Cookie("userID")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, errors.New("have to login first"))
		return
	}

	adminIDInt, err := strconv.Atoi(adminIDCookie.Value)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("malformed id. have to login again"))
		return
	}
	adminID := ID(adminIDInt)

	_, ok := cfg.playerIDtoUsername[adminID]
	if !ok {
		respondWithError(w, http.StatusUnauthorized, errors.New("invalid userID. have to login again"))
		return
	}

	lobbyName := r.FormValue("lobbyName")
	if lobbyName == "" {
		respondWithError(w, http.StatusBadRequest, errors.New("not found lobby name"))
		return
	}
	lobbySizeString := r.FormValue("lobbySize")
	lobbySize, err := strconv.Atoi(lobbySizeString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("not found lobby size"))
		return
	}
	if lobbySize >= MaxLobbySize || lobbySize < 2 {
		respondWithError(w, http.StatusBadRequest, errors.New("invalid lobby size"))
		return
	}

	newLobby := NewLobby(lobbyName, lobbySize, adminID)
	data, err := json.Marshal(&newLobby)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	cfg.lobbies[newLobby.ID] = newLobby
	respondWithJSON(w, http.StatusCreated, data)

	go newLobby.Start()
}
