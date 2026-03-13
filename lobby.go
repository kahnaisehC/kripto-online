package main

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/coder/websocket"
	"github.com/kahnaisehC/kripto_online/internal/kriptogame"
)

const MaxLobbySize = 32

type ID int

var _counter ID = 1

func GenLobbyID() ID {
	_counter++
	return _counter
}

type Connection struct {
	userID ID
	conn   *websocket.Conn
}

type lobbyChannelMessage struct {
	Issuer  ID
	Content string
	Conn    *websocket.Conn
}

const (
	CloseLobbyMessage = "finish"
	StartLobbyMessage = "start"
)

type Channel struct {
	ch chan lobbyChannelMessage
}

type LobbyState int

const (
	// LobbyStateClosed the game is over
	LobbyStateClosed LobbyState = iota
	// LobbyStatePlaying the game is currently being pllayed
	LobbyStatePlaying
	// LobbyStatePending means that the game have not started still. The admin should send StartMessage to start it
	LobbyStatePending
)

type Lobby struct {
	// Static Information
	ID       ID
	Name     string
	Link     string
	AdminID  ID
	CurrSize int
	Size     int

	//
	userIDTouserIdx map[ID]int

	// Dynamic Information
	conn      []Connection `json:"-"`
	connMutex sync.RWMutex `json:"-"`
	ch        Channel      `json:"-"`
	State     LobbyState
	// Game Information
	Game kriptogame.Game `json:"-"`
}

func NewLobby(Name string, Size int, AdminID ID) *Lobby {
	if Size < 2 {
		return nil
	}
	lobbyID := newLobbyID()
	// TODO: give the actual url??
	lobbyURL := "/lobby/" + strconv.Itoa(int(lobbyID))
	return &Lobby{
		ID:      lobbyID,
		Name:    Name,
		Link:    lobbyURL,
		AdminID: AdminID,
		userIDTouserIdx: map[ID]int{
			AdminID: 0,
		},
		Size:      Size,
		CurrSize:  1,
		conn:      nil,
		connMutex: sync.RWMutex{},
		ch: Channel{
			ch: make(chan lobbyChannelMessage, 10),
		},
		State: LobbyStatePending,
	}
}

func (l *Lobby) Broadcast(msg string) {
	l.connMutex.RLock()
	for _, con := range l.conn {
		err := con.conn.Write(context.Background(), websocket.MessageText, []byte(msg))
		if err != nil {
			println("ERR: " + err.Error())
		}
	}
	l.connMutex.RUnlock()
}

func (l *Lobby) Close() {
	l.connMutex.Lock()
	for _, conn := range l.conn {
		conn.conn.Close(websocket.StatusNormalClosure, "The lobby is closing")
	}
	l.conn = nil
	l.connMutex.Unlock()

	l.State = LobbyStateClosed
	l.ch.ch <- lobbyChannelMessage{
		Content: CloseLobbyMessage,
	}
	// TODO: Store the game
}

// Join adds a userID to the lobby
func (l *Lobby) Join(userID ID) error {
	if l.State != LobbyStatePending {
		return errors.New("game already started")
	}
	if l.CurrSize >= l.Size {
		return errors.New("lobby is full")
	}

	if _, ok := l.userIDTouserIdx[userID]; ok {
		return errors.New("user is already in the lobby")
	}
	l.userIDTouserIdx[userID] = l.CurrSize
	l.CurrSize++
	return nil
}

func (l *Lobby) Start() {
	for {
		msg := <-l.ch.ch
		if msg.Content == CloseLobbyMessage && l.State == LobbyStateClosed {
			// TODO: add logging of lobby closing
		L:
			for {
				select {
				case <-l.ch.ch:
				default:
					break L
				}
			}
			state := l.Game.GetStateString()
			l.Broadcast(state)
			return
		}
		if msg.Content == StartLobbyMessage {
			if l.State == LobbyStatePending && msg.Issuer == l.AdminID {
				l.Game = kriptogame.NewGame(l.CurrSize)
				l.State = LobbyStatePlaying
			}
			state := l.Game.GetStateString()
			l.Broadcast(state)
			continue
		}
		if l.State != LobbyStatePlaying {
			// DEBUG
			l.Broadcast("lobby is not being played")
			l.Broadcast(msg.Content)
			continue
		}

		kriptoMsg, err := l.Game.ParseMessage(msg.Content)
		if err != nil {
			// TODO: fix this login
			println(err.Error())
			continue
		}
		kriptoMsg.IssuerIdx = l.userIDTouserIdx[msg.Issuer]
		err = l.Game.CheckMessageValidity(kriptoMsg)
		if err != nil {
			// TODO: fix this login
			println(err.Error())
			continue
		}

		ok := l.Game.ExecuteUnsafe(kriptoMsg)
		if !ok {
			println(err.Error())
			continue
		}
		state := l.Game.GetStateString()
		l.Broadcast(state)

		if l.Game.Phase == kriptogame.PhaseFinished {
			l.Close()
		}

	}
}
