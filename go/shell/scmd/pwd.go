package scmd

import (
	"fmt"
	"os"
)

var pwdCmd = InternalCmd{
	Name: "pwd",
	Func: func(e Executer, args []string) Result {
		if len(args) > 0 {
			// return fmt.Errorf("exit: too many arguments")
			return Result{
				err:      fmt.Errorf("pwd: too many arguments"),
				exitcode: 1,
			}
		}

		wd, err := os.Getwd()
		if err != nil {
			// return fmt.Errorf("pwd: %s", err)
			return Result{
				err:      fmt.Errorf("pwd: %s", err),
				exitcode: 1,
			}
		}
		// fmt.Fprintln(e.TTY.Output(), wd)
		e.Puts(wd, "\n")

		return Result{
			err:      nil,
			exitcode: 0,
		}

	},
}

func init() {
	internalCmds = append(internalCmds, pwdCmd)
}
