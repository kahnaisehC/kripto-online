package main

import (
	"log"
	"net/http"
	"strconv"
)

type Lobby struct {
	id   int
	Name string
	Link string
}

var ExampleLobbiesParams = []Lobby{
	{
		id:   0,
		Name: "example lobby",
		Link: "https://localhost:12312/1",
	},
}

func (cfg *config) handlerGetLobbies(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Add("Content-Type", "text/html")
	err := cfg.temp.ExecuteTemplate(w, "lobby", ExampleLobbiesParams)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (cfg *config) handlerGetLobby(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Add("Content-Type", "text/html")
	err := cfg.temp.ExecuteTemplate(w, "", ExampleLobbiesParams[0])
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (cfg *config) handlerPostLobby(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type response struct{}

	defer r.Body.Close()

	_, err := r.Cookie("sessionID")
	if err != nil {
		currSessionIDString := strconv.Itoa(cfg.currSessionID)
		sessionCookie := http.Cookie{
			Name:  "sessionID",
			Value: currSessionIDString,
		}
		cfg.sessionIDs[cfg.currSessionID] = currSessionIDString
		cfg.currSessionID++
		http.SetCookie(w, &sessionCookie)
	}

	w.WriteHeader(http.StatusCreated)
}
