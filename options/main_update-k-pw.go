package options

import (
	"github.com/aivyss/jsonx"
	"github.com/aivyss/password-manager/options/parser"
	"github.com/urfave/cli/v2"
)

type mainUpdateOptKOptPw struct {
	Key string `json:"k" annotation:"@NotBlank"`
}

func (o *mainUpdateOptKOptPw) ToEntity(password string) (*MainUpdateOptKOptPw, error) {
	opts := MainUpdateOptKOptPw{
		Key:      o.Key,
		Password: password,
	}

	if err := jsonx.Validate(opts); err != nil {
		return nil, err
	}

	return &opts, nil
}

type MainUpdateOptKOptPw struct {
	Key      string `json:"k" annotation:"@NotBlank"`
	Password string `json:"pw"`
}

func NewmainUpdateOptKOptPw(c *cli.Context) (*mainUpdateOptKOptPw, error) {
	return parser.ParseOpts[mainUpdateOptKOptPw](c, []parser.OptKeyValue{
		{
			Key:     "k",
			OptType: parser.STRING,
		},
	})
}
