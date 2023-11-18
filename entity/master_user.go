package entity

import "time"

type MasterUser struct {
	Id        int       `db:"ID"`
	UserName  string    `db:"USERNAME"`
	Password  string    `db:"USER_PASSWORD"`
	CreatedAt time.Time `db:"CREATED_AT"`
	UpdatedAt time.Time `db:"UPDATED_AT"`
}
