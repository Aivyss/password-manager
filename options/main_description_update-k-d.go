package options

import (
	"github.com/aivyss/jsonx"
	"github.com/aivyss/password-manager/options/parser"
	"github.com/urfave/cli/v2"
)

type MainDescriptionUpdateOptKeyOptValue struct {
	Key      string `json:"key" annotation:"@NotBlank"`
	Value    string `json:"value"`
	Password string `json:"password" annotation:"@NotBlank"`
}

func (o mainDescriptionUpdateOptKeyOptValue) ToEntity(password string) (MainDescriptionUpdateOptKeyOptValue, error) {
	opts := MainDescriptionUpdateOptKeyOptValue{
		Key:      o.Key,
		Value:    o.Value,
		Password: password,
	}
	if err := jsonx.Validate(opts); err != nil {
		return MainDescriptionUpdateOptKeyOptValue{}, err
	}

	return opts, nil
}

type mainDescriptionUpdateOptKeyOptValue struct {
	Key   string `json:"key" annotation:"@NotBlank"`
	Value string `json:"value"`
}

func NewmainDescriptionUpdateOptKeyOptValue(c *cli.Context) (mainDescriptionUpdateOptKeyOptValue, error) {
	return parser.ParseOpts[mainDescriptionUpdateOptKeyOptValue](c, []parser.OptKeyValue{
		{
			Key:     "key",
			OptType: parser.STRING,
		},
		{
			Key:     "value",
			OptType: parser.STRING,
		},
	})
}
