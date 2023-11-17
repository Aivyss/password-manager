package repository

const (
	createMasterUserTable = `
		CREATE TABLE IF NOT EXISTS MASTER_USER (
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
	GetMasterUserById = `
		SELECT
			ID,
			USERNAME,
			USER_PASSWORD
		FROM
			MASTER_USER
		WHERE 1=1
			AND ID = :id
	`
)

const (
	createPasswordListTable = `
		CREATE TABLE IF NOT EXISTS PASSWORD_LIST(
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			USER_PK INTEGER NOT NULL,
			KEY_VALUE TEXT NOT NULL,
			PASSWORD TEXT NOT NULL
		)
	`
	InsertPassword = `
		INSERT INTO PASSWORD_LIST (USER_PK, KEY_VALUE, PASSWORD) VALUES (:userPk, :key, :password)
	`
	GetPasswordByUserPkAndKey = `
		SELECT
			ID,
			USER_PK,
			KEY_VALUE,
			PASSWORD
		FROM
			PASSWORD_LIST
		WHERE 1=1
			AND USER_PK = :userPk
			AND KEY_VALUE = :key
	`
	UpdatePasswordByUserPkAndKey = `
		UPDATE
			PASSWORD_LIST
		SET
			PASSWORD = :password
		WHERE 1=1
			AND USER_PK = :userPk
			AND KEY_VALUE = :key 
	`
)
