package Daemon

import (
	"BotWorks/src/Config"
	"BotWorks/src/Helpers"
	"BotWorks/src/Storage"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Begin() {
	cycleTime := 60 * time.Second
	time.Sleep(3 * time.Second) // startup wait

	for true {
		if Helpers.License() {
			config := Config.GetConfig()
			currentTime := Helpers.TimeString()
			splitTime := strings.Split(currentTime, ":")
			currentHour, _ := strconv.Atoi(splitTime[0])
			currentMin, _ := strconv.Atoi(splitTime[1])

			botsUndetected := 0
			botsClosed := 0
			botsStarted := 0
			botsDetected := 0
			botsRestarted := 0

			rows, err := Storage.GetDb().Query("SELECT id,name,schedule,pid,customer_id,logfile,script_id,group_id FROM accounts")
			if err != nil {
				continue
			}

			for rows.Next() {
				var id int
				var name string
				var schedule int
				var pid int
				var customer int
				var logfile string
				var scriptId int
				var groupId int

				err := rows.Scan(&id, &name, &schedule, &pid, &customer, &logfile, &scriptId, &groupId)
				if err != nil {
					continue
				}

				// check if any bots need to be marked as offline
				if pid > 0 {
					if !Helpers.IsBotRunning(pid) {
						_, err = Storage.GetDb().Exec("UPDATE accounts SET pid = 0 WHERE pid = ?", pid)
						if err != nil {
							continue
						}
						botsUndetected++
						continue
					} else {
						botStatus := Helpers.GetBotStatus(id)
						if botStatus != "BANNED" {
							if botStatus == "COMPLETE" {
								if config.CompleteAction != 0 {
									Helpers.KillBot(id, customer)
									if config.CompleteAction == -1 {
										_, err = Storage.GetDb().Exec("UPDATE accounts SET schedule = 0 WHERE id = ?", id)
										if err != nil {
											fmt.Println("Daemon:", err)
										}
									} else {
										if config.CompleteAction > 0 {
											_, err = Storage.GetDb().Exec("UPDATE accounts SET schedule = 0, group_id = ? WHERE id = ?", config.CompleteAction, id)
										}
									}
									fmt.Println("Daemon:", name, "has been tagged as completed")
								}
							} else {
								// restart any bots not running their intended script
								file := config.DreambotLogsDir + "/" + logfile
								fileData := Helpers.TailFile(file, 30)
								if strings.Contains(fileData, "[SYSTEM INFO]") {
									// Metrics here.
									// Restart the bot if its not running a script
									if strings.Contains(fileData, "Script Running: None") {
										Helpers.KillBot(id, customer)
										time.Sleep(650 * time.Millisecond)
										newPid, newLog := Helpers.StartBotFromId(customer, id)
										err := Helpers.UpdateBotStatus(id, newPid, newLog)
										if err != nil {
											fmt.Println("Daemon:", "Big error issue restarting bot, could not update status!")
											Helpers.KillBot(id, customer)
											continue
										}
										botsRestarted++
									}
								}
							}
						} else {
							if config.BannedAction != 0 {
								Helpers.KillBot(id, customer)
								if config.BannedAction == -1 {
									_, err = Storage.GetDb().Exec("DELETE FROM accounts WHERE id = ?", id)
									if err != nil {
										fmt.Println("Daemon:", err)
									}

								} else if config.BannedAction == -2 {
									// something but i forgot
								} else if config.BannedAction > 0 {
									_, err = Storage.GetDb().Exec("UPDATE accounts SET schedule = 0, proxy_id = 0, group_id = ? WHERE id = ?", config.BannedAction, id)
									if err != nil {
										fmt.Println("Daemon:", err)
									}
								}
								fmt.Println("Daemon:", name, "has been tagged as banned")
								botsClosed++
							}

							scriptData := Helpers.GetScriptById(scriptId)
							if scriptData.Id > 0 && scriptData.AutoReplenish > 0 {
								var newId int
								err := Storage.GetDb().QueryRow("UPDATE accounts SET group_id = ?, script_id = ?, schedule = ? WHERE pid = 0 AND group_id = ? AND customer_id = ? LIMIT 1 RETURNING id",
									groupId, scriptId, schedule, scriptData.AutoReplenish, customer).Scan(&newId)

								if err != nil {
									fmt.Println("Daemon:", err)
								}

								newPid, newLog := Helpers.StartBotFromId(customer, newId)

								err = Helpers.UpdateBotStatus(newId, newPid, newLog)
								if err != nil {
									fmt.Println("Daemon:", "Big error issue restarting bot, could not update status!")
									Helpers.KillBot(id, customer)
								}

								botsStarted++
							}
						}

						botsDetected++
					}
				}

				// check for scheduling
				var scheduleName string
				var startHour int
				var startMin int
				var endHour int
				var endMin int

				err = Storage.GetDb().QueryRow("SELECT name,start_hour,start_min,end_hour,end_min FROM account_schedules WHERE id = ?", schedule).Scan(&scheduleName, &startHour, &startMin, &endHour, &endMin)
				if err != nil {
					continue
				}

				if pid > 0 {
					if currentHour == endHour && currentMin >= endMin {
						botsClosed++
						Helpers.KillBot(id, customer)
						continue
					}
				} else {
					if currentHour == startHour && currentMin >= startMin {
						botsStarted++
						newPid, newLog := Helpers.StartBotFromId(customer, id)
						err := Helpers.UpdateBotStatus(id, newPid, newLog)
						if err != nil {
							fmt.Println("Big error issue starting bot, could not update status!")
							Helpers.KillBot(id, customer)
							continue
						}
						time.Sleep(5 * time.Second)
						continue
					}
				}
			}

			_ = rows.Close()

			fmt.Println("Daemon has detected", botsDetected, "bots. Restarted:", botsRestarted, "| Started:", botsStarted, "| Stopped:", botsClosed, "| Undetected:", botsUndetected)
			time.Sleep(cycleTime)
		} else {
			fmt.Println("Daemon is deactivated until a valid license code is present.")
			time.Sleep(cycleTime)
		}
	}
}
