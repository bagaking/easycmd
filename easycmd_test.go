package easycmd

import (
	"errors"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestMergeFlagsReturnsOnlyMergedFlags(t *testing.T) {
	flags, err := MergeFlags(
		[]cli.Flag{
			&cli.StringFlag{Name: "config", Aliases: []string{"c"}},
		},
		[]cli.Flag{
			&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}},
		},
	)
	if err != nil {
		t.Fatalf("MergeFlags returned error: %v", err)
	}

	if len(flags) != 2 {
		t.Fatalf("expected 2 flags, got %d: %#v", len(flags), flags)
	}
	for i, flag := range flags {
		if flag == nil {
			t.Fatalf("flag %d is nil", i)
		}
	}
}

func TestMergeFlagsRejectsDuplicateNames(t *testing.T) {
	_, err := MergeFlags(
		[]cli.Flag{&cli.StringFlag{Name: "config"}},
		[]cli.Flag{&cli.BoolFlag{Name: "config"}},
	)
	if !errors.Is(err, ErrFlagAlreadyExist) {
		t.Fatalf("expected ErrFlagAlreadyExist, got %v", err)
	}
}

func TestMergeFlagsRejectsNameAliasConflicts(t *testing.T) {
	tests := []struct {
		name   string
		flags1 []cli.Flag
		flags2 []cli.Flag
	}{
		{
			name:   "existing name conflicts with new alias",
			flags1: []cli.Flag{&cli.StringFlag{Name: "config"}},
			flags2: []cli.Flag{&cli.BoolFlag{Name: "color", Aliases: []string{"config"}}},
		},
		{
			name:   "existing alias conflicts with new name",
			flags1: []cli.Flag{&cli.StringFlag{Name: "config", Aliases: []string{"c"}}},
			flags2: []cli.Flag{&cli.BoolFlag{Name: "c"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MergeFlags(tt.flags1, tt.flags2)
			if !errors.Is(err, ErrFlagAlreadyExist) {
				t.Fatalf("expected ErrFlagAlreadyExist, got %v", err)
			}
		})
	}
}

func TestMergeFlagsCanIgnoreDuplicateNames(t *testing.T) {
	flags, err := MergeFlags(
		[]cli.Flag{&cli.StringFlag{Name: "config"}},
		[]cli.Flag{&cli.BoolFlag{Name: "config"}},
		true,
	)
	if err != nil {
		t.Fatalf("MergeFlags returned error: %v", err)
	}
	if len(flags) != 1 {
		t.Fatalf("expected duplicate flag to be skipped, got %d flags", len(flags))
	}
}

func TestMergeFlagsCanIgnoreDuplicateAliases(t *testing.T) {
	flags, err := MergeFlags(
		[]cli.Flag{&cli.StringFlag{Name: "config", Aliases: []string{"c"}}},
		[]cli.Flag{&cli.BoolFlag{Name: "color", Aliases: []string{"c"}}},
		true,
	)
	if err != nil {
		t.Fatalf("MergeFlags returned error: %v", err)
	}
	if len(flags) != 1 {
		t.Fatalf("expected alias-conflicting flag to be skipped, got %d flags", len(flags))
	}
	if flags[0].Names()[0] != "config" {
		t.Fatalf("expected original flag to be preserved, got names %#v", flags[0].Names())
	}
}

func TestBuilderFlagsWithNoArgsClearsCurrentFlags(t *testing.T) {
	builder := New("root").Flags(&cli.StringFlag{Name: "config"})

	builder.Flags()

	if flags := builder.BuildCur().Flags; len(flags) != 0 {
		t.Fatalf("expected flags to be cleared, got %#v", flags)
	}
}
