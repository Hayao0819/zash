package scmd

import (
	"fmt"
	"os"

	"github.com/mattn/go-tty"
)

var cdLastDir string

var cdCmd = InternalCmd{
	Name: "cd",
	Func: func(t *tty.TTY) func(args []string) error {
		return func(args []string) error {
			do := func(dir string) error {
				cdLastDir, _ = os.Getwd()
				if err := os.Chdir(dir); err != nil {
					return fmt.Errorf("cd: %s", err)
				}
				return nil
			}

			if len(args) == 0 {
				// return os.Chdir(os.Getenv("HOME"))
				return do(os.Getenv("HOME"))
			} else if len(args) > 1 {
				return fmt.Errorf("cd: too many arguments")
			} else if args[0] == "-" {
				if cdLastDir == "" {
					return fmt.Errorf("cd: no previous directory")
				}

				return do(cdLastDir)
			} else {
				return do(args[0])
			}
		}
	},
}
