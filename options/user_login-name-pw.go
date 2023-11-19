package options

import (
	"github.com/aivyss/jsonx"
	"github.com/aivyss/password-manager/options/parser"
	"github.com/urfave/cli/v2"
)

type userLoginOptNameOptPw struct {
	Name string `json:"name" annotation:"@NotBlank"`
}

func (o *userLoginOptNameOptPw) ToEntity(password string) (*UserLoginOptNameOptPw, error) {
	opts := UserLoginOptNameOptPw{
		Name:     o.Name,
		Password: password,
	}

	if err := jsonx.Validate(opts); err != nil {
		return nil, err
	}

	return &opts, nil
}

type UserLoginOptNameOptPw struct {
	Name     string `json:"name" annotation:"@NotBlank"`
	Password string `json:"pw" annotation:"@NotBlank"`
}

func NewuserLoginOptNameOptPw(c *cli.Context) (*userLoginOptNameOptPw, error) {
	return parser.ParseOpts[userLoginOptNameOptPw](c, []parser.OptKeyValue{
		{
			Key:     "name",
			OptType: parser.STRING,
		},
	})
}
