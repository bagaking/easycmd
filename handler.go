package easycmd

import "github.com/urfave/cli/v2"

type ICliHandler interface {
	Flags() []cli.Flag
	Parse(c *cli.Context) (ICliHandler, error)
	Handle(c *cli.Context) error
}
