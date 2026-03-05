package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/coder/websocket"
)

type Session string

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
	if _, ok := lobby.userIDTouserIdx[userID]; !ok {
		respondWithError(w, http.StatusNotFound, errors.New("player didn't join the lobby. Issue a Patch request first"))
		return
	}

	con, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{
			"kripto",
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("unable to upgrade websocket connection"))
		return
	}
	conn := Connection{
		userID: userID,
		conn:   con,
	}
	lobby.conn = append(lobby.conn, conn)
	// read
	go func() {
		for {
			msgType, data, err := conn.conn.Read(context.Background())
			if err != nil {
				return
			}
			msg := lobbyChannelMessage{
				Issuer:  userID,
				Content: string(data),
				Conn:    conn.conn,
			}
			if msgType != websocket.MessageText {
				log.Println("Unsupported message type")
				return
			}

			lobby.ch.ch <- msg
		}
	}()
}

func textatthebottom() {
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
	// yyeah bro this aint doing anythinng
}
