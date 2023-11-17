package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

func InitDB(db *sqlx.DB) {
	_, err := db.Exec(createMasterUserTable)
	if err != nil {
		e, ok := err.(sqlite3.Error)
		if ok && e.Code == 1 {
			return
		}
	}
}
