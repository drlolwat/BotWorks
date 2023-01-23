package Handlers

import (
	"BotWorks/src/Web/Sessions"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Sessions.Store.Get(r, Sessions.SessionName)
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
	return
}
