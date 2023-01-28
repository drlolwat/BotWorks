package Handlers

import (
	"BotWorks/src/Helpers"
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"database/sql"
	"fmt"
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
		Helpers.StartBotFromId(user.ID, id)
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
		return
	}
}

func HandleStartBotsByGroup(w http.ResponseWriter, r *http.Request) {
	if !Helpers.IsSessionValid(r) {
		http.Redirect(w, r, "/", http.StatusForbidden)
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
		rows, err := Storage.GetDb().Query("SELECT id FROM accounts WHERE group_id = ? AND customer_id = ? AND pid = 0", id, user.ID)
		if err != nil {
			return
		}

		Helpers.StartFromDbRows(rows, user.ID)
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
		return
	}
}

func HandleStartAll(w http.ResponseWriter, r *http.Request) {
	if !Helpers.IsSessionValid(r) {
		http.Redirect(w, r, "/", http.StatusForbidden)
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

	rows, err := Storage.GetDb().Query("SELECT id FROM accounts WHERE customer_id = ? AND pid = 0", user.ID)
	if err != nil {
		return
	}

	Helpers.StartFromDbRows(rows, user.ID)
	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
	return
}

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
		rows, err := Storage.GetDb().Query("SELECT pid FROM accounts WHERE group_id = ? AND customer_id = ? AND pid > 0", id, user.ID)
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

func HandleStopAll(w http.ResponseWriter, r *http.Request) {
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
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rows, err := Storage.GetDb().Query("SELECT pid FROM accounts WHERE customer_id = ? AND pid > 0", user.ID)
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
		Helpers.KillBot(botId, user.ID)
	}

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
	return
}
