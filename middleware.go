package easycmd

import "github.com/urfave/cli/v2"

type Middleware func(cli.ActionFunc) cli.ActionFunc

func chain(mws ...Middleware) Middleware {
	return func(handler cli.ActionFunc) cli.ActionFunc {
		for _, mw := range mws {
			handler = mw(handler)
		}
		return handler
	}
}
