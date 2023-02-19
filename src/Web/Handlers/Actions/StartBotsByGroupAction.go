package Actions

import (
	"BotWorks/src/Helpers"
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func StartBotsByGroupAction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if !Helpers.IsSessionValid(r) {
		http.Redirect(w, r, "/", http.StatusForbidden)
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

	if id > 0 {
		rows, err := Storage.GetDb().Query("SELECT id FROM accounts WHERE group_id = ? AND customer_id = ? AND pid = 0", id, user.Id)
		if err != nil {
			return
		}

		delay := 0

		if r.PostForm.Has("delay") {
			delay, err = strconv.Atoi(r.FormValue("delay"))
		}

		if err != nil {
			delay = 0
			fmt.Print(err)
		}

		Helpers.StartFromDbRows(rows, user.Id, delay)
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
		return
	}
}
