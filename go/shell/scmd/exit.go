package scmd

import (
	"fmt"
	"os"
)

var exitCmd = InternalCmd{
	Name: "exit",
	Func: func(e Executer, args []string) Result {
		if len(args) > 0 {
			return Result{
				err:      fmt.Errorf("exit: too many arguments"),
				exitcode: 1,
			}
		}
		os.Exit(0)
		return Result{
			exitcode: 0,
		}
	},
}

func init() {
	internalCmds = append(internalCmds, exitCmd)
}
