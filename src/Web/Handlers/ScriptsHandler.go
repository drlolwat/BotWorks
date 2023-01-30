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

func ScriptsHandler(w http.ResponseWriter, r *http.Request) {
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
		Name:     "Scripts",
		Url:      "scripts",
	}

	if r.PostForm.Has("update") || r.PostForm.Has("add") {
		if r.FormValue("nickname") != "" && r.FormValue("realname") != "" {

			nickname := r.FormValue("nickname")
			sdnName := r.FormValue("realname")

			if r.PostForm.Has("update") && r.FormValue("id") != "" {
				scriptId := r.FormValue("id")
				_, err := Storage.GetDb().Exec("UPDATE scripts SET name = ?, backend_name = ? WHERE id = ?", nickname, sdnName, scriptId)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			if r.PostForm.Has("add") {
				_, err := Storage.GetDb().Exec("INSERT INTO scripts (name, backend_name) VALUES (?, ?)", nickname, sdnName)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}

	if r.URL.Query().Has("d") {
		_, err := Storage.GetDb().Exec("DELETE FROM scripts WHERE id = ? AND id > 0", r.URL.Query().Get("d"))
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = Storage.GetDb().Exec("UPDATE accounts SET script_id = 0 WHERE script_id = ? AND script_id > 0", r.URL.Query().Get("d"))
	}

	scriptsPageData := map[string]interface{}{
		"Scripts":        Helpers.GetScripts(),
		"AddingScript":   r.URL.Query().Get("a"),
		"EditingScript":  r.URL.Query().Get("e"),
		"DeletingScript": r.URL.Query().Get("d"),
	}

	if r.URL.Query().Has("e") {
		rows, err := Storage.GetDb().Query("SELECT id,name,backend_name FROM scripts WHERE id = ?", r.URL.Query().Get("e"))
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
			var sdnName string

			err = rows.Scan(&id, &name, &sdnName)
			if err != nil {
				fmt.Println(err)
				return
			}

			scriptsPageData["ScriptData"] = Structs.ScriptData{
				Id:         id,
				PublicName: name,
				SdnName:    sdnName,
			}
		}

	} else if r.URL.Query().Has("a") {
		scriptsPageData["ScriptData"] = Structs.ScriptData{
			Id:         0,
			PublicName: "New script",
			SdnName:    "Name on SDN or local script",
		}
	}

	err := Template.GetTpl().ExecuteTemplate(w, "bootstrap_header.html", headerInfo)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "scripts.html", scriptsPageData)
	if err != nil {
		println("failed to execute tpl: scripts")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_footer.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_footer")
	}
}
