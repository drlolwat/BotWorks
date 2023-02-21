package Actions

import (
	"BotWorks/src/Helpers"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func StopBotByIdAction(w http.ResponseWriter, r *http.Request) {
	if !Helpers.IsSessionValid(r) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var user Structs.UserData
	session, _ := Sessions.Store.Get(r, Sessions.SessionName)
	userID := session.Values[Sessions.UserIDKey]
	if userID != nil {
		user.Id = userID.(int)
	} else {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	Helpers.KillBot(id, user.Id)

	http.Redirect(w, r, "/dashboard?ks", http.StatusTemporaryRedirect)
	return
}
