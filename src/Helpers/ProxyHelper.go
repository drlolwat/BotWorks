package Helpers

import (
	"BotWorks/src/Storage"
	"BotWorks/src/Structs"
	"database/sql"
	"fmt"
)

func NumAccountsUsing(proxyId int) int {
	id := 0
	rows, err := Storage.GetDb().Query("SELECT count(id) AS num FROM accounts WHERE proxy_id = ?", proxyId)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0
		}
	}

	return id
}

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
		var port string
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
			Accounts:  NumAccountsUsing(id),
		})
	}

	return allProxies
}
