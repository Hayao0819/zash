package scmd

import (
	"fmt"
	"os"

	"github.com/mattn/go-tty"
)

var exitCmd = InternalCmd{
	Name: "exit",
	Func: func(t *tty.TTY) func(args []string) error {
		return func(args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("exit: too many arguments")
			}
			os.Exit(0)
			return nil
		}
	},
}
