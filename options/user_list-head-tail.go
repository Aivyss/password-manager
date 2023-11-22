package options

import (
	"github.com/aivyss/password-manager/options/parser"
	"github.com/urfave/cli/v2"
)

type UserListOptHeadOptTail struct {
	Head int `json:"head"`
	Tail int `json:"tail"`
}

func NewUserListOptHeadOptTail(c *cli.Context) (UserListOptHeadOptTail, error) {
	return parser.ParseOpts[UserListOptHeadOptTail](c, []parser.OptKeyValue{
		{
			Key:     "head",
			OptType: parser.INT,
		},
		{
			Key:     "tail",
			OptType: parser.INT,
		},
	})
}
