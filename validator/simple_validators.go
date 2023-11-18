package validator

import (
	"github.com/aivyss/password-manager/options"
	"github.com/aivyss/password-manager/pwmErr"
)

type UserCreateOptNameOptPwValidator int

func (vv *UserCreateOptNameOptPwValidator) Validate(v options.UserCreateOptNameOptPw) error {
	if len(v.Password) < 16 {
		return pwmErr.InvalidOpt
	}

	return nil
}

type UserLoginOptNameOptPwValidator int

func (vv *UserLoginOptNameOptPwValidator) Validate(v options.UserLoginOptNameOptPw) error {
	if len(v.Password) < 16 {
		return pwmErr.InvalidOpt
	}

	return nil
}
