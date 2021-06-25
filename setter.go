package easycmd

import "github.com/urfave/cli/v2"

type (
	Setter struct {
		End *Builder
	}
)

func (set *Setter) Custom(setter func(*cli.Command)) *Setter {
	setter(set.End.curCli)
	return set
}

func (set *Setter) Usage(usage string) *Setter {
	set.End.curCli.Usage = usage
	return set
}

func (set *Setter) Alias(a string) *Setter {
	if set.End.curCli.Aliases == nil {
		set.End.curCli.Aliases = []string{a}
	} else {
		set.End.curCli.Aliases = append(set.End.curCli.Aliases, a)
	}
	return set
}
