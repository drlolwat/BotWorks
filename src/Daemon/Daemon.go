package Daemon

import (
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
			currentTime := Helpers.TimeString()
			splitTime := strings.Split(currentTime, ":")
			currentHour, _ := strconv.Atoi(splitTime[0])
			currentMin, _ := strconv.Atoi(splitTime[1])

			botsUndetected := 0
			botsClosed := 0
			botsStarted := 0
			botsDetected := 0
			botsRestarted := 0

			rows, err := Storage.GetDb().Query("SELECT id,name,schedule,pid,customer_id,logfile,script_id,group_id,breaks,script_params,proxy_id FROM accounts")
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
				var breaks string
				var params string
				var proxy int

				err := rows.Scan(&id, &name, &schedule, &pid, &customer, &logfile, &scriptId, &groupId, &breaks, &params, &proxy)
				if err != nil {
					continue
				}

				if pid > 0 {
					//STATS-Attack:69 Defence:68 Strength:70 Hitpoints:74 Ranged:70 Prayer:44 Magic:70 Cooking:75 Woodcutting:75 Fletching:1 Fishing:71 Firemaking:67 Crafting:72 Smithing:56 Mining:66 Herblore:1 Agility:1 Thieving:1 Slayer:1 Farming:1 Runecrafting:45 Hunter:1 Construction:1 TOTAL:1000 TOTAL_THIS_SESSION:0 CMB:86 QP:16 QP_THIS_SESSION:0 GP:-1 GP_THIS_SESSION:0 LOGIN_TIME:1291616101 LOGIN_TIME_THIS_SESSION:1981 RESTRICTED:0 BANNED:0
					botStatus := Helpers.GetBotStatus(id)
					if !Helpers.IsBotRunning(pid) {
						_, err = Storage.GetDb().Exec("UPDATE accounts SET pid = 0 WHERE pid = ?", pid)
						if err != nil {
							continue
						}
						botsUndetected++
						continue
					} else {
						if len(logfile) == 0 || logfile == "" {
							logfile = Helpers.GetLogFileWithKeyword(name)
							if logfile != "" {
								fmt.Println("Daemon:", name, "using corrected logfile")
								_, err = Storage.GetDb().Exec("UPDATE accounts SET logfile = ? WHERE id = ?", logfile, id)
							}
						}

						if groupId > 0 {
							groupData := Helpers.GetGroupById(groupId, customer)
							if botStatus != "BANNED" {
								if botStatus == "COMPLETE" {
									if groupData.CompletedAction != 0 {
										Helpers.KillBot(id, customer)
										if groupData.CompletedAction == -1 {
											_, err = Storage.GetDb().Exec("UPDATE accounts SET schedule = 0 WHERE id = ?", id)
											if err != nil {
												fmt.Println("Daemon:", err)
											}
										} else {
											if groupData.CompletedAction > 0 {
												_, err = Storage.GetDb().Exec("UPDATE accounts SET schedule = 0, proxy_id = 0, group_id = ? WHERE id = ?", groupData.CompletedAction, id)
											}
										}

										if Helpers.Replenish(scriptId, groupId, schedule, customer, params, proxy) {
											botsStarted++
										}

										fmt.Println("Daemon:", name, "(", pid, ") has been tagged as completed")
										botsClosed++
									}
								} else if botStatus == "LOCKED" {
									if groupData.LockedAction != 0 {
										Helpers.KillBot(id, customer)
										if groupData.CompletedAction == -1 {
											_, err = Storage.GetDb().Exec("UPDATE accounts SET schedule = 0 WHERE id = ?", id)
											if err != nil {
												fmt.Println("Daemon:", err)
											}
										} else {
											if groupData.CompletedAction > 0 {
												_, err = Storage.GetDb().Exec("UPDATE accounts SET schedule = 0, proxy_id = 0, group_id = ? WHERE id = ?", groupData.LockedAction, id)
											}
										}

										if Helpers.Replenish(scriptId, groupId, schedule, customer, params, proxy) {
											botsStarted++
										}

										fmt.Println("Daemon:", name, "(", pid, ") has been tagged as locked")
										botsClosed++
									}
								} else {
									// restart any bots not running their intended script
									if botStatus == "NOSCRIPT" {
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
							} else {
								if groupData.BannedAction != 0 {
									Helpers.KillBot(id, customer)
									if groupData.BannedAction == -1 {
										_, err = Storage.GetDb().Exec("DELETE FROM accounts WHERE id = ?", id)
										if err != nil {
											fmt.Println("Daemon:", err)
										}

									} else if groupData.BannedAction == -2 {
										// something but i forgot
									} else if groupData.BannedAction > 0 {
										_, err = Storage.GetDb().Exec("UPDATE accounts SET schedule = 0, proxy_id = 0, group_id = ? WHERE id = ?", groupData.BannedAction, id)
										if err != nil {
											fmt.Println("Daemon:", err)
										}
									}

									if proxy > 0 {
										_, err = Storage.GetDb().Exec("UPDATE proxies SET bans = bans + 1 WHERE id = ?", proxy)
									}

									fmt.Println("Daemon:", name, "(", pid, ") has been tagged as banned")
									botsClosed++
								}

								if Helpers.Replenish(scriptId, groupId, schedule, customer, params, proxy) {
									botsStarted++
								}
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
					botStatus := Helpers.GetBotStatus(id)
					if currentHour == endHour && currentMin >= endMin {
						if botStatus == "BREAK" || currentMin+1 >= 60 || len(breaks) <= 2 { //TODO store a queue of bots to close and how often to recheck
							botsClosed++
							Helpers.KillBot(id, customer)
						}
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

			fmt.Println("Daemon:", botsDetected-botsClosed, "detected |", botsRestarted, "restarted |", botsStarted, "started |", botsClosed, "closed |", botsUndetected, "undetected")
			time.Sleep(cycleTime)
		} else {
			fmt.Println("Daemon is deactivated until a valid license code is present.")
			time.Sleep(cycleTime)
		}
	}
}
