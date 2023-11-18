package entity

import "time"

type PasswordList struct {
	Id        int       `db:"ID"`
	UserPk    int       `db:"USER_PK"`
	Key       string    `db:"KEY_VALUE"`
	Password  string    `db:"PASSWORD"`
	CreatedAt time.Time `db:"CREATED_AT"`
	UpdatedAt time.Time `db:"UPDATED_AT"`
}
