package easycmd

import (
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

var ErrFlagAlreadyExist = errors.New("flag already exists")

func Print(cliCtx *cli.Context, msgLst ...interface{}) (err error) {
	_, err = cliCtx.App.Writer.Write([]byte(fmt.Sprint(msgLst...)))
	return
}

func Println(cliCtx *cli.Context, msgLst ...interface{}) (err error) {
	_, err = cliCtx.App.Writer.Write([]byte(fmt.Sprintln(msgLst...)))
	return
}
