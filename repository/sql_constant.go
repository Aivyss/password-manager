package repository

const (
	createMasterUserTable = `
		CREATE TABLE MASTER_USER (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			USERNAME TEXT NOT NULL UNIQUE,
			USER_PASSWORD TEXT NOT NULL
		)
	`
	InsertMasterUser = `
		INSERT INTO MASTER_USER (USERNAME, USER_PASSWORD) VALUES (:userName, :userPassword)
	`
	GetMasterUserByUserName = `
		SELECT
			ID,
			USERNAME,
			USER_PASSWORD
		FROM
			MASTER_USER
		WHERE 1=1
			AND USERNAME = :userName
	`
)
