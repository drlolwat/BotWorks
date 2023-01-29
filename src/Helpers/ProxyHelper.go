package Helpers

import (
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"database/sql"
	"fmt"
)

func GetProxies(customerId int) []Structs.ProxyData {
	rows, err := Storage.GetDb().Query("SELECT id,ip,port,username,password,nickname FROM proxies WHERE customer_id= ?", customerId)

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

	var allProxies []Structs.ProxyData

	for rows.Next() {
		var id int
		var name string
		var ip string
		var port int
		var username string
		var password string

		err := rows.Scan(&id, &ip, &port, &username, &password, &name)
		if err != nil {
			return nil
		}

		allProxies = append(allProxies, Structs.ProxyData{
			Id:        id,
			Name:      name,
			IpAddress: ip,
			Port:      port,
			Username:  username,
			Password:  password,
		})
	}

	return allProxies
}
