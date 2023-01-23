package Storage

import "database/sql"

var (
	db *sql.DB
)

func Init() {
	var err error
	db, err = sql.Open("sqlite3", "./Botworks.db")
	if err != nil {
		println(err)
	}
}

func GetDb() *sql.DB {
	return db
}
