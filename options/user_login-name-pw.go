package options

import (
	"github.com/aivyss/password-manager/options/parser"
	"github.com/urfave/cli/v2"
)

type UserLoginOptNameOptPw struct {
	Name     string `json:"name" annotation:"@NotBlank"`
	Password string `json:"pw" annotation:"@NotBlank"`
}

func NewUserLoginOptNameOptPw(c *cli.Context) (*UserLoginOptNameOptPw, error) {
	return parser.ParseOpts[UserLoginOptNameOptPw](c, []parser.OptKeyValue{
		{
			Key:     "name",
			OptType: parser.STRING,
		},
		{
			Key:     "pw",
			OptType: parser.STRING,
		},
	})
}
