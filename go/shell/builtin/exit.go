package builtin

import (
	"fmt"
	"os"
)

var exitCmd = internalCmd{
	Name: "exit",
	Func: func(args []string, files []*os.File) Result {
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
	Cmds = append(Cmds, exitCmd)
}
