package Helpers

import (
	"BotWorks/src/Config"
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// cheers chatgpt
func GetLatestLogFile() string {
	config := Config.GetConfig()
	files, err := ReadDir(config.DreambotLogsDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return ""
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})
	logFile := files[0].Name()
	return logFile
}

func GetLogFileWithKeyword(keyword string) string {
	config := Config.GetConfig()
	files, err := ReadDir(config.DreambotLogsDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return ""
	}

	// Sort files by creation time in descending order
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})

	// Look at the last 10 files created
	for i := 0; i < 100 && i < len(files); i++ {
		if !files[i].IsDir() {
			logFile := files[i].Name()
			if strings.Contains(logFile, keyword) {
				return logFile
			} else {
				filePath := filepath.Join(config.DreambotLogsDir, logFile)
				file, err := os.Open(filePath)
				if err != nil {
					fmt.Println("Error opening file:", err)
					continue
				}
				defer file.Close()

				// Read the first 30 lines of the file
				scanner := bufio.NewScanner(file)
				for j := 0; scanner.Scan() && j < 100; j++ {
					if strings.Contains(scanner.Text(), keyword) {
						return logFile
					}
				}

				if err := scanner.Err(); err != nil {
					fmt.Println("Error scanning file:", err)
					continue
				}
			}
		}
	}

	return ""
}

func GetBotStatus(botId int) string {
	config := Config.GetConfig()
	var fileName string

	rows, _ := Storage.GetDb().Query("SELECT logfile FROM accounts WHERE id = ?", botId)

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(rows)

	for rows.Next() {
		err := rows.Scan(&fileName)
		if err != nil {
			return "No Status"
		}
	}

	filepath := config.DreambotLogsDir + "/" + fileName

	lastLine := TailFile(filepath, 4)
	lastLine = strings.ToLower(lastLine)

	if strings.Contains(lastLine, "starting break") {
		return "BREAK"
	}

	lastTen := TailFile(filepath, 60)
	lastTen = strings.ToLower(lastTen)

	if strings.Contains(lastTen, "members_world") {
		return "BOND"
	} else if strings.Contains(lastTen, "severity") {
		if strings.Contains(lastTen, "disabled") {
			return "BANNED"
		} else if strings.Contains(lastTen, "locked") {
			return "LOCKED"
		}
	} else if strings.Contains(lastTen, "reached target ttl and qp") || strings.Contains(lastTen, "reached non-99 target levels and qp") || strings.Contains(lastTen, "trade unrestricted, stopping") {
		return "COMPLETE"
	} else if strings.Contains(lastTen, "running: none") || strings.Contains(lastTen, "problem starting") {
		return "NOSCRIPT"
	}

	return "ACTIVE"
}

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

func StartFromDbRows(rows *sql.Rows, userId int, delay int) {
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

			botPid, logFile := StartBotFromId(userId, botId)
			err = UpdateBotStatus(botId, botPid, logFile)
			if err != nil {
				fmt.Println(err)
				return
			}

			delayTime := time.Duration(delay) * time.Minute

			if delay <= 0 {
				time.Sleep(3 * time.Second) // mass starts wait an additional 3 seconds to make sure they get the logfile
			} else {
				time.Sleep(delayTime)
			}
		}
	}()
}

// returns pid,logfile
func StartBotFromId(userId int, pBotId int) (int, string) {
	if !authenticated {
		fmt.Println("Starting of bots is deactivated until a valid license code is present.")
		return 0, ""
	}

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
	var botWorld string

	rows, err := Storage.GetDb().Query(
		"SELECT a.id,a.name,a.password,a.bank_pin,a.proxy_id,a.script_id,a.script_params,a.breaks,a.fps,a.world,"+
			"b.id,b.backend_name FROM accounts AS a JOIN scripts AS b ON (a.script_id=b.id) WHERE a.id = ? AND a.customer_id = ?", pBotId, userId)
	if err != nil {
		fmt.Println(err)
		return 0, ""
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(rows)

	for rows.Next() {
		err := rows.Scan(&botId, &botLoginName, &botLoginPass, &botBankPin, &proxyId, &scriptId, &scriptParams, &botBreaks, &botFps, &botWorld, &scriptDbId, &scriptSdnName)
		if err != nil {
			fmt.Println(err)
			return 0, ""
		}

		var proxyUser string
		var proxyPass string
		var proxyIp string
		var proxyPort string

		args := []string{
			"-jar",
			"-Xms" + config.JavaXms,
			"-Xmx" + config.JavaXmx,
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
			"-fps",
			botFps,
		}

		if proxyId > 0 {
			prox, _ := Storage.GetDb().Query("SELECT ip,port,username,password FROM proxies WHERE id = ?", proxyId)

			for prox.Next() {
				_ = prox.Scan(&proxyIp, &proxyPort, &proxyUser, &proxyPass)
				args = append(args, "-proxyHost", proxyIp, "-proxyPort", proxyPort, "-proxyUser", proxyUser, "-proxyPass", proxyPass)
			}

			err := prox.Close()
			if err != nil {
				return 0, ""
			}
		}

		args = append(args, "-destroy", "-fresh", "-covert")

		if strings.Contains(strings.ToLower(scriptSdnName), "f2p") {
			args = append(args, "-world", "f2p")
		} else if strings.Contains(strings.ToLower(scriptSdnName), "mule") {
			args = append(args, "-world", "308")
		} else if strings.Contains(strings.ToLower(scriptSdnName), "todt") {
			args = append(args, "-world", "307")
		} else if strings.Contains(strings.ToLower(scriptSdnName), "watminer") {
			args = append(args, "-world", "f2p")
		} else {
			args = append(args, "-world", botWorld)
		}

		if scriptParams != "" {
			args = append(args, "-params", scriptParams)
		}

		if botBreaks != "" {
			args = append(args, "-breaks", botBreaks)
		} else {
			args = append(args, "-breaks", "\"\"")
		}

		execType := "javaw"
		if OperatingSystem != "windows" {
			execType = "java"
		}

		//fmt.Println("started bot with:", args)

		command := exec.Command(execType, args...)
		err = command.Start()
		if err != nil {
			fmt.Println(err)
			return 0, ""
		}

		pid := command.Process.Pid
		time.Sleep(60 * time.Second)
		return pid, GetLogFileWithKeyword(botLoginName) //GetLatestLogFile()
	}

	return 0, ""
}

func UpdateBotStatus(botId int, pid int, logFile string) error {
	if !authenticated {
		return errors.New("status updating of bots is deactivated until a valid license code is present")
	}

	stmt, err := Storage.GetDb().Prepare("UPDATE `accounts` SET `pid` = ?, `logfile` = ? WHERE `id` = ?")
	if err != nil {
		fmt.Println(err)
		return err
	}

	res, err := stmt.Exec(pid, logFile, botId)

	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = res.RowsAffected()

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

func GetBotsByCustomerId(userId int) ([]Structs.DashboardBotData, []Structs.DashboardBotData, int) {
	var runningBotData []Structs.DashboardBotData
	var offlineBotData []Structs.DashboardBotData

	rows, err := Storage.GetDb().Query("SELECT a.id,a.name,a.proxy_id,a.script_id,a.pid,a.group_id,a.fps,b.name AS script_name FROM accounts AS a JOIN scripts AS b ON (a.script_id=b.id) WHERE a.customer_id = ?", userId)
	if err != nil {
		fmt.Println(err)
		return []Structs.DashboardBotData{}, []Structs.DashboardBotData{}, 0
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var runningBots = 0
	var stoppedBots = 0
	var scripts []int

	for rows.Next() {
		var accountId int
		var name string
		var proxyId int
		var scriptId int
		var pid int
		var groupId int
		var scriptName string
		var status string
		var fps int

		err = rows.Scan(&accountId, &name, &proxyId, &scriptId, &pid, &groupId, &fps, &scriptName)
		if err != nil {
			fmt.Println(err)
			return []Structs.DashboardBotData{}, []Structs.DashboardBotData{}, 0
		}

		if pid > 0 {
			if !ContainsInt(scriptId, scripts) {
				scripts = append(scripts, scriptId)
			}
			runningBots++
		} else {
			stoppedBots++
		}

		//TODO status comes from the logfile
		status = GetBotStatus(accountId)

		var botData = Structs.DashboardBotData{
			AccountId:   accountId,
			AccountName: name,
			Script:      scriptName,
			Status:      status,
			Pid:         pid,
			GroupId:     groupId,
			ProxyId:     proxyId,
			Fps:         fps,
		}

		if pid > 0 {
			if !ContainsBotData(botData, runningBotData) {
				runningBotData = append(runningBotData, botData)
			}
		} else {
			if !ContainsBotData(botData, offlineBotData) {
				offlineBotData = append(offlineBotData, botData)
			}
		}

	}

	return runningBotData, offlineBotData, len(scripts)
}
