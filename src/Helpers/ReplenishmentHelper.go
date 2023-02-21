package Helpers

import (
	"BotWorks/src/Storage"
	"database/sql"
	"fmt"
)

func Replenish(scriptId int, groupId int, schedule int, customer int, params string, proxy int) bool {
	scriptData := GetScriptById(scriptId)
	if scriptData.Id > 0 && scriptData.AutoReplenish > 0 {
		var newId int
		err := Storage.GetDb().QueryRow("SELECT id FROM accounts WHERE group_id = ? LIMIT 1", scriptData.AutoReplenish).Scan(&newId)

		if err != nil || err == sql.ErrNoRows {
			fmt.Println("Daemon: No bots are in auto-replenishment group for script:", scriptData.PublicName)
			return false
		}

		_, err = Storage.GetDb().Exec("UPDATE accounts SET group_id = ?, script_id = ?, schedule = ?, script_params = ?, proxy_id = ? WHERE pid = 0 AND group_id = ? AND customer_id = ?", groupId, scriptId, schedule, params, proxy, scriptData.AutoReplenish, customer)

		if err != nil {
			fmt.Println("Daemon:", err)
			return false
		}

		newPid, newLog := StartBotFromId(customer, newId)
		err = UpdateBotStatus(newId, newPid, newLog)
		if err != nil {
			fmt.Println("Daemon:", "Big error issue restarting bot, could not update status!")
			KillBot(newId, customer)
		}

		return true
	}

	return false
}
