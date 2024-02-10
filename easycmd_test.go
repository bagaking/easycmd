package easycmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
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

func TestBuilderChildReusesExistingAlias(t *testing.T) {
	builder := New("root")

	byName := builder.Child("hello").Set.Alias("hi").Usage("Print a greeting").End
	byAlias := builder.Child("hi")
	byAlias.Set.Usage("Print a warm greeting")

	if byAlias.BuildCur() != byName.BuildCur() {
		t.Fatalf("Child(%q) returned command %p, want existing alias command %p", "hi", byAlias.BuildCur(), byName.BuildCur())
	}

	subcommands := builder.BuildBase().Subcommands
	if len(subcommands) != 1 {
		t.Fatalf("Child(%q) created %d subcommands, want 1: %#v", "hi", len(subcommands), subcommands)
	}
	if usage := subcommands[0].Usage; usage != "Print a warm greeting" {
		t.Fatalf("Child(%q) reused command usage = %q, want %q", "hi", usage, "Print a warm greeting")
	}
}

func TestToAppDoesNotClearCommandFlags(t *testing.T) {
	restoreHelpPrinterAfter(t)

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
	restoreHelpPrinterAfter(t)

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

func TestToAppPreservesCustomHelpPrinter(t *testing.T) {
	restoreHelpPrinterAfter(t)

	tests := []struct {
		name        string
		helpPrinter func(io.Writer, string, interface{})
	}{
		{
			name:        "custom help printer",
			helpPrinter: func(out io.Writer, tpl string, data interface{}) {},
		},
		{
			name:        "alternate custom help printer",
			helpPrinter: func(out io.Writer, tpl string, data interface{}) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli.HelpPrinter = tt.helpPrinter
			want := reflect.ValueOf(tt.helpPrinter).Pointer()

			if _, err := ToApp(New("root").BuildBase()); err != nil {
				t.Fatalf("ToApp(root command) returned error: %v", err)
			}

			got := reflect.ValueOf(cli.HelpPrinter).Pointer()
			if got != want {
				t.Fatalf("ToApp(root command) HelpPrinter pointer = %v, want %v", got, want)
			}
		})
	}
}

func TestREADMEQuickStartAliasRunsSubcommand(t *testing.T) {
	restoreHelpPrinterAfter(t)

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

	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "primary command",
			args: []string{"example", "hello", "Alice"},
			want: "hello, Alice\n",
		},
		{
			name: "alias",
			args: []string{"example", "hi", "Alice"},
			want: "hello, Alice\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			app.Writer = &output

			if err := app.Run(tt.args); err != nil {
				t.Fatalf("App.Run(%q) got error %v, want nil", tt.args, err)
			}

			got := output.String()
			if got != tt.want {
				t.Fatalf("App.Run(%q) output got %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

func TestNestedAliasHelpUsesCustomHelpPrinter(t *testing.T) {
	restoreHelpPrinterAfter(t)

	cmd := New("example").
		Child("admin").
		Set.Alias("adm").Usage("Manage resources").End.
		Child("users").
		Set.Alias("u").Usage("Manage users").End.
		BuildBase()

	app, err := ToApp(cmd)
	if err != nil {
		t.Fatalf("ToApp(nested alias command) got error %v, want nil", err)
	}

	type helpCall struct {
		template string
		name     string
		aliases  []string
	}

	var calls []helpCall
	cli.HelpPrinter = func(out io.Writer, tpl string, data interface{}) {
		command, ok := data.(*cli.Command)
		if !ok {
			t.Errorf("HelpPrinter data type = %T, want *cli.Command", data)
			return
		}
		calls = append(calls, helpCall{
			template: tpl,
			name:     command.Name,
			aliases:  append([]string(nil), command.Aliases...),
		})
		_, _ = fmt.Fprintf(out, "custom help for %s\n", command.Name)
	}

	var output bytes.Buffer
	app.Writer = &output
	if err := app.Run([]string{"example", "adm", "u", "--help"}); err != nil {
		t.Fatalf("App.Run(nested alias help) got error %v, want nil", err)
	}

	if got, want := output.String(), "custom help for users\n"; got != want {
		t.Fatalf("App.Run(nested alias help) output got %q, want %q", got, want)
	}
	if len(calls) != 1 {
		t.Fatalf("HelpPrinter call count = %d, want 1: %#v", len(calls), calls)
	}
	gotCall := calls[0]
	if gotCall.template != cli.CommandHelpTemplate {
		t.Fatalf("HelpPrinter template = %q, want command help template", gotCall.template)
	}
	if gotCall.name != "users" {
		t.Fatalf("HelpPrinter command name = %q, want %q", gotCall.name, "users")
	}
	if !reflect.DeepEqual(gotCall.aliases, []string{"u"}) {
		t.Fatalf("HelpPrinter command aliases = %#v, want %#v", gotCall.aliases, []string{"u"})
	}
}

func restoreHelpPrinterAfter(t *testing.T) {
	t.Helper()

	helpPrinter := cli.HelpPrinter
	t.Cleanup(func() {
		cli.HelpPrinter = helpPrinter
	})
}
