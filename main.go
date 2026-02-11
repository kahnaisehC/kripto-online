package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func respondWithError(w http.ResponseWriter, code int, err error) {
	log.Println(err.Error())
	w.WriteHeader(code)
}

func respondWithJSON(w http.ResponseWriter, code int, data []byte) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

type config struct {
	connString    string
	temp          *template.Template
	sessionIDs    map[int]string
	currSessionID int
	playerIDs     map[int]struct{}
	games         map[int]GameState
	sessions      map[Session]ID
}

func main() {
	temp, err := template.ParseGlob("front/*.html")
	if err != nil {
		panic(err)
	}
	cfg := config{
		connString:    sampleConnString,
		temp:          temp,
		sessionIDs:    map[int]string{},
		currSessionID: 1,
		playerIDs:     map[int]struct{}{},
		games:         sampleGames,
	}
	mux := http.NewServeMux()

	// static assets
	mux.Handle("/css/", http.StripPrefix("", http.FileServer(http.Dir("./front"))))

	// frontend endpoints
	// TODO: Implement this
	mux.HandleFunc("GET /lobby", cfg.handlerGetLobby)
	// TODO: Implement this
	mux.HandleFunc("POST /lobby", middlewareLogParty(cfg.handlerPostLobby))

	// API
	mux.HandleFunc("GET /api/game/", cfg.handlerGetAllGames)
	// TODO: Implement this
	mux.HandleFunc("POST /api/game/", cfg.handlerCreateGame)

	mux.HandleFunc("GET /api/game/{gameID}", cfg.handlerGetGame)
	// TODO: Implement this
	mux.HandleFunc("POST /api/game/{gameID}", cfg.handlerJoinGame)

	// TODO: Implement this
	mux.HandleFunc("api/game/{gameID}/ws", cfg.handlerGameWebsocket)

	serverChannel := make(chan error, 1)
	go func() {
		serverChannel <- http.ListenAndServe(cfg.connString, mux)
	}()

	fmt.Println("Listening in port http://" + cfg.connString)

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
