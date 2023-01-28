package Handlers

import (
	"BotWorks/src/Helpers"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func HandleStopBotById(w http.ResponseWriter, r *http.Request) {
	if !Helpers.IsSessionValid(r) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var user Structs.UserData
	session, _ := Sessions.Store.Get(r, Sessions.SessionName)
	userID := session.Values[Sessions.UserIDKey]
	if userID != nil {
		user.ID = userID.(int)
	} else {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	Helpers.KillBot(id, user.ID)

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
	return
}
