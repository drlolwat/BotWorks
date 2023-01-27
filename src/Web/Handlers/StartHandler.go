package Handlers

import (
	"BotWorks/src/Helpers"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func HandleStartBotById(w http.ResponseWriter, r *http.Request) {
	if !Helpers.IsSessionValid(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
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

	if id > 0 {
		if Helpers.StartBotFromId(user.ID, id) {
			// send to dashboard, maybe use session for success msg
		}
	}
}

func HandleStartBotsByGroup(w http.ResponseWriter, r *http.Request) {

}

func HandleStopBotById(w http.ResponseWriter, r *http.Request) {

}

func HandleStopBotByGroup(w http.ResponseWriter, r *http.Request) {

}
