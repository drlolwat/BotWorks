package Handlers

import (
	"BotWorks/src/Config"
	"BotWorks/src/Helpers"
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"BotWorks/src/Web/Template"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func MassActionsHandler(w http.ResponseWriter, r *http.Request) {
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
		Name:     "Mass Actions",
		Url:      "massactions",
	}

	if r.PostForm.Has("add_proxies") {
		if r.FormValue("proxies") != "" {
			split := strings.Split(r.FormValue("proxies"), "\r\n")
			if len(split) > 0 {
				i := 1
				for _, data := range split {
					proxData := strings.Split(data, ":")
					if len(proxData) == 4 {
						ip := proxData[0]
						port := proxData[1]
						username := proxData[2]
						password := proxData[3]
						proxName := "New Proxy " + strconv.Itoa(i)
						_, err := Storage.GetDb().Exec("INSERT INTO proxies (ip, port, username, password, nickname, customer_id) VALUES (?, ?, ?, ?, ?, ?)", ip, port, username, password, proxName, user.Id)
						if err != nil {
							fmt.Println(err)
							return
						}
						i++
					}
				}
			}
		}
	} else if r.PostForm.Has("add_accounts") {
		if r.FormValue("accounts") != "" && r.FormValue("script_id") != "" && r.FormValue("params") != "" &&
			r.FormValue("schedule_id") != "" && r.FormValue("group_id") != "" && r.FormValue("fps") != "" {

			pin := r.FormValue("bank_pin")
			script := r.FormValue("script_id")
			params := r.FormValue("params")
			sched := r.FormValue("schedule_id")
			grp := r.FormValue("group_id")
			breaks := r.FormValue("breaks")
			fps := r.FormValue("fps")

			split := strings.Split(r.FormValue("accounts"), "\r\n")

			if len(split) > 0 {
				for _, data := range split {
					accData := strings.Split(data, ":")
					username := accData[0]
					password := accData[1]

					_, err := Storage.GetDb().Exec("INSERT INTO accounts (name, password, bank_pin, proxy_id, script_id, script_params, customer_id, logfile, pid, group_id, breaks, schedule, fps) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)", username, password, pin, -1, script, params, user.Id, "", 0, grp, breaks, sched, fps)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		}
	}

	pageData := map[string]interface{}{
		"Config":    Config.GetConfig(),
		"Scripts":   Helpers.GetScripts(),
		"Schedules": Helpers.GetSchedules(user.Id),
		"Groups":    Helpers.GetGroups(user.Id),
	}

	err := Template.GetTpl().ExecuteTemplate(w, "bootstrap_header.html", headerInfo)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "massactions.html", pageData)
	if err != nil {
		println("failed to execute tpl: massactions")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_footer.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_footer")
	}
}
