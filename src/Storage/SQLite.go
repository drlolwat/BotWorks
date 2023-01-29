package Storage

import (
	"database/sql"
	"strings"
)

var (
	db *sql.DB
)

func Init() {
	var err error
	db, err = sql.Open("sqlite3", "./BotBuddy.db")
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return
	}

	if err != nil {
		println(err)
	}
}

func SanitizeInput(input string) string {
	// Remove any semicolons or single quotes from the input
	input = strings.Replace(input, ";", "", -1)
	input = strings.Replace(input, "'", "", -1)
	input = strings.Replace(input, "\"", "\"", -1)

	// Return the sanitized input
	return input
}

func GetDb() *sql.DB {
	return db
}
