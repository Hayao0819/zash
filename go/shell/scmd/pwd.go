package scmd

import (
	"fmt"
	"os"

	"github.com/mattn/go-tty"
)

var pwdCmd = InternalCmd{
	Name: "pwd",
	Func: func(t *tty.TTY) func(args []string) error {
		return func(args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("exit: too many arguments")
			}

			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("pwd: %s", err)
			}
			fmt.Fprintln(t.Output(), wd)
			return nil
		}
	},
}
