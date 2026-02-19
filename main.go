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

type DB interface{}

type voidDB struct{}

type config struct {
	connString         string
	temp               *template.Template
	sessionIDs         map[int]string
	currSessionID      int
	playerIDtoUsername map[ID]string
	games              map[int]GameState
	sessions           map[string]ID
	db                 DB
	lobbies            []Lobby
	jwtKey             string
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
		games:              sampleGames,
		db:                 voidDB{},
		jwtKey:             os.Getenv("JWT_KEY"),
		lobbies: []Lobby{
			{
				ID:      10,
				Name:    "IAN",
				Link:    "!!",
				AdminID: 1000,
				Conn: []Connection{
					{
						UserName: "ianuser",
						userID:   1000,
						conn:     nil,
					},
				},
			},
		},
	}

	mux := http.NewServeMux()

	// static assets
	mux.Handle("/css/", http.StripPrefix("", http.FileServer(http.Dir("./front"))))

	// frontend endpoints
	// TODO: Implement this
	mux.HandleFunc("GET /lobby", cfg.handlerTemplate("lobby"))

	// TODO: Implement htis
	mux.HandleFunc("GET /login", cfg.handlerTemplate("login"))
	mux.HandleFunc("/", cfg.handlerTemplate("login"))
	// API
	// NOTE: This only gives the client a cookie with a random number
	// that will be used to identify it later
	// TODO: Use JWT or something more sophisticated
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)

	mux.HandleFunc("POST /api/lobby", cfg.handlerPostLobby)

	// TODO: Implement this
	mux.HandleFunc("GET /api/lobby", cfg.handlerGetLobbies)
	// TODO: Implement this
	mux.HandleFunc("GET /api/game", cfg.handlerGetAllGames)
	// TODO: Implement this
	mux.HandleFunc("POST /api/game", cfg.handlerCreateGame)

	mux.HandleFunc("GET /api/game/{gameID}", cfg.handlerGetGame)
	// TODO: Implement this
	mux.HandleFunc("POST /api/game/{gameID}", cfg.handlerJoinGame)

	// TODO: Implement this
	mux.HandleFunc("/api/game/{gameID}/ws", cfg.handlerGameWebsocket)

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
