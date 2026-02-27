package main

import (
	"context"
	"log"

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
type Channel struct {
	ch chan KriptoMessage
}

type Phase int

const (
	PhasePending = iota
)

type Lobby struct {
	// Static Information
	ID      ID
	Name    string
	Link    string
	Size    int
	AdminID ID

	// Dynamic Information
	Order   []ID
	Result  int
	Players map[ID]string
	conn    []Connection
	ch      Channel

	// Game Information
	Phase     Phase
	PointedID ID
	Cards     []Card
	Leftover  map[ID]struct{}
}

func (l *Lobby) Broadcast(msg string) {
	for _, con := range l.conn {
		err := con.conn.Write(context.Background(), websocket.MessageText, []byte(msg))
		if err != nil {
			println("ERR: " + err.Error())
			return
		}
	}
}

func (l *Lobby) Start() {
	for {
		msg := <-l.ch.ch
		switch msg.Type {
		case KriptoStart:
		case KriptoJoin:
			l.Broadcast("Someone joined")
		case KriptoPlay:
		case KriptoDelete:
		case KriptoPoint:
		case KriptoSolution:
		case KriptoDisconnect:
			// l.AddPlayer()
		case KriptoInvalid:
			fallthrough
		default:
			println("The message is invalid")
			log.Printf("%v\n", msg)
		}
	}
}
