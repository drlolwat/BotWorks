package Storage

import (
	"database/sql"
	"fmt"
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

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='configuration';")
	if err != nil {
		return
	}
	defer rows.Close()

	setupQueries := []string{
		`CREATE TABLE "account_groups" (
					  "id" integer PRIMARY KEY AUTOINCREMENT,
					  "name" text,
					  "customer_id" integer
					);`,
		`INSERT INTO "account_groups" VALUES (-1, 'Ungrouped', -1);`,
		`CREATE TABLE "account_schedules" (
					  "id" integer PRIMARY KEY AUTOINCREMENT,
					  "name" text,
					  "start_hour" integer,
					  "start_min" integer,
					  "end_hour" integer,
					  "end_min" integer,
					  "customer_id" integer
					);`,
		`DROP TABLE IF EXISTS "accounts";`,
		`CREATE TABLE "accounts" (
					  "id" integer PRIMARY KEY AUTOINCREMENT,
					  "name" text,
					  "password" text,
					  "bank_pin" integer,
					  "proxy_id" integer,
					  "script_id" integer,
					  "script_params" text,
					  "customer_id" integer,
					  "logfile" integer,
					  "pid" integer,
					  "group_id" integer,
					  "breaks" text,
					  "schedule" integer,
					  "fps" integer DEFAULT 20
					);`,
		`DROP TABLE IF EXISTS "configuration";`,
		`CREATE TABLE "configuration" (
					  "db_jar" text,
					  "db_logs" text,
					  "db_username" text,
					  "db_password" text
					);`,
		`INSERT INTO "configuration" VALUES ('/home', '/home', 'dreambot', 'dreambot');`,
		`DROP TABLE IF EXISTS "customers";`,
		`CREATE TABLE "customers" (
					  "id" integer PRIMARY KEY AUTOINCREMENT,
					  "handle" text,
					  "password" text,
					  "bot_limit" integer,
					  "permission_level" integer
					);`,
		`DROP TABLE IF EXISTS "proxies";`,
		`CREATE TABLE "proxies" (
					  "id" integer PRIMARY KEY AUTOINCREMENT,
					  "ip" text,
					  "port" integer,
					  "username" tex,
					  "password" text,
					  "nickname" text,
					  "customer_id" integer
					);`,
		`DROP TABLE IF EXISTS "scripts";`,
		`CREATE TABLE "scripts" (
					  "id" integer PRIMARY KEY AUTOINCREMENT,
					  "name" text,
					  "backend_name" text
					);`,
		`INSERT INTO "scripts" VALUES (-1, 'No Script', 'No Script');`,
		`INSERT INTO "sqlite_sequence" VALUES ('accounts', 0);`,
		`INSERT INTO "sqlite_sequence" VALUES ('account_groups', 0);`,
		`INSERT INTO "sqlite_sequence" VALUES ('account_schedules', 0);`,
		`INSERT INTO "sqlite_sequence" VALUES ('customers', 0);`,
		`INSERT INTO "sqlite_sequence" VALUES ('proxies', 0);`,
		`INSERT INTO "sqlite_sequence" VALUES ('scripts', 0);`,
		`UPDATE "sqlite_sequence" SET seq = 0 WHERE name = 'account_groups';`,
		`UPDATE "sqlite_sequence" SET seq = 0 WHERE name = 'account_schedules';`,
		`UPDATE "sqlite_sequence" SET seq = 0 WHERE name = 'accounts';`,
		`UPDATE "sqlite_sequence" SET seq = 0 WHERE name = 'customers';`,
		`UPDATE "sqlite_sequence" SET seq = 0 WHERE name = 'proxies';`,
		`UPDATE "sqlite_sequence" SET seq = 0 WHERE name = 'scripts';`,
		`PRAGMA foreign_keys = true;`,
	}

	if !rows.Next() {
		for _, query := range setupQueries {
			_, err = db.Exec(query)
			if err != nil {
				fmt.Print(err)
			}
		}

	}
}

func GetDb() *sql.DB {
	return db
}
