package Helpers

import (
	"BotWorks/src/Config"
	"BotWorks/src/Storage"
	"database/sql"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

func GetPid(botId int) int {
	var pid int

	rows, _ := Storage.GetDb().Query("SELECT pid FROM accounts WHERE id = ?", botId)

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	for rows.Next() {
		err := rows.Scan(&pid)
		if err != nil {
			return 0
		}
	}

	return pid
}

func KillBot(botId int, userId int) {
	pid := GetPid(botId)

	if pid > 0 {
		cmd := exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/F")
		_ = cmd.Start()

		SetBotPidById(botId, 0)
	}

	return
}

func StartFromDbRows(rows *sql.Rows, userId int) {
	go func() {
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {

			}
		}(rows)

		for rows.Next() {
			var botId int
			err := rows.Scan(&botId)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = SetBotPidById(botId, StartBotFromId(userId, botId))
			if err != nil {
				fmt.Println(err)
				return
			}
			time.Sleep(1)
		}
	}()
}

func StartBotFromId(userId int, pBotId int) int {
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
		return 0
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(rows)

	for rows.Next() {
		err := rows.Scan(&botId, &botLoginName, &botLoginPass, &botBankPin, &proxyId, &scriptId, &scriptParams, &botBreaks, &botFps, &scriptDbId, &scriptSdnName)
		if err != nil {
			fmt.Println(err)
			return 0
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
			"-username",
			config.DreambotUsername,
			"-password",
			config.DreambotPassword,
			"-script",
			scriptSdnName,
			"-accountUsername",
			botLoginName,
			"-accountPassword",
			botLoginPass,
		}

		if proxyId > 0 {
			prox, _ := Storage.GetDb().Query("SELECT ip,port,username,password FROM proxies WHERE id = ?", proxyId)

			defer func(prox *sql.Rows) {
				err := prox.Close()
				if err != nil {
					fmt.Println(err)
				}
			}(prox)

			for prox.Next() {
				_ = prox.Scan(&proxyIp, &proxyPort, &proxyUser, &proxyPass)
				args = append(args, "-proxyHost", proxyIp, "-proxyPort", proxyPort, "-proxyUser", proxyUser, "-proxyPass", proxyPass)
			}
		}

		args = append(args, "-destroy", "-fresh", "-covert", scriptParams)
		command := exec.Command("C:\\Program Files\\Java\\jre1.8.0_361\\bin\\javaw.exe", args...)
		err = command.Start()
		if err != nil {
			fmt.Println(err)
		}

		pid := command.Process.Pid
		return pid
	}

	return 0
}

func SetBotPidById(botId int, pid int) error {
	stmt, err := Storage.GetDb().Prepare("UPDATE `accounts` SET `pid` = ? WHERE `id` = ?")
	if err != nil {
		fmt.Println(err)
		return err
	}

	res, err := stmt.Exec(pid, botId)
	rows, err := res.RowsAffected()

	if rows > 0 {
		fmt.Println("updating acc", botId, "to pid", pid)
	}

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = stmt.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
