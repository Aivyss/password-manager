package repository

import (
	"github.com/aivyss/bean"
	"github.com/aivyss/password-manager/constant"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/aivyss/password-manager/util"
	"github.com/jmoiron/sqlx"
	"strings"
)

func InitDB(db *sqlx.DB) error {
	_, err := db.Exec(createAppVersionTable)
	if err != nil {
		return pwmErr.DBInit
	}

	_, err = db.Exec(createMasterUserTable)
	if err != nil {
		return pwmErr.DBInit
	}

	_, err = db.Exec(createPasswordListTable)
	if err != nil {
		return pwmErr.DBInit
	}

	_, err = db.Exec(createPasswordListDetailTable)
	if err != nil {
		return pwmErr.DBInit
	}

	return nil
}

func CheckAppVersion() error {
	appVersionRepo, err := bean.GetBean[AppVersionRepository]()
	if err != nil {
		return err
	}

	versionCount, err := appVersionRepo.CountVersions()
	if err != nil {
		return err
	}

	if versionCount > 0 {
		version, err := appVersionRepo.GetLatestAppVersion()
		if err != nil {
			return err
		}

		parseRegisteredVerStrs := strings.Split(version, ".")
		currentAppVerStrs := strings.Split(constant.AppVersion, ".")

		parseRegisteredVer := appVersionDataBind{
			main:   util.MustAtoi(parseRegisteredVerStrs[0]),
			middle: util.MustAtoi(parseRegisteredVerStrs[1]),
			minor:  util.MustAtoi(parseRegisteredVerStrs[2]),
		}
		currentAppVer := appVersionDataBind{
			main:   util.MustAtoi(currentAppVerStrs[0]),
			middle: util.MustAtoi(currentAppVerStrs[1]),
			minor:  util.MustAtoi(currentAppVerStrs[2]),
		}

		if currentAppVer.After(parseRegisteredVer) {
			err := appVersionRepo.InsertAppVersion()
			if err != nil {
				return err
			}

			return nil
		} else if currentAppVer.Equal(parseRegisteredVer) {
			return nil
		}

		return pwmErr.AppVersionUnknown
	}

	return appVersionRepo.InsertAppVersion()
}

type appVersionDataBind struct {
	main   int
	middle int
	minor  int
}

func (ver *appVersionDataBind) After(v appVersionDataBind) bool {
	thisVersion := ver.main*constant.MainVersionUnit + ver.middle*constant.MiddleVersionUnit + ver.minor*constant.MinorVersionUnit
	thatVersion := v.main*constant.MainVersionUnit + v.middle*constant.MiddleVersionUnit + v.minor*constant.MinorVersionUnit

	return thisVersion > thatVersion
}

func (ver *appVersionDataBind) Equal(v appVersionDataBind) bool {
	thisVersion := ver.main*constant.MainVersionUnit + ver.middle*constant.MiddleVersionUnit + ver.minor*constant.MinorVersionUnit
	thatVersion := v.main*constant.MainVersionUnit + v.middle*constant.MiddleVersionUnit + v.minor*constant.MinorVersionUnit

	return thisVersion == thatVersion
}
