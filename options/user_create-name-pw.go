package options

import (
	"github.com/aivyss/password-manager/options/parser"
	"github.com/urfave/cli/v2"
)

type UserCreateOptNameOptPw struct {
	Name     string `json:"name" annotation:"@NotBlank"`
	Password string `json:"pw" annotation:"@NotBlank"`
}

func NewUserCreateOptNameOptPw(c *cli.Context) (*UserCreateOptNameOptPw, error) {
	return parser.ParseOpts[UserCreateOptNameOptPw](c, []parser.OptKeyValue{
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
