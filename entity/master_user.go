package entity

type MasterUser struct {
	Id       int    `db:"ID"`
	UserName string `db:"USERNAME"`
	Password string `db:"USER_PASSWORD"`
}
