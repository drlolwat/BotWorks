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

func GroupsHandler(w http.ResponseWriter, r *http.Request) {
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
		Name:     "Groups",
		Url:      "groups",
	}

	if r.PostForm.Has("update") || r.PostForm.Has("add") {
		if r.FormValue("name") != "" {

			nickname := r.FormValue("name")

			if r.PostForm.Has("update") && r.FormValue("id") != "" {
				groupId := r.FormValue("id")
				_, err := Storage.GetDb().Exec("UPDATE account_groups SET name = ? WHERE id = ? AND customer_id = ?", nickname, groupId, user.Id)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			if r.PostForm.Has("add") {
				_, err := Storage.GetDb().Exec("INSERT INTO account_groups (name, customer_id) VALUES (?, ?)", nickname, user.Id)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}

	if r.URL.Query().Has("d") {
		_, err := Storage.GetDb().Exec("DELETE FROM account_groups WHERE id = ? AND customer_id = ?", r.URL.Query().Get("d"), user.Id)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = Storage.GetDb().Exec("UPDATE accounts SET group_id = 0 WHERE group_id = ? AND customer_id = ?", r.URL.Query().Get("d"), user.Id)
	}

	groupsPageData := map[string]interface{}{
		"Groups":        Helpers.GetGroups(user.Id),
		"AddingGroup":   r.URL.Query().Get("a"),
		"EditingGroup":  r.URL.Query().Get("e"),
		"DeletingGroup": r.URL.Query().Get("d"),
	}

	if r.URL.Query().Has("e") {
		rows, err := Storage.GetDb().Query("SELECT id,name FROM account_groups WHERE id = ? AND customer_id = ?", r.URL.Query().Get("e"), user.Id)
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

			err = rows.Scan(&id, &name)
			if err != nil {
				fmt.Println(err)
				return
			}

			groupsPageData["GroupData"] = Structs.GroupData{
				Id:   id,
				Name: name,
			}
		}

	} else if r.URL.Query().Has("a") {
		groupsPageData["GroupData"] = Structs.GroupData{
			Id:   0,
			Name: "New Group",
		}
	}

	err := Template.GetTpl().ExecuteTemplate(w, "bootstrap_header.html", headerInfo)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "groups.html", groupsPageData)
	if err != nil {
		println("failed to execute tpl: groups")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_footer.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_footer")
	}
}
