package options

import (
	"github.com/aivyss/jsonx"
	"github.com/aivyss/password-manager/options/parser"
	"github.com/urfave/cli/v2"
)

type mainSaveOptKOptPw struct {
	Key string `json:"k" annotation:"@NotBlank"`
}

func (o *mainSaveOptKOptPw) ToEntity(password string) (*MainSaveOptKOptPw, error) {
	opts := MainSaveOptKOptPw{
		Key:      o.Key,
		Password: password,
	}

	if err := jsonx.Validate(opts); err != nil {
		return nil, err
	}

	return &opts, nil
}

type MainSaveOptKOptPw struct {
	Key      string `json:"k" annotation:"@NotBlank"`
	Password string `json:"pw" annotation:"@NotBlank"`
}

func NewmainSaveOptKOptPw(c *cli.Context) (*mainSaveOptKOptPw, error) {
	return parser.ParseOpts[mainSaveOptKOptPw](c, []parser.OptKeyValue{
		{
			Key:     "k",
			OptType: parser.STRING,
		},
	})
}
