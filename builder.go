package easycmd

import (
	"github.com/thoas/go-funk"
	"github.com/urfave/cli/v2"
)

type (
	CMDBuilder struct {
		Set  *CMDSetter
		root *cli.Command
		cur  *cli.Command
	}
)

func newBuilder(root, cur *cli.Command) *CMDBuilder {
	builder := &CMDBuilder{
		Set:  &CMDSetter{},
		root: root,
		cur:  cur,
	}
	builder.Set.End = builder
	return builder
}

func New(cmd string) *CMDBuilder {
	cmdObj := &cli.Command{Name: cmd}
	return newBuilder(cmdObj, cmdObj)
}

func (cb *CMDBuilder) Root() *cli.Command {
	return cb.root
}

func (cb *CMDBuilder) Cur() *CMDBuilder {
	return newBuilder(cb.root, cb.cur)
}

func (cb *CMDBuilder) Flags(flags ...cli.Flag) *CMDBuilder {
	if len(flags) > 0 {
		cb.cur.Flags = flags
	}
	return cb
}

func (cb *CMDBuilder) Child(cmd string) *CMDBuilder {
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

func (cb *CMDBuilder) SubCmd(child *cli.Command) *CMDBuilder {
	if cb.cur.Subcommands == nil {
		cb.cur.Subcommands = make([]*cli.Command, 0)
	}
	cb.cur.Subcommands = append(cb.cur.Subcommands, child)
	return cb
}

func (cb *CMDBuilder) Action(action cli.ActionFunc) *CMDBuilder {
	cb.cur.Action = action
	return cb
}

func (cb *CMDBuilder) Handler(handler ICliHandler, mws ...Middleware) *CMDBuilder {
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
