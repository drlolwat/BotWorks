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

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if !Helpers.IsSessionValid(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var user Structs.UserData
	session, _ := Sessions.Store.Get(r, Sessions.SessionName)
	userID := session.Values[Sessions.UserIDKey]
	if userID != nil {
		user.ID = userID.(int)
	} else {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	var headerInfo = Structs.PageData{
		LoggedIn: user.ID > 0,
		Name:     "Proxies",
		Url:      "proxies",
	}

	if r.PostForm.Has("update") || r.PostForm.Has("add") {
		if r.FormValue("nickname") != "" && r.FormValue("ip") != "" && r.FormValue("port") != "" &&
			r.FormValue("username") != "" && r.FormValue("password") != "" {

			nickname := r.FormValue("nickname")
			ip := r.FormValue("ip")
			port := r.FormValue("port")
			username := r.FormValue("username")
			password := r.FormValue("password")

			if r.PostForm.Has("update") && r.FormValue("proxy_id") != "" {
				proxyid := r.FormValue("proxy_id")
				_, err := Storage.GetDb().Exec("UPDATE proxies SET ip = ?, port = ?, username = ?, password = ?, nickname = ? WHERE id = ? AND customer_id = ?", ip, port, username, password, nickname, proxyid, user.ID)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			if r.PostForm.Has("add") {
				_, err := Storage.GetDb().Exec("INSERT INTO proxies (ip, port, username, password, nickname, customer_id) VALUES (?, ?, ?, ?, ?, ?)", ip, port, username, password, nickname, user.ID)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}

	if r.URL.Query().Has("d") {
		_, err := Storage.GetDb().Exec("DELETE FROM proxies WHERE id = ? AND customer_id = ?", r.URL.Query().Get("d"), user.ID)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	proxyPageData := map[string]interface{}{
		"Proxies":       Helpers.GetProxies(user.ID),
		"AddingProxy":   r.URL.Query().Get("a"),
		"EditingProxy":  r.URL.Query().Get("e"),
		"DeletingProxy": r.URL.Query().Get("d"),
		"Scripts":       Helpers.GetScripts(),
	}

	if r.URL.Query().Has("e") {
		rows, err := Storage.GetDb().Query("SELECT id,ip,port,username,password,nickname FROM proxies WHERE id = ? AND customer_id = ?", r.URL.Query().Get("e"), user.ID)
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
			var ip string
			var port string
			var username string
			var password string
			var nickname string

			err = rows.Scan(&id, &ip, &port, &username, &password, &nickname)
			if err != nil {
				fmt.Println(err)
				return
			}

			proxyPageData["ProxyData"] = Structs.ProxyData{
				Id:        id,
				Name:      nickname,
				IpAddress: ip,
				Port:      port,
				Username:  username,
				Password:  password,
				Accounts:  0,
			}
		}

	} else if r.URL.Query().Has("a") {
		proxyPageData["ProxyData"] = Structs.ProxyData{
			Id:        0,
			Name:      "New Proxy",
			IpAddress: "127.0.0.1",
			Port:      "1234",
			Username:  "username",
			Password:  "password",
			Accounts:  0,
		}
	}

	err := Template.GetTpl().ExecuteTemplate(w, "bootstrap_header.html", headerInfo)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "proxies.html", proxyPageData)
	if err != nil {
		println("failed to execute tpl: accounts")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_footer.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_footer")
	}
}
