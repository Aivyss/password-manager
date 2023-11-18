package util

import (
	"github.com/aivyss/password-manager/pwmErr"
	"strconv"
)

func MustAtoi(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(pwmErr.Unknown)
	}

	return num
}
