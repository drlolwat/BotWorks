package Helpers

import (
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"database/sql"
	"fmt"
)

func GetSchedules(customerId int) []Structs.ScheduleData {
	rows, err := Storage.GetDb().Query("SELECT id,name,start_hour,start_min,end_hour,end_min FROM account_schedules WHERE customer_id= ?", customerId)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var allSchedules []Structs.ScheduleData

	for rows.Next() {
		var id int
		var name string
		var starth int
		var startm int
		var endh int
		var endm int

		err := rows.Scan(&id, &name, &starth, &startm, &endh, &endm)
		if err != nil {
			return nil
		}

		allSchedules = append(allSchedules, Structs.ScheduleData{
			Id:          id,
			Name:        name,
			StartHour:   starth,
			StartMinute: startm,
			EndHour:     endh,
			EndMinute:   endm,
		})
	}

	return allSchedules
}
