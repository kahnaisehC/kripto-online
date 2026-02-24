package main

import (
	"html"
	"net/http"
	"strconv"
)

var _userID ID = 0

func genUserID() ID {
	_userID++
	return _userID
}

func getUserID(r *http.Request) (ID, error) {
	idString, err := r.Cookie("userID")
	if err != nil {
		// annon player
		return 0, err
	}

	id, err := strconv.Atoi(idString.Value)
	if err != nil {
		// not found
		return 0, err
	}
	return ID(id), nil
}

func randomUserName() string {
	return "sloppy doggy"
}

func (cfg *config) annonLogin(w http.ResponseWriter) (ID, string) {
	userName := randomUserName()
	userName = html.EscapeString(userName)
	CookieName := http.Cookie{
		Name:  "userName",
		Value: userName,
		Path:  "/",
	}

	userId := genUserID()
	userIdString := strconv.Itoa(int(userId))
	CookieID := http.Cookie{
		Name:  "userID",
		Value: userIdString,
		Path:  "/",
	}

	http.SetCookie(w, &CookieID)
	http.SetCookie(w, &CookieName)

	cfg.playerIDtoUsername[userId] = userName
	// respondWithJSON(w, http.StatusCreated, nil)
	return userId, userName
}

func (cfg *config) handlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.ParseForm()

	userName := r.FormValue("userName")
	userName = html.EscapeString(userName)
	CookieName := http.Cookie{
		Name:  "userName",
		Value: userName,
		Path:  "/",
	}

	userId := genUserID()
	userIdString := strconv.Itoa(int(userId))
	CookieID := http.Cookie{
		Name:  "userID",
		Value: userIdString,
		Path:  "/",
	}

	http.SetCookie(w, &CookieID)
	http.SetCookie(w, &CookieName)

	cfg.playerIDtoUsername[userId] = userName
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Successfully logged in"))
}
