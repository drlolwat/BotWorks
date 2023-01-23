package Handlers

import (
	"BotWorks/src/Helpers"
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"BotWorks/src/Web/Template"
	"fmt"
	"net/http"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("userid: ", user.ID)

	// get bot data for dash
	rows, err := Storage.GetDb().Query("SELECT name,proxy_id,script_id,pid,group_id FROM accounts WHERE customer_id = ?", user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rows.Close()

	var runningBots = 0
	var stoppedBots = 0
	var uniqueScripts = 0
	var scripts []int

	for rows.Next() {
		var name string
		var proxyId int
		var scriptId int
		var pid int
		var groupId int

		err = rows.Scan(&name, &proxyId, &scriptId, &pid, &groupId)
		if err != nil {
			fmt.Println(err)
			return
		}

		if !Helpers.ContainsInt(scriptId, scripts) {
			scripts = append(scripts, scriptId)
		}

		if pid > 0 {
			runningBots += 1
		} else {
			stoppedBots += 1
		}
	}

	uniqueScripts = len(scripts)

	fmt.Println("Bots running:", runningBots)
	fmt.Println("Bots offline:", stoppedBots)
	fmt.Println("Unique scripts:", uniqueScripts)

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_header.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	data := Structs.DashboardData{RunningBots: runningBots, StoppedBots: stoppedBots, UniqueScripts: uniqueScripts, CurrentTime: "1234"}
	err = Template.GetTpl().ExecuteTemplate(w, "dashboard.html", data)
	if err != nil {
		println("failed to execute tpl: dashboard")
		return
	}
}
