package Helpers

import (
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"database/sql"
	"fmt"
)

func GetGroups(customerId int) []Structs.GroupData {
	botsOnline, botsOffline, _ := GetBotsByCustomerId(customerId)
	rows, err := Storage.GetDb().Query("SELECT id,name,complete_action,locked_action,banned_action FROM account_groups WHERE customer_id= ?", customerId)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var allGroups []Structs.GroupData

	for rows.Next() {
		var id int
		var name string
		var completeAction int
		var lockedAction int
		var bannedAction int

		err := rows.Scan(&id, &name, &completeAction, &lockedAction, &bannedAction)
		if err != nil {
			return nil
		}

		var myBots []Structs.DashboardBotData
		for _, bots := range botsOnline {
			if bots.GroupId == id {
				myBots = append(myBots, bots)
			}
		}
		for _, bots := range botsOffline {
			if bots.GroupId == id {
				myBots = append(myBots, bots)
			}
		}

		allGroups = append(allGroups, Structs.GroupData{
			Id:              id,
			Name:            name,
			BotsOnline:      botsOnline,
			BotsOffline:     botsOffline,
			BotsCount:       len(myBots),
			CompletedAction: completeAction,
			LockedAction:    lockedAction,
			BannedAction:    bannedAction,
		})
	}

	return allGroups
}

func GetGroupById(groupId int, customerId int) Structs.GroupData {
	var id int
	var name string
	var completeAction int
	var lockedAction int
	var bannedAction int
	err := Storage.GetDb().QueryRow("SELECT id,name,complete_action,locked_action,banned_action FROM account_groups WHERE id = ?", groupId).Scan(&id, &name, &completeAction, &lockedAction, &bannedAction)

	if err != nil {
		fmt.Println("GroupHelper:", err)
		return Structs.GroupData{Id: -1}
	}

	botsOnline, botsOffline, _ := GetBotsByCustomerId(customerId)

	return Structs.GroupData{
		Id:              id,
		Name:            name,
		BotsOnline:      botsOnline,
		BotsOffline:     botsOffline,
		BotsCount:       len(botsOffline) + len(botsOnline),
		CompletedAction: completeAction,
		LockedAction:    lockedAction,
		BannedAction:    bannedAction,
	}
}
