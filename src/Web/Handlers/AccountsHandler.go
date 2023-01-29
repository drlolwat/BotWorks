package Handlers

import (
	"BotWorks/src/Helpers"
	"BotWorks/src/Structs"
	"BotWorks/src/Web/Sessions"
	"BotWorks/src/Web/Template"
	"net/http"
)

func AccountsHandler(w http.ResponseWriter, r *http.Request) {
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
		Name:     "Accounts",
		Url:      "accounts",
	}

	err := Template.GetTpl().ExecuteTemplate(w, "bootstrap_header.html", headerInfo)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	err = r.ParseForm()
	if err != nil {
		return
	}

	runningBotData, offlineBotData, _ := Helpers.GetBotsByCustomerId(user.ID)

	var allBots []Structs.BotData
	allBots = append(allBots, runningBotData...)
	allBots = append(allBots, offlineBotData...)

	accountPageData := map[string]interface{}{
		"Bots":           allBots,
		"AddingAccount":  r.URL.Query().Get("a"),
		"EditingAccount": r.URL.Query().Get("e"),
		"DeletingBot":    r.URL.Query().Get("d"),
		"Scripts":        Helpers.GetScripts(),
		"Proxies":        Helpers.GetProxies(user.ID),
		"Schedules":      Helpers.GetSchedules(user.ID),
		"Groups":         Helpers.GetGroups(user.ID),
	}

	if r.PostForm.Has("update") || r.PostForm.Has("add") {
		if r.PostForm.Has("update") {
			// do update
		} else if r.PostForm.Has("add") {
			// add
		}
	}

	err = Template.GetTpl().ExecuteTemplate(w, "accounts.html", accountPageData)
	if err != nil {
		println("failed to execute tpl: accounts")
	}

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_footer.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_footer")
	}

}
