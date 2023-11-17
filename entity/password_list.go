package entity

type PasswordList struct {
	Id       int    `db:"ID"`
	UserPk   int    `db:"USER_PK"`
	Key      string `db:"KEY_VALUE"`
	Password string `db:"PASSWORD"`
}
