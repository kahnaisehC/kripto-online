package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type PhaseState int

type ID int

type Session [64]byte

type Card struct {
	Value int
	Palo  string
}

type GameState struct {
	Phase     PhaseState
	PointedID ID
	Leftover  map[ID]struct{}
	Order     []ID
	Admin     ID
	Players   map[ID]struct{}
	ID        int
	Cards     []Card
	Result    int
}

func (cfg *config) handlerGameWebsocket(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("gameID")
	con, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		Subprotocols: []string{
			"kripto",
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("unable to upgrade websocket connection"))
		return
	}
	log.Println("connected ws with" + id)

	wsChan := make(chan []byte, 1)
	go func() {
		for {

			_, msg, err := con.Read(context.TODO())
			if err != nil {
				con.Close(websocket.StatusNormalClosure, "finished reading")
				break
			}
			wsChan <- msg
		}
	}()

	con.Write(context.Background(), websocket.MessageText, []byte("hello"))
}
