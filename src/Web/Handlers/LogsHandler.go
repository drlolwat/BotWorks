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
)

func LogsHandler(w http.ResponseWriter, r *http.Request) {
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
		Name:     "Dashboard",
		Url:      "dashboard",
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("b"))
	logFile := ""

	if id > 0 {
		config := Config.GetConfig()
		var fileName string

		rows, _ := Storage.GetDb().Query("SELECT logfile FROM accounts WHERE id = ? AND customer_id = ?", id, user.Id)

		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
		}(rows)

		for rows.Next() {
			err := rows.Scan(&fileName)
			if err != nil {
				logFile = "Could not fetch log file"
			}
		}

		filepath := config.DreambotLogsDir + "/" + fileName
		logFile = Helpers.ReadFile(filepath)
	}

	logsData := map[string]interface{}{
		"LogFile": logFile,
	}

	err := Template.GetTpl().ExecuteTemplate(w, "bootstrap_header.html", headerInfo)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "logs.html", logsData)
	if err != nil {
		println("failed to execute tpl: logs")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_footer.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_footer")
	}
}
