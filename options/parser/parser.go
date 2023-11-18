package parser

import (
	"github.com/aivyss/jsonx"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/urfave/cli/v2"
)

type OptType int

const (
	STRING OptType = iota
	INT
)

type OptKeyValue struct {
	Key     string
	OptType OptType
}

func ParseOpts[T any](c *cli.Context, opts []OptKeyValue) (*T, error) {
	m := map[string]any{}
	for _, opt := range opts {
		switch opt.OptType {
		case STRING:
			m[opt.Key] = c.String(opt.Key)
		case INT:
			m[opt.Key] = c.Int(opt.Key)
		}
	}

	j, err := jsonx.Marshal(m)
	if err != nil {
		return nil, err
	}

	unmarshal, err := jsonx.Unmarshal[T](j)
	if err != nil {
		return nil, pwmErr.OptParseErr
	}

	return unmarshal, nil
}
