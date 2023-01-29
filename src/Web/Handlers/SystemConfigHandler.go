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
	}

	systemData := map[string]interface{}{
		"Config": Config.GetConfig(),
	}

	err := Template.GetTpl().ExecuteTemplate(w, "bootstrap_header.html", headerInfo)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "system.html", systemData)
	if err != nil {
		println("failed to execute tpl: schedules")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_footer.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_footer")
	}
}
