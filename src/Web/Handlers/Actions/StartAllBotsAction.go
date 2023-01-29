package Actions

import (
	"BotWorks/src/Helpers"
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"net/http"
)

func StartAllBotsAction(w http.ResponseWriter, r *http.Request) {
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

	rows, err := Storage.GetDb().Query("SELECT id FROM accounts WHERE customer_id = ? AND pid = 0", user.Id)
	if err != nil {
		return
	}

	Helpers.StartFromDbRows(rows, user.Id)
	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
	return
}
