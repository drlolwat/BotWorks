package Helpers

import (
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"database/sql"
	"fmt"
)

func GetScripts() []Structs.ScriptData {
	rows, err := Storage.GetDb().Query("SELECT id,name,backend_name,auto_replenish FROM scripts WHERE id > 0")

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

	var allScripts []Structs.ScriptData

	for rows.Next() {
		var id int
		var name string
		var sdnName string
		var autoReplenish int

		rows.Scan(&id, &name, &sdnName, &autoReplenish)
		allScripts = append(allScripts, Structs.ScriptData{
			Id:            id,
			PublicName:    name,
			SdnName:       sdnName,
			AutoReplenish: autoReplenish,
		})
	}

	return allScripts
}

func GetScriptById(id int) Structs.ScriptData {
	rows, err := Storage.GetDb().Query("SELECT id,name,backend_name,auto_replenish FROM scripts WHERE id = ?", id)

	if err != nil {
		fmt.Println(err)
		return Structs.ScriptData{}
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var script Structs.ScriptData

	for rows.Next() {
		var id int
		var name string
		var sdnName string
		var autoReplenish int

		err := rows.Scan(&id, &name, &sdnName, &autoReplenish)
		if err != nil {
			fmt.Println(err)
			return Structs.ScriptData{}
		}

		script = Structs.ScriptData{
			Id:            id,
			PublicName:    name,
			SdnName:       sdnName,
			AutoReplenish: autoReplenish,
		}
	}

	return script
}
