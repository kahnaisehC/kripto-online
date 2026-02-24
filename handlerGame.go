package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/coder/websocket"
)

type PhaseState int

type ID int

type Session string

type Card struct {
	Value int
	Palo  string
}

type GameState struct {
	Phase       PhaseState
	PointedID   ID
	Admin       ID
	Leftover    map[ID]struct{}
	Order       []ID
	Players     map[ID]struct{}
	ID          int
	Cards       []Card
	Result      int
	Connections []Connection
}

func getGameID(r *http.Request) (ID, error) {
	idString := r.PathValue("gameID")
	if idString == "" {
		return 0, errors.New("no game ID in path parameter")
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, err
	}
	return ID(id), nil
}

func (cfg *config) handlerJoinLobbyWebsocket(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userID, err := getUserID(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	gameID, err := getGameID(r)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err)
		return
	}
	lobby, ok := cfg.lobbies[gameID]
	if !ok {
		respondWithError(w, http.StatusNotFound, errors.New("not found game id"))
		return
	}
	if _, ok := lobby.Players[userID]; !ok {
		respondWithError(w, http.StatusNotFound, errors.New("player didn't join the lobby. Issue a Patch request first"))
		return
	}

	// TODO: authenticate and relate userID to connection

	con, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// WARNING:INSECURE
		Subprotocols: []string{
			"kripto",
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("unable to upgrade websocket connection"))
		return
	}
	conn := Connection{
		UserName: cfg.playerIDtoUsername[userID],
		userID:   userID,
		conn:     con,
	}
	lobby.conn = append(lobby.conn, conn)
	// read
	go func() {
		for {
			_, data, err := conn.conn.Read(context.Background())
			if err != nil {
				return
			}
			lobby.ch.ch <- string(data)
		}
	}()
	// send
	//
}
