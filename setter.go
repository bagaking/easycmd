package easycmd

import "github.com/urfave/cli/v2"

type (
	Setter struct {
		End *Builder
	}
)

func (set *Setter) Custom(setter func(*cli.Command)) *Setter {
	setter(set.End.cur)
	return set
}

func (set *Setter) Usage(usage string) *Setter {
	set.End.cur.Usage = usage
	return set
}

func (set *Setter) Alias(a string) *Setter {
	if set.End.cur.Aliases == nil {
		set.End.cur.Aliases = []string{a}
	} else {
		set.End.cur.Aliases = append(set.End.cur.Aliases, a)
	}
	return set
}
