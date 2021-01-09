package easycmd

import (
	"io"
	"os"

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
