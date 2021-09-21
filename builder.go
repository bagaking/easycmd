package easycmd

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"
)

type (
	Builder struct {
		Set     *Setter
		baseCli *cli.Command
		curCli  *cli.Command
	}
)

func newBuilder(base, cur *cli.Command) *Builder {
	builder := &Builder{
		Set:     &Setter{},
		baseCli: base,
		curCli:  cur,
	}
	builder.Set.End = builder
	return builder
}

// New method can create a New Builder by cmd name
func New(cmd string) *Builder {
	cmdObj := &cli.Command{Name: cmd}
	return newBuilder(cmdObj, cmdObj)
}

// Base returns base Builder
func (cb *Builder) Base() *Builder {
	return newBuilder(cb.baseCli, cb.baseCli)
}

// Cur returns current Builder
func (cb *Builder) Cur() *Builder {
	return newBuilder(cb.baseCli, cb.curCli)
}

// Flags override all flags of current Builder
func (cb *Builder) Flags(flags ...cli.Flag) *Builder {
	if len(flags) > 0 {
		cb.curCli.Flags = flags
	}
	return cb
}

// Child creates a child of current Builder by a given cmd name
func (cb *Builder) Child(cmd string) *Builder {
	if cb.curCli.Subcommands == nil {
		cb.curCli.Subcommands = make([]*cli.Command, 0)
	}

	var find *cli.Command = nil
	for _, v := range cb.curCli.Subcommands {
		if v.Name == cmd {
			find = v
			break
		}
	}

	if find != nil {
		return newBuilder(cb.baseCli, find)
	}

	newCmd := New(cmd)
	newCmd.baseCli = cb.baseCli

	cb.curCli.Subcommands = append(cb.curCli.Subcommands, newCmd.BuildCur())
	return newCmd
}

// SubCmd creates a child of current Builder by a given cli.Command
func (cb *Builder) SubCmd(child *cli.Command) *Builder {
	if cb.curCli.Subcommands == nil {
		cb.curCli.Subcommands = make([]*cli.Command, 0)
	}
	cb.curCli.Subcommands = append(cb.curCli.Subcommands, child)
	return cb
}

// Action sets the action of current Builder
func (cb *Builder) Action(action cli.ActionFunc) *Builder {
	cb.curCli.Action = action
	return cb
}

// Handler method sets the cmd by a ICliHandler
func (cb *Builder) Handler(handler ICliHandler, mws ...Middleware) *Builder {
	cb.curCli.Flags = handler.Flags()
	cb.curCli.Action = chain(append(mws, func(next cli.ActionFunc) cli.ActionFunc {
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

// RunBaseAsApp runs the command as a single app
func (cb *Builder) RunBaseAsApp() error {
	app, err := ToApp(cb.BuildBase())
	if err != nil {
		return err
	}
	return app.Run(os.Args)
}

// RunBaseAsAppWithCtx runs the command as a single app
func (cb *Builder) RunBaseAsAppWithCtx(ctx context.Context) error {
	app, err := ToApp(cb.BuildBase())
	if err != nil {
		return err
	}
	return app.RunContext(ctx, os.Args)
}

// BuildBase returns root Builder's cmd
func (cb *Builder) BuildBase() *cli.Command {
	return cb.baseCli
}

// BuildCur returns current Builder's cmd
func (cb *Builder) BuildCur() *cli.Command {
	return cb.curCli
}
