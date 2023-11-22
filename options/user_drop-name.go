package options

import (
	"github.com/aivyss/jsonx"
	"github.com/aivyss/password-manager/options/parser"
	"github.com/urfave/cli/v2"
)

type userDropOptNameOptPw struct {
	Name string `json:"name" annotation:"@NotBlank"`
}

func (o *userDropOptNameOptPw) ToEntity(password string) (UserDropOptNameOptPw, error) {
	var zeroValue UserDropOptNameOptPw
	opts := UserDropOptNameOptPw{
		Name:     o.Name,
		Password: password,
	}

	if err := jsonx.Validate(opts); err != nil {
		return zeroValue, err
	}

	return opts, nil
}

type UserDropOptNameOptPw struct {
	Name     string `json:"name" annotation:"@NotBlank"`
	Password string `json:"pw" annotation:"@NotBlank"`
}

func NewuserDropOptName(c *cli.Context) (userDropOptNameOptPw, error) {
	return parser.ParseOpts[userDropOptNameOptPw](c, []parser.OptKeyValue{
		{
			Key:     "name",
			OptType: parser.STRING,
		},
	})
}
