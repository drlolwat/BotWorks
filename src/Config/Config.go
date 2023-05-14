package Config

import (
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"database/sql"
	"fmt"
)

func GetConfig() Structs.Configuration {
	var ret Structs.Configuration

	rows, err := Storage.GetDb().Query("SELECT db_jar,db_logs,db_username,db_password,java_xms,java_xmx,port FROM configuration")
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
	var javaXms string
	var javaXmx string
	var bbPort int

	for rows.Next() {
		err = rows.Scan(&dbJar, &dbLogs, &dbUser, &dbPass, &javaXms, &javaXmx, &bbPort)
		ret = Structs.Configuration{
			DreambotJar:      dbJar,
			DreambotLogsDir:  dbLogs,
			DreambotUsername: dbUser,
			DreambotPassword: dbPass,
			JavaXms:          javaXms,
			JavaXmx:          javaXmx,
			BBPort:           bbPort,
		}
	}

	return ret
}
