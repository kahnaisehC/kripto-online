package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/coder/websocket"
)

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
	Conn    []Connection
}

func (cfg *config) handlerGetLobbies(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data, err := json.Marshal(cfg.lobbies)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	respondWithJSON(w, http.StatusOK, data)
}

func (cfg *config) handlerGetLobby(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Add("Content-Type", "text/html")
	err := cfg.temp.ExecuteTemplate(w, "lobby", nil)
	if err != nil {
		log.Println("its mee")
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

var _lobbyID ID

func newLobbyID() ID {
	_lobbyID++
	return _lobbyID
}

func (cfg *config) handlerPostLobby(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	r.ParseForm()

	adminIdCookie, err := r.Cookie("userID")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, errors.New("have to login first"))
		return
	}

	adminId, err := strconv.Atoi(adminIdCookie.Name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("malformed id. have to login again"))
		return
	}

	lobbyName := r.FormValue("lobbyName")
	if lobbyName == "" {
		respondWithError(w, http.StatusBadRequest, errors.New("not found lobby name"))
		return
	}

	newLobby := Lobby{
		ID:      newLobbyID(),
		Name:    lobbyName,
		Link:    "",
		AdminID: ID(adminId),
		Conn:    nil,
	}

	cfg.lobbies = append(cfg.lobbies, newLobby)

	newLobbyData, _ := json.Marshal(newLobby)

	respondWithJSON(w, http.StatusCreated, newLobbyData)
}
