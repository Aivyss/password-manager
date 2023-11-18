package options

import (
	"github.com/aivyss/password-manager/options/parser"
	"github.com/urfave/cli/v2"
)

type MainSaveOptKOptPw struct {
	Key      string `json:"k" annotation:"@NotBlank"`
	Password string `json:"pw" annotation:"@NotBlank"`
}

func NewMainSaveOptKOptPw(c *cli.Context) (*MainSaveOptKOptPw, error) {
	return parser.ParseOpts[MainSaveOptKOptPw](c, []parser.OptKeyValue{
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
