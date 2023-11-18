package options

import (
	"github.com/aivyss/password-manager/options/parser"
	"github.com/urfave/cli/v2"
)

type MainUpdateOptKOptPw struct {
	Key      string `json:"k" annotation:"@NotBlank"`
	Password string `json:"pw" annotation:"@NotBlank"`
}

func NewMainUpdateOptKOptPw(c *cli.Context) (*MainUpdateOptKOptPw, error) {
	return parser.ParseOpts[MainUpdateOptKOptPw](c, []parser.OptKeyValue{
		{
			Key:     "k",
			OptType: parser.STRING,
		},
		{
			Key:     "pw",
			OptType: parser.STRING,
		},
	})
}
