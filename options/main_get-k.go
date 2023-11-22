package options

import (
	"github.com/aivyss/password-manager/options/parser"
	"github.com/urfave/cli/v2"
)

type MainGetOptK struct {
	Key string `json:"k" annotation:"@NotBlank"`
}

func NewMainGetOptK(c *cli.Context) (MainGetOptK, error) {
	return parser.ParseOpts[MainGetOptK](c, []parser.OptKeyValue{
		{
			Key:     "k",
			OptType: parser.STRING,
		},
	})
}
