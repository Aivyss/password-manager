package repository

import (
	"github.com/aivyss/password-manager/constant"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/jmoiron/sqlx"
)

type AppVersionRepository interface {
	CountVersions() (int, error)
	InsertAppVersion() error
	GetLatestAppVersion() (string, error)
}

type appVersionDbDataBindObject struct {
	Version string `db:"VERSION"`
}

type appVersionRepository struct {
	db                  *sqlx.DB
	insertAppVersion    *sqlx.NamedStmt
	getLatestAppVersion *sqlx.NamedStmt
	countVersions       *sqlx.NamedStmt
}

func (r *appVersionRepository) CountVersions() (int, error) {
	type countDbBindObject struct {
		Count int `db:"COUNT"`
	}

	result := new(countDbBindObject)
	if err := r.countVersions.Get(result, map[string]any{}); err != nil {
		return 0, pwmErr.AppVersionUnknown
	}

	return result.Count, nil
}

func (r *appVersionRepository) GetLatestAppVersion() (string, error) {
	result := new(appVersionDbDataBindObject)
	if err := r.getLatestAppVersion.Get(result, map[string]any{}); err != nil {
		return "", pwmErr.AppVersionUnknown
	}

	return result.Version, nil
}

func (r *appVersionRepository) InsertAppVersion() error {
	if _, err := r.insertAppVersion.Exec(map[string]any{
		"version": constant.AppVersion,
	}); err != nil {
		return pwmErr.AppVersionUnknown
	}

	return nil
}

func NewAppVersionRepository(db *sqlx.DB) (AppVersionRepository, error) {
	insertAppVersion, err := db.PrepareNamed(InsertAppVersion)
	if err != nil {
		return nil, err
	}
	getLatestAppVersion, err := db.PrepareNamed(GetLatestAppVersion)
	if err != nil {
		return nil, err
	}
	countVersions, err := db.PrepareNamed(CountVersions)
	if err != nil {
		return nil, err
	}

	return &appVersionRepository{
		db:                  db,
		insertAppVersion:    insertAppVersion,
		getLatestAppVersion: getLatestAppVersion,
		countVersions:       countVersions,
	}, nil
}
