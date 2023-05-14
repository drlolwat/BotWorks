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

func AssignProxies(customer int, group int) {
	rows, err := Storage.GetDb().Query("SELECT id FROM accounts WHERE group_id = ? AND proxy_id <= 0", group)

	if err != nil {
		fmt.Println(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	for rows.Next() {
		var id int

		err = rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
		}

		newProxy := GetRandomProxy(customer, 0, true)
		if newProxy > 0 {
			_, err = Storage.GetDb().Exec("UPDATE accounts SET proxy_id = ? WHERE id = ?", newProxy, id)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func RemoveProxies(group int) {
	_, err := Storage.GetDb().Exec("UPDATE accounts SET proxy_id = -1 WHERE group_id = ?", group)
	if err != nil {
		fmt.Println(err)
	}
}

func ExportAccounts(group int) string {
	rows, err := Storage.GetDb().Query("SELECT name,password FROM accounts WHERE group_id = ?", group)
	if err != nil {
		fmt.Println(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	builder := ""

	for rows.Next() {
		var name string
		var pass string

		err = rows.Scan(&name, &pass)
		if err != nil {
			fmt.Println(err)
		}

		builder = builder + "\r\n" + name + ":" + pass
	}

	return builder
}

func DeleteAccounts(group int) {
	_, err := Storage.GetDb().Exec("DELETE FROM accounts WHERE group_id = ?", group)
	if err != nil {
		fmt.Println(err)
	}
}
