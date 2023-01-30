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

			rows, err := Storage.GetDb().Query("SELECT id,name,schedule,pid,customer_id FROM accounts")
			if err != nil {
				continue
			}

			for rows.Next() {
				var id int
				var name string
				var schedule int
				var pid int
				var customer int

				err := rows.Scan(&id, &name, &schedule, &pid, &customer)
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

			fmt.Println("Daemon has detected", botsDetected, "bots. Started:", botsStarted, "| Stopped:", botsClosed, "| Undetected:", botsUndetected)
			time.Sleep(cycleTime)
		} else {
			fmt.Println("Daemon is deactivated until a valid license code is present.")
			time.Sleep(cycleTime)
		}
	}
}
