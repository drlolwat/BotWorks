package Helpers

import (
	"BotWorks/src/Web/Sessions"
	"net/http"
)

func IsSessionValid(r *http.Request) bool {
	session, _ := Sessions.Store.Get(r, Sessions.SessionName)
	userID := session.Values[Sessions.UserIDKey]
	if userID == nil {
		return false
	}

	return true
}
