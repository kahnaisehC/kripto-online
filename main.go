package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var defaultError = errors.New("default error")

func respondWithError(w http.ResponseWriter, code int, err error, data []byte) {
	if err == nil {
		err = defaultError
	}
	log.Println(err.Error())
	w.Header().Add("Content-type", "text/plain")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, data []byte) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

type config struct {
	connString         string
	temp               *template.Template
	sessionIDs         map[int]string
	currSessionID      int
	playerIDtoUsername map[ID]string
	sessions           map[string]ID
	lobbies            map[ID]*Lobby
}

func main() {
	temp, err := template.ParseGlob("front/*.html")
	if err != nil {
		panic(err)
	}
	cfg := config{
		connString:         sampleConnString,
		temp:               temp,
		sessionIDs:         map[int]string{},
		currSessionID:      1,
		playerIDtoUsername: map[ID]string{},
		lobbies:            make(map[ID]*Lobby),
	}

	mux := http.NewServeMux()

	// static assets
	mux.Handle("/css/", http.StripPrefix("", http.FileServer(http.Dir("./front"))))
	mux.Handle("/js/", http.StripPrefix("", http.FileServer(http.Dir("./front"))))

	// frontend endpoints
	// TODO: Implement this
	mux.HandleFunc("GET /lobby", cfg.handlerTemplate("lobbies"))
	mux.HandleFunc("GET /lobby/{gameID}", cfg.handlerTemplate("lobby"))
	mux.HandleFunc("GET /login", cfg.handlerTemplate("login"))

	// mux.HandleFunc("/", cfg.handlerTemplate("login"))
	// API
	// NOTE: This only gives the client a cookie with a random number
	// that will be used to identify it later
	// TODO: Use JWT or something more sophisticated
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)

	// TODO: ADD PAGINATION
	mux.HandleFunc("GET /api/lobby", cfg.handlerGetLobbies)
	mux.HandleFunc("POST /api/lobby", cfg.handlerPostLobby)

	mux.HandleFunc("GET /api/lobby/{gameID}", cfg.handlerGetLobby)
	mux.HandleFunc("DELETE /api/lobby/{gameID}", cfg.handlerDeleteLobby)
	mux.HandleFunc("PATCH /api/lobby/{gameID}", cfg.handlerPatchLobby)

	// Websockets
	// TODO: Implement this
	mux.HandleFunc("/api/game/{gameID}/join", cfg.handlerJoinLobbyWebsocket)

	serverChannel := make(chan error, 1)
	go func() {
		fmt.Println("Serving at: http://" + cfg.connString)
		serverChannel <- http.ListenAndServe(cfg.connString, mux)
	}()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)

	select {
	case msg1 := <-sigChannel:
		if msg1 == os.Interrupt {
			fmt.Println("Shutting down server")
		}
	case msg2 := <-serverChannel:
		if msg2 == <-serverChannel {
			fmt.Println("Error listening and serving")
		}
	}
}
