package Config

import (
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"database/sql"
	"fmt"
)

func GetConfig() Structs.Configuration {
	var ret Structs.Configuration

	rows, err := Storage.GetDb().Query("SELECT db_jar,db_logs,db_username,db_password,banned_action,complete_action FROM configuration")
	if err != nil {
		fmt.Println(err)
		return ret
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var dbJar string
	var dbLogs string
	var dbUser string
	var dbPass string
	var bannedAction int
	var completeAction int

	for rows.Next() {
		err = rows.Scan(&dbJar, &dbLogs, &dbUser, &dbPass, &bannedAction)
		ret = Structs.Configuration{
			DreambotJar:      dbJar,
			DreambotLogsDir:  dbLogs,
			DreambotUsername: dbUser,
			DreambotPassword: dbPass,
			BannedAction:     bannedAction,
			CompleteAction:   completeAction,
		}
	}

	return ret
}
