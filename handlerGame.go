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

func (cfg *config) handlerGameWebsocket(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("gameID")
	if idString == "" {
		respondWithError(w, http.StatusBadRequest, errors.New("no game ID in path parameter"))
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("game id is not a valid number"))
		return
	}
	game, ok := cfg.games[id]
	if !ok {
		respondWithError(w, http.StatusNotFound, errors.New("not found game id"))
		return
	}
	_ = game

	// TODO: authenticate and relate userID to connection
	userID, _ := getUserID(r)

	con, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// WARNING:INSECURE
		InsecureSkipVerify: true,
		Subprotocols: []string{
			"kripto",
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("unable to upgrade websocket connection"))
		return
	}

	game.Connections = append(game.Connections, Connection{
		conn:   con,
		userID: userID,
	})
	con.Write(context.Background(), websocket.MessageText, []byte("hello"))
}
