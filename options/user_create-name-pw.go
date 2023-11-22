package options

import (
	"github.com/aivyss/jsonx"
	"github.com/aivyss/password-manager/options/parser"
	"github.com/urfave/cli/v2"
)

type userCreateOptNameOptPw struct {
	Name string `json:"name" annotation:"@NotBlank"`
}

func (o *userCreateOptNameOptPw) ToEntity(password string) (UserCreateOptNameOptPw, error) {
	var zeroValue UserCreateOptNameOptPw
	opts := UserCreateOptNameOptPw{
		Name:     o.Name,
		Password: password,
	}

	if err := jsonx.Validate(opts); err != nil {
		return zeroValue, err
	}

	return opts, nil
}

type UserCreateOptNameOptPw struct {
	Name     string `json:"name" annotation:"@NotBlank"`
	Password string `json:"pw" annotation:"@NotBlank"`
}

func NewuserCreateOptNameOptPw(c *cli.Context) (userCreateOptNameOptPw, error) {
	return parser.ParseOpts[userCreateOptNameOptPw](c, []parser.OptKeyValue{
		{
			Key:     "name",
			OptType: parser.STRING,
		},
	})
}
