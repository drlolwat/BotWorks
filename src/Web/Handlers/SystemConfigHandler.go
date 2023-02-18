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
)

func SystemConfigHandler(w http.ResponseWriter, r *http.Request) {
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
		Name:     "System",
		Url:      "system",
	}

	if r.PostForm.Has("update") {
		if r.FormValue("jar") != "" && r.FormValue("logs") != "" && r.FormValue("user") != "" && r.FormValue("pass") != "" {
			dbJar := r.FormValue("jar")
			dbLogs := r.FormValue("logs")
			dbUser := r.FormValue("user")
			dbPass := r.FormValue("pass")

			_, err := Storage.GetDb().Exec("UPDATE configuration SET db_jar = ?, db_logs = ?, db_username = ?, db_password = ?", dbJar, dbLogs, dbUser, dbPass)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else if r.PostForm.Has("update_system") {
		if r.FormValue("banned_action") != "" && r.FormValue("complete_action") != "" {
			fmt.Println(r.PostForm)
			action := r.FormValue("banned_action")
			completeAction := r.FormValue("complete_action")

			fmt.Println(action, completeAction)
			fmt.Println(Config.GetConfig().BannedAction)
			fmt.Println(Config.GetConfig().CompleteAction)

			_, err := Storage.GetDb().Exec("UPDATE configuration SET banned_action = ?", action, completeAction)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else if r.PostForm.Has("oh_shit_button") {
		_, err := Storage.GetDb().Exec("DELETE FROM account_groups WHERE id > 0")

		if err != nil {
			fmt.Println("A", err)
		}

		_, err = Storage.GetDb().Exec("DELETE FROM account_schedules WHERE id > 0")

		if err != nil {
			fmt.Println("B", err)
		}
		_, err = Storage.GetDb().Exec("DELETE FROM accounts")

		if err != nil {
			fmt.Println("C", err)
		}
		_, err = Storage.GetDb().Exec("DELETE FROM customers")

		if err != nil {
			fmt.Println("D", err)
		}
		_, err = Storage.GetDb().Exec("DELETE FROM proxies WHERE id > 0")

		if err != nil {
			fmt.Println("E", err)
		}
		_, err = Storage.GetDb().Exec("DELETE FROM scripts WHERE id > 0")

		if err != nil {
			fmt.Println("F", err)
		}

		_, err = Storage.GetDb().Exec("UPDATE configuration SET db_jar = ?, db_logs = ?, db_username = ?, db_password = ?", "C:\\", "C:\\", "user", "pass")

		if err != nil {
			fmt.Println("G", err)
		}

		fmt.Println("The entire panel has been wiped")

		session, _ := Sessions.Store.Get(r, Sessions.SessionName)
		session.Options.MaxAge = -1
		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err)
		}

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	systemData := map[string]interface{}{
		"Config": Config.GetConfig(),
		"Groups": Helpers.GetGroups(user.Id),
	}

	err := Template.GetTpl().ExecuteTemplate(w, "bootstrap_header.html", headerInfo)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "system.html", systemData)
	if err != nil {
		println("failed to execute tpl: system")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_footer.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_footer")
	}
}
