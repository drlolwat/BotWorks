package Actions

import (
	"BotWorks/src/Helpers"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func StartBotByIdAction(w http.ResponseWriter, r *http.Request) {
	if !Helpers.IsSessionValid(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
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

	go func() {
		if id > 0 {
			botPid, logFile := Helpers.StartBotFromId(user.Id, id)
			err := Helpers.UpdateBotStatus(id, botPid, logFile)
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}
	}()

	http.Redirect(w, r, "/dashboard?ss", http.StatusTemporaryRedirect)
}
