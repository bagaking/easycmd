package easycmd

import (
	"bytes"
	"errors"
	"fmt"
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

func TestToAppDoesNotClearCommandFlags(t *testing.T) {
	cmd := New("root").Flags(&cli.StringFlag{Name: "config"}).BuildBase()

	app, err := ToApp(cmd)
	if err != nil {
		t.Fatalf("ToApp(root command with config flag) returned error: %v", err)
	}
	if flags := app.Flags; len(flags) != 1 {
		t.Fatalf("ToApp(root command with config flag) app flags length = %d, want 1", len(flags))
	}
	if flags := cmd.Flags; len(flags) != 1 {
		t.Fatalf("ToApp(root command with config flag) command flags length = %d, want 1", len(flags))
	}

	app, err = ToApp(cmd)
	if err != nil {
		t.Fatalf("second ToApp(root command with config flag) returned error: %v", err)
	}
	if flags := app.Flags; len(flags) != 1 {
		t.Fatalf("second ToApp(root command with config flag) app flags length = %d, want 1", len(flags))
	}
}

func TestToAppRootActionReadsParsedAppFlagValue(t *testing.T) {
	var gotConfig string
	cmd := New("root").Flags(&cli.StringFlag{
		Name:  "config",
		Value: "default-config",
	}).Action(func(c *cli.Context) error {
		gotConfig = c.String("config")
		return nil
	}).BuildBase()

	app, err := ToApp(cmd)
	if err != nil {
		t.Fatalf("ToApp(root command with config flag) returned error: %v", err)
	}

	if err := app.Run([]string{"test", "--config", "parsed-config"}); err != nil {
		t.Fatalf("ToApp(root command with config flag).Run(--config parsed-config) returned error: %v", err)
	}

	if gotConfig != "parsed-config" {
		t.Fatalf("ToApp(root command with config flag) root action config = %q, want %q", gotConfig, "parsed-config")
	}
}

func TestREADMEQuickStartAliasRunsSubcommand(t *testing.T) {
	args := []string{"example", "hi", "Alice"}
	var output bytes.Buffer
	cmd := New("example").
		Set.Usage("Small easycmd example").End.
		Child("hello").
		Set.Alias("hi").Usage("Print a greeting").End.
		Action(func(c *cli.Context) error {
			name := c.Args().First()
			if name == "" {
				name = "world"
			}
			_, err := fmt.Fprintf(c.App.Writer, "hello, %s\n", name)
			return err
		}).
		BuildBase()

	app, err := ToApp(cmd)
	if err != nil {
		t.Fatalf("ToApp(README quickstart command) got error %v, want nil", err)
	}
	app.Writer = &output

	if err := app.Run(args); err != nil {
		t.Fatalf("App.Run(%q) got error %v, want nil", args, err)
	}

	got := output.String()
	want := "hello, Alice\n"
	if got != want {
		t.Fatalf("App.Run(%q) output got %q, want %q", args, got, want)
	}
}
