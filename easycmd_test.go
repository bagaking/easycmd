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
