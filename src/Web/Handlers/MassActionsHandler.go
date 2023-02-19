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
			r.FormValue("proxies") != "" {

			breaks := ""
			params := ""

			script := r.FormValue("script_id")
			sched := r.FormValue("schedule_id")
			grp := r.FormValue("group_id")
			fps := r.FormValue("fps")

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
						res, _ := Storage.GetDb().Query("SELECT id FROM proxies WHERE customer_id = ? AND bans = 0 ORDER BY RANDOM()", user.Id)

						for res.Next() {
							var dbId int
							err = res.Scan(&dbId)

							if err != nil {
								fmt.Println(err)
								break
							}

							if dbId > 0 {
								var used int
								err = Storage.GetDb().QueryRow("SELECT proxy_id FROM accounts WHERE proxy_id = ?", dbId).Scan(&used)

								if err != nil || err == sql.ErrNoRows {
									proxyId = dbId
									break
								}
							}
						}

						_ = res.Close()
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

					_, err := Storage.GetDb().Exec("INSERT INTO accounts (name, password, bank_pin, proxy_id, script_id, script_params, customer_id, logfile, pid, group_id, breaks, schedule, fps) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)", username, password, "-1", proxyId, script, params, user.Id, "", 0, grp, breaks, sched, fps)
					if err != nil {
						fmt.Println(err)
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
