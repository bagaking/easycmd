package easycmd

import (
	"github.com/urfave/cli/v2"
)

type (
	CMDSetter struct {
		End *CMDBuilder
	}
)

func (set *CMDSetter) Custom(setter func(*cli.Command)) *CMDSetter {
	setter(set.End.cur)
	return set
}

func (set *CMDSetter) Usage(usage string) *CMDSetter {
	set.End.cur.Usage = usage
	return set
}

func (set *CMDSetter) Alias(a string) *CMDSetter {
	if set.End.cur.Aliases == nil {
		set.End.cur.Aliases = []string{a}
	} else {
		set.End.cur.Aliases = append(set.End.cur.Aliases, a)
	}
	return set
}
