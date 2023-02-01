package Handlers

import (
	"BotWorks/src/Helpers"
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"BotWorks/src/Web/Template"
	"database/sql"
	"fmt"
	"net/http"
)

func AccountsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

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

	var headerInfo = Structs.PageData{
		LoggedIn: user.Id > 0,
		Name:     "Accounts",
		Url:      "accounts",
	}

	if r.PostForm.Has("update") || r.PostForm.Has("add") {
		id := "0"
		if r.FormValue("name") != "" && r.FormValue("password") != "" && r.FormValue("bank_pin") != "" &&
			r.FormValue("proxy_id") != "" && r.FormValue("script_id") != "" && r.FormValue("params") != "" &&
			r.FormValue("schedule_id") != "" && r.FormValue("group_id") != "" && r.FormValue("fps") != "" {

			username := r.FormValue("name")
			password := r.FormValue("password")
			pin := r.FormValue("bank_pin")
			proxy := r.FormValue("proxy_id")
			script := r.FormValue("script_id")
			params := r.FormValue("params")
			sched := r.FormValue("schedule_id")
			grp := r.FormValue("group_id")
			breaks := r.FormValue("breaks")
			fps := r.FormValue("fps")

			if r.PostForm.Has("update") && r.FormValue("account_id") != "" {
				id = r.FormValue("account_id")
				_, err := Storage.GetDb().Exec("UPDATE accounts SET name = ?, password = ?, bank_pin = ?, proxy_id = ?, script_id = ?, script_params = ?, schedule = ?, group_id = ?, breaks = ?, fps = ? WHERE id = ? AND customer_id=?", username, password, pin, proxy, script, params, sched, grp, breaks, fps, id, user.Id)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			if r.PostForm.Has("add") {
				_, err := Storage.GetDb().Exec("INSERT INTO accounts (name, password, bank_pin, proxy_id, script_id, script_params, customer_id, logfile, pid, group_id, breaks, schedule, fps) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)", username, password, pin, proxy, script, params, user.Id, "", 0, grp, breaks, sched, fps)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}

	if r.URL.Query().Has("d") {
		_, err := Storage.GetDb().Exec("DELETE FROM accounts WHERE id = ? AND customer_id = ?", r.URL.Query().Get("d"), user.Id)
		if err != nil {
			fmt.Println(err)
		}
	}

	runningBotData, offlineBotData, _ := Helpers.GetBotsByCustomerId(user.Id)

	var allBots []Structs.DashboardBotData
	allBots = append(allBots, runningBotData...)
	allBots = append(allBots, offlineBotData...)

	accountPageData := map[string]interface{}{
		"Bots":           allBots,
		"AddingAccount":  r.URL.Query().Get("a"),
		"EditingAccount": r.URL.Query().Get("e"),
		"DeletingBot":    r.URL.Query().Get("d"),
		"Scripts":        Helpers.GetScripts(),
		"Proxies":        Helpers.GetProxies(user.Id),
		"Schedules":      Helpers.GetSchedules(user.Id),
		"Groups":         Helpers.GetGroups(user.Id),
	}

	if r.URL.Query().Has("e") {
		rows, err := Storage.GetDb().Query("SELECT id,name,password,bank_pin,script_params,breaks,fps,script_id,proxy_id,schedule,group_id FROM accounts WHERE id = ? AND customer_id = ?", r.URL.Query().Get("e"), user.Id)
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
			var id int
			var name string
			var pass string
			var pin int
			var params string
			var breaks string
			var fps string
			var scriptId int
			var proxyId int
			var sched int
			var groupId int

			err := rows.Scan(&id, &name, &pass, &pin, &params, &breaks, &fps, &scriptId, &proxyId, &sched, &groupId)
			if err != nil {
				fmt.Println(err)
				return
			}

			accountPageData["BotData"] = Structs.BotData{
				Id:         id,
				Name:       name,
				Pass:       pass,
				Pin:        pin,
				Params:     params,
				Breaks:     breaks,
				Fps:        fps,
				ScriptId:   scriptId,
				ProxyId:    proxyId,
				ScheduleId: sched,
				GroupId:    groupId,
			}
		}
	} else if r.URL.Query().Has("a") {
		accountPageData["BotData"] = Structs.BotData{
			Id:         0,
			Name:       "lol@wat.com",
			Pass:       "lol123",
			Pin:        -1,
			Params:     "\"\"",
			Breaks:     "\"\"",
			Fps:        "20",
			ScriptId:   -1,
			ProxyId:    -1,
			ScheduleId: -1,
			GroupId:    -1,
		}
	}

	err := Template.GetTpl().ExecuteTemplate(w, "bootstrap_header.html", headerInfo)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "accounts.html", accountPageData)
	if err != nil {
		println("failed to execute tpl: accounts")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_footer.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_footer")
	}
}
