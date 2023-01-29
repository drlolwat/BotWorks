package Helpers

import (
	"BotWorks/src/Config"
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"database/sql"
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
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

	lastTen := TailFile(filepath, 10)
	lastTen = strings.ToLower(lastTen)

	if strings.Contains(lastTen, "members_world") {
		return "BOND"
	} else if strings.Contains(lastTen, "disabled") {
		if strings.Contains(lastTen, "severity") {
			return "BANNED"
		}
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

func KillBot(botId int, userId int) {
	pid := GetPid(botId)

	if pid > 0 {
		cmd := exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/F")
		_ = cmd.Start()

		err := SetBotPidById(botId, 0, "")
		if err != nil {
			fmt.Println(err)
			return
		}
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

			botPid, logFile := StartBotFromId(userId, botId)
			err = SetBotPidById(botId, botPid, logFile)
			if err != nil {
				fmt.Println(err)
				return
			}
			time.Sleep(3 * time.Second) // mass starts wait an additional 3 seconds to make sure they get the logfile
		}
	}()
}

func StartBotFromId(userId int, pBotId int) (int, string) {
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
		err := rows.Scan(&botId, &botLoginName, &botLoginPass, &botBankPin, &proxyId, &scriptId, &scriptParams, &botBreaks, &botFps, &scriptDbId, &scriptSdnName)
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
			"-fps",
			botFps,
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
		time.Sleep(3 * time.Second)
		return pid, GetLatestLogFile()
	}

	return 0, ""
}

func IsBotRunning(pid int) bool {
	var command string
	var output []byte
	var err error

	if runtime.GOOS == "windows" {
		command = fmt.Sprintf("tasklist /FI \"PID eq %d\"", pid)
		output, err = exec.Command("cmd", "/C", command).Output()
	} else {
		command = fmt.Sprintf("ps %d", pid)
		output, err = exec.Command("sh", "-c", command).Output()
	}

	if err != nil {
		return false
	}

	return !strings.Contains(strings.ToLower(string(output)), "no tasks")
}

func SetBotPidById(botId int, pid int, logFile string) error {
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
