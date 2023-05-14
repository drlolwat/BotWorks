package Handlers

import (
	"BotWorks/src/Config"
	"BotWorks/src/Helpers"
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"BotWorks/src/Web/Template"
	"database/sql"
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

	config := Config.GetConfig()

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
						proxName := ip + ":" + port
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
		if r.FormValue("accounts") != "" && r.FormValue("script_id") != "" &&
			r.FormValue("schedule_id") != "" && r.FormValue("group_id") != "" && r.FormValue("fps") != "" &&
			r.FormValue("proxies") != "" && r.FormValue("world") != "" {

			breaks := ""
			params := ""

			script := r.FormValue("script_id")
			sched := r.FormValue("schedule_id")
			grp := r.FormValue("group_id")
			fps := r.FormValue("fps")
			world := r.FormValue("world")

			if r.FormValue("breaks") != "" {
				breaks = r.FormValue("breaks")
			}

			if r.FormValue("params") != "" {
				params = r.FormValue("params")
			}

			proxyOption, err := strconv.Atoi(r.FormValue("proxies"))

			if err != nil {
				fmt.Println(err)
			}

			split := strings.Split(r.FormValue("accounts"), "\r\n")

			if len(split) > 0 {
				for _, data := range split {
					accData := strings.Split(data, ":")
					username := accData[0]
					password := accData[1]
					proxyId := -1

					if proxyOption == 1 {
						proxyId = Helpers.GetRandomProxy(user.Id, 1, true)
						if proxyId <= 0 {
							fmt.Println("No available unused proxies in the system")
							break
						}

					} else if proxyOption == 2 {
						err = Storage.GetDb().QueryRow("SELECT id FROM proxies WHERE customer_id = ? ORDER BY RANDOM() LIMIT 1", user.Id).Scan(&proxyId)

						if err != nil || err == sql.ErrNoRows {
							fmt.Println("No available proxies in the system")
							break
						}
					}

					_, err := Storage.GetDb().Exec("INSERT INTO accounts (name, password, bank_pin, proxy_id, script_id, script_params, customer_id, logfile, pid, group_id, breaks, schedule, fps, world) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)", username, password, "-1", proxyId, script, params, user.Id, "", 0, grp, breaks, sched, fps, world)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	} else if r.PostForm.Has("assign_all") {
		rows, err := Storage.GetDb().Query("SELECT id,group_id FROM accounts WHERE proxy_id <= 0 AND customer_id = ?", user.Id)

		if err != nil {
			fmt.Println(err)
		}

		for rows.Next() {
			var id int
			var group int
			err := rows.Scan(&id, &group)

			if err != nil {
				fmt.Println(err)
				continue
			}

			newProxy := Helpers.GetRandomProxy(user.Id, 1, true)
			_, _ = Storage.GetDb().Exec("UPDATE accounts SET proxy_id = ? WHERE id = ?", newProxy, id)
		}
	}

	pageData := map[string]interface{}{
		"Config":    config,
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
