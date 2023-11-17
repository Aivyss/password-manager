package repository

import (
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/jmoiron/sqlx"
)

func InitDB(db *sqlx.DB) error {
	_, err := db.Exec(createMasterUserTable)
	if err != nil {
		return pwmErr.DBInit
	}

	_, err = db.Exec(createPasswordListTable)
	if err != nil {
		return pwmErr.DBInit
	}

	return nil
}
