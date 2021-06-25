package easycmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
)

type CustomOption struct {
	ExitAfterPrintHelpMsg bool
}

// SetCustomOptions can set some easycmd-layer options
func SetCustomOptions(opt CustomOption) {
	if opt.ExitAfterPrintHelpMsg {
		cli.HelpPrinter = func(out io.Writer, tpl string, data interface{}) {
			cli.HelpPrinterCustom(out, tpl, data, nil)
			os.Exit(0)
		}
	}
}

// SetCustomOptions can merge two list of cli.Flags
func MergeFlags(flags1, flags2 []cli.Flag, ignoreError ...bool) ([]cli.Flag, error) {
	all := append(flags1, flags2...)
	ier := ignoreError != nil && len(ignoreError) > 0 && ignoreError[0]

	result := make([]cli.Flag, len(all))

	exist := map[string]bool{}
	for _, v := range all {
		conflictKey, keys := "", v.Names()

		for _, key := range keys {
			if exist[key] {
				conflictKey = key
				break
			}
		}

		if conflictKey == "" {
			for _, key := range keys {
				exist[key] = true
			}
			result = append(result, v)
			continue
		}

		if !ier {
			return nil, fmt.Errorf("%w, key= %v", ErrFlagAlreadyExist, conflictKey)
		}
	}

	return result, nil
}

// ToApp creates an app from a single command
// when flags conflict, error will be thrown and the procedure will be stock
// to ignore the error, you can stash the flags of the command, and the call MergeFlags by yourself
func ToApp(cmd *cli.Command) (*cli.App, error) {
	SetCustomOptions(CustomOption{ExitAfterPrintHelpMsg: true})

	app := &cli.App{
		Name:        filepath.Base(os.Args[0]),
		HelpName:    filepath.Base(os.Args[0]),
		Description: cmd.Usage,

		BashComplete: cli.DefaultAppComplete,
		Action:       cmd.Run,
		Reader:       os.Stdin,
		Writer:       os.Stdout,
		ErrWriter:    os.Stderr,
		Compiled: func() time.Time {
			info, err := os.Stat(os.Args[0])
			if err != nil {
				return time.Now()
			}
			return info.ModTime()
		}(),
	}

	if cmd.Flags != nil {
		if app.Flags == nil {
			app.Flags = cmd.Flags
		} else {
			flags, err := MergeFlags(app.Flags, cmd.Flags)
			if err != nil {
				return nil, err
			}

			app.Flags = flags
		}
		cmd.Flags = nil
	}

	app.Commands = cmd.Subcommands

	return app, nil
}
