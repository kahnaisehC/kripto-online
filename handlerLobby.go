package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/coder/websocket"
)

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
	Conn    []Connection
	Size    int
}
type LobbyResponse struct {
	ID      ID
	Name    string
	AdminID ID
	Link    string
	Players map[ID]string
	Size    int
}

func (cfg *config) handlerGetLobbies(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	lobbiesData := make([]LobbyResponse, len(cfg.lobbies))
	for i := range cfg.lobbies {
		lobbiesData[i].ID = cfg.lobbies[i].ID
		lobbiesData[i].Name = cfg.lobbies[i].Name
		lobbiesData[i].AdminID = cfg.lobbies[i].AdminID
		lobbiesData[i].Link = cfg.lobbies[i].Link
		lobbiesData[i].Size = cfg.lobbies[i].Size
		for _, con := range cfg.lobbies[i].Conn {
			lobbiesData[i].Players[con.userID] = con.UserName
		}
	}

	data, err := json.Marshal(lobbiesData)
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

	lobbyData := LobbyResponse{
		ID:      cfg.lobbies[ID(lobbyID)].ID,
		Name:    cfg.lobbies[ID(lobbyID)].Name,
		AdminID: cfg.lobbies[lobbyID].AdminID,
		Link:    cfg.lobbies[lobbyID].Link,
		Size:    cfg.lobbies[lobbyID].Size,
	}
	for _, con := range cfg.lobbies[lobbyID].Conn {
		lobbyData.Players[con.userID] = con.UserName
	}
	data, err := json.Marshal(lobbyData)
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

	userName, ok := cfg.playerIDtoUsername[adminID]
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
		Conn:    nil,
		Size:    lobbySize,
	}
	newLobbyData := LobbyResponse{
		ID:      lobbyID,
		Name:    lobbyName,
		Link:    lobbyURL,
		AdminID: ID(adminID),
		Players: map[ID]string{
			ID(adminID): userName,
		},
		Size: lobbySize,
	}
	data, err := json.Marshal(newLobbyData)
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
}
