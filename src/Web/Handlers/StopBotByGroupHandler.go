package Handlers

import (
	"BotWorks/src/Helpers"
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func HandleStopBotByGroup(w http.ResponseWriter, r *http.Request) {
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

	if id > 0 {
		rows, err := Storage.GetDb().Query("SELECT id FROM accounts WHERE group_id = ? AND customer_id = ? AND pid > 0", id, user.ID)
		if err != nil {
			return
		}

		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {

			}
		}(rows)

		for rows.Next() {
			var botId int
			rows.Scan(&botId)
			Helpers.KillBot(botId, user.ID)
		}

		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
		return
	}
}
