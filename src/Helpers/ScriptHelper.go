package Helpers

import (
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"database/sql"
	"fmt"
)

func GetScripts() []Structs.ScriptData {
	rows, err := Storage.GetDb().Query("SELECT id,name,backend_name FROM scripts WHERE id > 0")

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

		rows.Scan(&id, &name, &sdnName)
		allScripts = append(allScripts, Structs.ScriptData{
			Id:         id,
			PublicName: name,
			SdnName:    sdnName,
		})
	}

	return allScripts
}
