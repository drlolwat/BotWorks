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

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
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

	// fetch my bots
	runningBotData, offlineBotData, uniqueScripts := Helpers.GetBotsByCustomerId(user.ID)

	var allBots []Structs.DashboardBotData
	var groups []Structs.GroupData

	allBots = append(allBots, runningBotData...)
	allBots = append(allBots, offlineBotData...)

	// fetch my groups
	rows, err := Storage.GetDb().Query("SELECT id,name FROM account_groups WHERE customer_id = ?", user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var groupId int
		var groupName string
		var groupData = Structs.GroupData{}
		var groupRunningBotData []Structs.DashboardBotData
		var groupOfflineBotData []Structs.DashboardBotData

		err = rows.Scan(&groupId, &groupName)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, bots := range runningBotData {
			if bots.GroupId == groupId {
				groupRunningBotData = append(groupRunningBotData, bots)
			}
		}

		for _, bots := range offlineBotData {
			if bots.GroupId == groupId {
				groupOfflineBotData = append(groupOfflineBotData, bots)
			}
		}

		groupData = Structs.GroupData{
			Id:          groupId,
			Name:        groupName,
			BotsOnline:  groupRunningBotData,
			BotsOffline: groupOfflineBotData,
		}

		if !Helpers.ContainsGroupData(groupData, groups) {
			groups = append(groups, groupData)
		}
	}

	var headerInfo = Structs.PageData{
		LoggedIn: user.ID > 0,
		Name:     "Dashboard",
		Url:      "dashboard",
	}

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_header.html", headerInfo)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	data := Structs.DashboardData{RunningBots: len(runningBotData), StoppedBots: len(offlineBotData), UniqueScripts: uniqueScripts, CurrentTime: Helpers.TimeString(), OnlineBotData: runningBotData, OfflineBotData: offlineBotData, GroupData: groups}
	err = Template.GetTpl().ExecuteTemplate(w, "dashboard.html", data)
	if err != nil {
		println(err.Error())
		return
	}

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_footer.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_footer")
	}
}
