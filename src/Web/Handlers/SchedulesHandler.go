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

func SchedulesHandler(w http.ResponseWriter, r *http.Request) {
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
		Name:     "Schedules",
		Url:      "schedules",
	}

	fmt.Println(r.PostForm)
	if r.PostForm.Has("update") || r.PostForm.Has("add") {
		if r.FormValue("name") != "" && r.FormValue("start_hour") != "" && r.FormValue("start_min") != "" &&
			r.FormValue("end_hour") != "" && r.FormValue("end_min") != "" {

			nickname := r.FormValue("name")
			startHour := r.FormValue("start_hour")
			startMin := r.FormValue("start_min")
			endHour := r.FormValue("end_hour")
			endMin := r.FormValue("end_min")

			if r.PostForm.Has("update") && r.FormValue("id") != "" {
				scheduleId := r.FormValue("id")
				_, err := Storage.GetDb().Exec("UPDATE account_schedules SET name = ?, start_hour = ?, start_min = ?, end_hour = ?, end_min = ? WHERE id = ? AND customer_id = ?", nickname, startHour, startMin, endHour, endMin, scheduleId, user.Id)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			if r.PostForm.Has("add") {
				_, err := Storage.GetDb().Exec("INSERT INTO account_schedules (name, start_hour, start_min, end_hour, end_min, customer_id) VALUES (?, ?, ?, ?, ?, ?)", nickname, startHour, startMin, endHour, endMin, user.Id)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}

	if r.URL.Query().Has("d") {
		_, err := Storage.GetDb().Exec("DELETE FROM account_schedules WHERE id = ? AND customer_id = ?", r.URL.Query().Get("d"), user.Id)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = Storage.GetDb().Exec("UPDATE accounts SET schedule = 0 WHERE schedule = ? AND customer_id = ?", r.URL.Query().Get("d"), user.Id)
	}

	groupsPageData := map[string]interface{}{
		"Schedules":        Helpers.GetSchedules(user.Id),
		"AddingSchedule":   r.URL.Query().Get("a"),
		"EditingSchedule":  r.URL.Query().Get("e"),
		"DeletingSchedule": r.URL.Query().Get("d"),
	}

	if r.URL.Query().Has("e") {
		rows, err := Storage.GetDb().Query("SELECT id,name,start_hour,start_min,end_hour,end_min FROM account_schedules WHERE id = ? AND customer_id = ?", r.URL.Query().Get("e"), user.Id)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(rows)

		for rows.Next() {
			var id int
			var name string
			var startHour int
			var startMin int
			var endHour int
			var endMin int

			err = rows.Scan(&id, &name, &startHour, &startMin, &endHour, &endMin)
			if err != nil {
				fmt.Println(err)
				return
			}

			groupsPageData["ScheduleData"] = Structs.ScheduleData{
				Id:          id,
				Name:        name,
				StartHour:   startHour,
				StartMinute: startMin,
				EndHour:     endHour,
				EndMinute:   endMin,
			}
		}

	} else if r.URL.Query().Has("a") {
		groupsPageData["ScheduleData"] = Structs.ScheduleData{
			Id:          0,
			Name:        "New schedule",
			StartHour:   00,
			StartMinute: 00,
			EndHour:     00,
			EndMinute:   00,
		}
	}

	err := Template.GetTpl().ExecuteTemplate(w, "bootstrap_header.html", headerInfo)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "schedules.html", groupsPageData)
	if err != nil {
		println("failed to execute tpl: schedules")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_footer.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_footer")
	}
}
