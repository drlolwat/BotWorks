package Actions

import (
	"BotWorks/src/Helpers"
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"database/sql"
	"fmt"
	"net/http"
)

func StopAllBotsAction(w http.ResponseWriter, r *http.Request) {
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
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rows, err := Storage.GetDb().Query("SELECT id FROM accounts WHERE customer_id = ? AND pid > 0", user.Id)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(rows)

	for rows.Next() {
		var botId int
		err := rows.Scan(&botId)
		if err != nil {
			fmt.Println(err)
			return
		}

		Helpers.KillBot(botId, user.Id)
	}

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
	return
}
