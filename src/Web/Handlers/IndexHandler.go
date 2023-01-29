package Handlers

import (
	"BotWorks/src/Helpers"
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"BotWorks/src/Web/Template"
	"github.com/niktheblak/passwordhash/pkg/scrypt"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	if Helpers.IsSessionValid(r) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}

	if r.PostForm.Has("handle") && r.PostForm.Has("password") {
		username := r.FormValue("handle")
		password := r.FormValue("password")

		if len(username) > 0 && len(password) > 0 {
			var user Structs.UserData
			if err := Storage.GetDb().QueryRow("SELECT id, handle, password FROM customers WHERE handle = ?", username).Scan(&user.Id, &user.Handle, &user.Password); err != nil {
				println("login failed for user: " + username)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			// Compare the provided password with the hashed password from the database
			if err := scrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
				println("login failed for user: " + username)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			session, _ := Sessions.Store.Get(r, Sessions.SessionName)
			session.Values[Sessions.UserIDKey] = user.Id
			session.Save(r, w)

			println("User has authenticated: " + username)
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}
	}

	var headerInfo = Structs.PageData{
		LoggedIn: Helpers.IsSessionValid(r),
		Name:     "Login",
		Url:      "login",
	}

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_header.html", headerInfo)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		println("failed to execute tpl: index")
	}
}
