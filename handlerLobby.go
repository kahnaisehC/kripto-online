package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/coder/websocket"
)

type void struct{}

const MaxLobbySize = 32

var _counter ID = 1

func GenLobbyID() ID {
	_counter++
	return _counter
}

type Connection struct {
	UserName string
	userID   ID
	conn     *websocket.Conn
}

type Lobby struct {
	ID      ID
	Name    string
	Link    string
	AdminID ID
	conn    []Connection
	Size    int
	Players map[ID]struct{}
}

const PageSize = 10

func (cfg *config) handlerGetLobbies(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	lobbies := make([]Lobby, 0, PageSize)

	for _, lobby := range cfg.lobbies {
		lobbies = append(lobbies, lobby)
		if len(lobbies) == PageSize {
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
	lobbyID := newLobbyID()
	lobbyURL := "/lobby/" + strconv.Itoa(int(lobbyID))

	newLobby := Lobby{
		ID:      lobbyID,
		Name:    lobbyName,
		Link:    lobbyURL,
		AdminID: ID(adminID),
		conn:    nil,
		Size:    lobbySize,
		Players: map[ID]struct{}{
			ID(adminID): {},
		},
	}
	data, err := json.Marshal(newLobby)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	cfg.lobbies[lobbyID] = newLobby
	respondWithJSON(w, http.StatusCreated, data)
}

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
	userID, err := strconv.Atoi(userIDCookie.Value)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err)
		return
	}
	if ID(userID) != cfg.lobbies[ID(gameID)].AdminID {
		respondWithError(w, http.StatusUnauthorized, fmt.Errorf("%v is not the admin of %v, %v is", userID, gameID, cfg.lobbies[ID(gameID)].AdminID))
		return
	}
	delete(cfg.lobbies, ID(gameID))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully deleted the lobby"))
}

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
	cfg.lobbies[gameID].Players[userID] = struct{}{}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully deleted the lobby"))
}
