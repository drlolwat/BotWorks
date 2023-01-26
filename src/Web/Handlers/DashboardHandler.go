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
	rows, err := Storage.GetDb().Query("SELECT a.id,a.name,a.proxy_id,a.script_id,a.pid,a.group_id,b.name AS script_name FROM accounts AS a JOIN scripts AS b ON (a.script_id=b.id) WHERE a.customer_id = ?", user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var runningBots = 0
	var stoppedBots = 0
	var uniqueScripts = 0
	var scripts []int
	var runningBotData []Structs.DashboardBotData
	var offlineBotData []Structs.DashboardBotData

	for rows.Next() {
		var accountId int
		var name string
		var proxyId int
		var scriptId int
		var pid int
		var groupId int
		var scriptName string
		var status string

		err = rows.Scan(&accountId, &name, &proxyId, &scriptId, &pid, &groupId, &scriptName)
		if err != nil {
			fmt.Println(err)
			return
		}

		if !Helpers.ContainsInt(scriptId, scripts) {
			scripts = append(scripts, scriptId)
		}

		if pid > 0 {
			runningBots++
		} else {
			stoppedBots++
		}

		//TODO status comes from the logfile
		status = "No Status"

		var botData = Structs.DashboardBotData{
			AccountId:   accountId,
			AccountName: name,
			Script:      scriptName,
			Status:      status,
			Pid:         pid,
			GroupId:     groupId,
		}

		if pid > 0 {
			if !Helpers.ContainsDashboardBotData(botData, runningBotData) {
				runningBotData = append(runningBotData, botData)
			}
		} else {
			if !Helpers.ContainsDashboardBotData(botData, offlineBotData) {
				offlineBotData = append(offlineBotData, botData)
			}
		}
	}

	uniqueScripts = len(scripts)
	var groups []Structs.GroupData

	// fetch my groups
	rows, err = Storage.GetDb().Query("SELECT id,name FROM account_groups WHERE customer_id = ?", user.ID)
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
				var localRunningBotData = Structs.DashboardBotData{
					AccountId:   bots.AccountId,
					AccountName: bots.AccountName,
					Script:      bots.Script,
					Status:      bots.Status,
					Pid:         bots.Pid,
					GroupId:     bots.GroupId,
				}

				groupRunningBotData = append(groupRunningBotData, localRunningBotData)
			}
		}

		for _, bots := range offlineBotData {
			if bots.GroupId == groupId {
				var localOfflineBotData = Structs.DashboardBotData{
					AccountId:   bots.AccountId,
					AccountName: bots.AccountName,
					Script:      bots.Script,
					Status:      bots.Status,
					Pid:         bots.Pid,
					GroupId:     bots.GroupId,
				}

				groupOfflineBotData = append(groupOfflineBotData, localOfflineBotData)
			}
		}

		groupData = Structs.GroupData{
			Id:         groupId,
			Name:       groupName,
			BotsOnline: groupRunningBotData,
		}

		if !Helpers.ContainsGroupData(groupData, groups) {
			groups = append(groups, groupData)
		}
	}

	err = Template.GetTpl().ExecuteTemplate(w, "bootstrap_header.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	data := Structs.DashboardData{RunningBots: runningBots, StoppedBots: stoppedBots, UniqueScripts: uniqueScripts, CurrentTime: Helpers.TimeString(), OnlineBotData: runningBotData, OfflineBotData: offlineBotData, GroupData: groups}
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
