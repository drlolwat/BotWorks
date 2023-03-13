package Handlers

import (
	"BotWorks/src/Helpers"
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"BotWorks/src/Web/Template"
	"fmt"
	"net/http"
	"strconv"
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
		if r.FormValue("name") != "" && r.FormValue("banned_action") != "" && r.FormValue("locked_action") != "" &&
			r.FormValue("complete_action") != "" {
			nickname := r.FormValue("name")
			bannedAction := r.FormValue("banned_action")
			lockedAction := r.FormValue("locked_action")
			completeAction := r.FormValue("complete_action")

			if r.PostForm.Has("update") && r.FormValue("id") != "" {
				groupId := r.FormValue("id")
				_, err := Storage.GetDb().Exec("UPDATE account_groups SET name = ?, complete_action = ?, locked_action = ?, banned_action = ? WHERE id = ? AND customer_id = ?", nickname, completeAction, lockedAction, bannedAction, groupId, user.Id)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			if r.PostForm.Has("add") {
				_, err := Storage.GetDb().Exec("INSERT INTO account_groups (name, customer_id, complete_action, locked_action, banned_action) VALUES (?, ?, ?, ?, ?)", nickname, user.Id, completeAction, lockedAction, bannedAction)
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
		groupId, _ := strconv.Atoi(r.URL.Query().Get("e"))
		groupStructData := Helpers.GetGroupById(groupId, user.Id)

		if groupStructData.Id <= 0 {
			fmt.Println("Unable to find group:", groupStructData.Id)
			return
		}

		groupsPageData["GroupData"] = groupStructData

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
