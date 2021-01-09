package easycmd

import (
	"os"

	"github.com/thoas/go-funk"
	"github.com/urfave/cli/v2"
)

type (
	Builder struct {
		Set  *Setter
		root *cli.Command
		cur  *cli.Command
	}
)

func newBuilder(root, cur *cli.Command) *Builder {
	builder := &Builder{
		Set:  &Setter{},
		root: root,
		cur:  cur,
	}
	builder.Set.End = builder
	return builder
}

func New(cmd string) *Builder {
	cmdObj := &cli.Command{Name: cmd}
	return newBuilder(cmdObj, cmdObj)
}

func (cb *Builder) Root() *cli.Command {
	return cb.root
}

func (cb *Builder) Cur() *Builder {
	return newBuilder(cb.root, cb.cur)
}

func (cb *Builder) Flags(flags ...cli.Flag) *Builder {
	if len(flags) > 0 {
		cb.cur.Flags = flags
	}
	return cb
}

func (cb *Builder) Child(cmd string) *Builder {
	if cb.cur.Subcommands == nil {
		cb.cur.Subcommands = make([]*cli.Command, 0)
	}

	_, v := funk.FindKey(cb.cur.Subcommands, func(sc *cli.Command) bool {
		return sc.Name == cmd
	})

	if v != nil {
		return newBuilder(cb.root, v.(*cli.Command))
	}

	newCmd := New(cmd)
	cb.cur.Subcommands = append(cb.cur.Subcommands, newCmd.Root())
	return newCmd
}

func (cb *Builder) SubCmd(child *cli.Command) *Builder {
	if cb.cur.Subcommands == nil {
		cb.cur.Subcommands = make([]*cli.Command, 0)
	}
	cb.cur.Subcommands = append(cb.cur.Subcommands, child)
	return cb
}

func (cb *Builder) Action(action cli.ActionFunc) *Builder {
	cb.cur.Action = action
	return cb
}

// Handler method sets the cmd by a ICliHandler
func (cb *Builder) Handler(handler ICliHandler, mws ...Middleware) *Builder {
	cb.cur.Flags = handler.Flags()
	cb.cur.Action = chain(append(mws, func(next cli.ActionFunc) cli.ActionFunc {
		return func(c *cli.Context) (err error) {
			handler, err = handler.Parse(c)
			if err != nil {
				return err
			}
			return next(c)
		}
	})...)(handler.Handle)
	return cb
}

// RunAsApp runs the command as a single app
func (cb *Builder) RunAsApp() error {
	app, err := ToApp(cb.Root())
	if err != nil {
		return err
	}
	return app.Run(os.Args)
}
