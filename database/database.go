package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func OpenDatabase() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite", "database/users.db")
	if err != nil {
		fmt.Printf("error opening database: %v", err)
	}
	return
}
