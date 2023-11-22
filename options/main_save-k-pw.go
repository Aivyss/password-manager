package options

import (
	"github.com/aivyss/jsonx"
	"github.com/aivyss/password-manager/options/parser"
	"github.com/urfave/cli/v2"
)

type mainSaveOptKOptPw struct {
	Key         string `json:"k" annotation:"@NotBlank"`
	Description string `json:"description"`
}

func (o *mainSaveOptKOptPw) ToEntity(password string) (MainSaveOptKOptPw, error) {
	var zeroValue MainSaveOptKOptPw
	opts := MainSaveOptKOptPw{
		Key:         o.Key,
		Password:    password,
		Description: o.Description,
	}

	if err := jsonx.Validate(opts); err != nil {
		return zeroValue, err
	}

	return opts, nil
}

type MainSaveOptKOptPw struct {
	Key         string `json:"k" annotation:"@NotBlank"`
	Password    string `json:"pw" annotation:"@NotBlank"`
	Description string `json:"description"`
}

func NewmainSaveOptKOptPw(c *cli.Context) (mainSaveOptKOptPw, error) {
	return parser.ParseOpts[mainSaveOptKOptPw](c, []parser.OptKeyValue{
		{
			Key:     "k",
			OptType: parser.STRING,
		},
		{
			Key:     "description",
			OptType: parser.STRING,
		},
	})
}
