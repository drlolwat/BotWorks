package Helpers

import (
	"BotWorks/src/Config"
	"BotWorks/src/Storage"
	"database/sql"
	"fmt"
	"os/exec"
)

func StartBotFromId(userId int, pBotId int) bool {
	var config = Config.GetConfig()
	var botId int
	var botLoginName string
	var botLoginPass string
	var botBankPin string
	var proxyId int
	var scriptId int
	var scriptParams string
	var botBreaks string
	var botFps string
	var scriptDbId int
	var scriptSdnName string

	rows, err := Storage.GetDb().Query(
		"SELECT a.id,a.name,a.password,a.bank_pin,a.proxy_id,a.script_id,a.script_params,a.breaks,a.fps,"+
			"b.id,b.backend_name FROM accounts AS a JOIN scripts AS b ON (a.script_id=b.id) WHERE a.id = ? AND a.customer_id = ?", pBotId, userId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		err := rows.Scan(&botId, &botLoginName, &botLoginPass, &botBankPin, &proxyId, &scriptId, &scriptParams, &botBreaks, &botFps, &scriptDbId, &scriptSdnName)
		if err != nil {
			return false
		}

		var proxyUser string
		var proxyPass string
		var proxyIp string
		var proxyPort string

		args := []string{
			"-jar",
			"-Xms256M",
			"-Xmx4G",
			config.DreambotJar,
			"-fresh",
			"-covert",
			"-username",
			"\"" + config.DreambotUsername + "\"",
			"-password",
			"\"" + config.DreambotPassword + "\"",
			"-script",
			"\"" + scriptSdnName + "\"",
			"-accountUsername",
			"\"" + botLoginName + "\"",
			"-accountPassword",
			"\"" + botLoginPass + "\"",
			"-accountPin",
			"\"" + botBankPin + "\"",
		}

		if proxyId > 0 {
			prox, _ := Storage.GetDb().Query("SELECT ip,port,username,password FROM proxies WHERE id = ?", proxyId)

			defer func(prox *sql.Rows) {
				err := prox.Close()
				if err != nil {
					return
				}
			}(prox)

			for prox.Next() {
				_ = prox.Scan(&proxyIp, &proxyPort, &proxyUser, &proxyPass)
				args = append(args, "-proxyHost", "\""+proxyIp+"\"", "-proxyPort", "\""+proxyPort+"\"", "-proxyUser", "\""+proxyUser+"\"", "-proxyPass", "\""+proxyPass+"\"")
			}
		}
		args = append(args, "-destroy", "-params", scriptParams)
		fmt.Println(args)

		command := exec.Command("java", args...)
		if err := command.Start(); err != nil {
			fmt.Println(err.Error())
			return false
		}

		fmt.Println(args)

		return true
	}

	return false
}
