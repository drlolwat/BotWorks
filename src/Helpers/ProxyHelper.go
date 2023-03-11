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
	rows, err := Storage.GetDb().Query("SELECT id,ip,port,username,password,nickname,bans FROM proxies WHERE customer_id= ? ORDER BY nickname ASC", customerId)

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
		var bans int

		err := rows.Scan(&id, &ip, &port, &username, &password, &name, &bans)
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
			Bans:      bans,
		})
	}

	return allProxies
}

func GetUnusedProxy(customer int, maxBans int) int {
	res, _ := Storage.GetDb().Query("SELECT id FROM proxies WHERE customer_id = ? AND bans <= ? ORDER BY RANDOM()", customer, maxBans)
	ret := -1

	for res.Next() {
		var dbId int
		err := res.Scan(&dbId)

		if err != nil {
			fmt.Println(err)
			break
		}

		if dbId > 0 {
			var used int
			err = Storage.GetDb().QueryRow("SELECT proxy_id FROM accounts WHERE proxy_id = ?", dbId).Scan(&used)

			if err != nil || err == sql.ErrNoRows {
				ret = dbId
			}
		}
	}

	_ = res.Close()
	return ret
}
