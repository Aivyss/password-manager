package parser

import (
	"github.com/aivyss/jsonx"
	"github.com/aivyss/password-manager/pwmErr"
	"github.com/urfave/cli/v2"
	"regexp"
	"strings"
)

type OptType int

const (
	STRING OptType = iota
	INT
)

var optValuePattern = regexp.MustCompile("-[a-zA-Z\\d]*\\s*[^-]*")
var optPattern = regexp.MustCompile("-[a-zA-Z\\d]*\\s*")
var patternMultipleSpacePattern = regexp.MustCompile("\\b\\s+\\b")

type OptKeyValue struct {
	Key     string
	OptType OptType
}

func ParseOpts[T any](c *cli.Context, opts []OptKeyValue) (T, error) {
	var result T
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
		return result, err
	}

	unmarshal, err := jsonx.Unmarshal[T](j)
	if err != nil {
		return result, pwmErr.OptParseErr
	}
	result = *unmarshal

	return result, nil
}

func ParseCommand(s string) []string {
	opts := optValuePattern.FindAllString(s, -1)
	nonOptsStr := s

	for _, opt := range opts {
		nonOptsStr = strings.TrimSpace(strings.ReplaceAll(nonOptsStr, opt, ""))
	}
	nonOptsStr = string(patternMultipleSpacePattern.ReplaceAll([]byte(nonOptsStr), []byte(" ")))

	command := make([]string, 0, len(opts)+1)
	for _, nonOpt := range strings.Split(nonOptsStr, " ") {
		command = append(command, strings.TrimSpace(nonOpt))
	}

	for _, opt := range opts {
		opt = string(patternMultipleSpacePattern.ReplaceAll([]byte(opt), []byte(" ")))
		value := string(optPattern.ReplaceAll([]byte(opt), []byte("")))
		flag := strings.ReplaceAll(opt, value, "")
		command = append(command, strings.TrimSpace(flag))
		command = append(command, strings.TrimSpace(value))
	}

	return command
}
