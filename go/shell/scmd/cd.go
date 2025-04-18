package scmd

import (
	"fmt"
	"os"
)

var cdLastDir string

var cdCmd = InternalCmd{
	Name: "cd",
	Func: func(e Executer, args []string) Result {
		do := func(dir string) Result {
			cdLastDir, _ = os.Getwd()
			if err := os.Chdir(dir); err != nil {
				// return fmt.Errorf("cd: %s", err)
				return Result{
					exitcode: 1,
					err:      fmt.Errorf("cd: %s", err),
				}
			}
			return Result{
				exitcode: 0,
			}

		}

		if len(args) == 0 {
			// return os.Chdir(os.Getenv("HOME"))
			return do(os.Getenv("HOME"))
		} else if len(args) > 1 {
			// return fmt.Errorf("cd: too many arguments")
			return Result{
				exitcode: 1,
				err:      fmt.Errorf("cd: too many arguments"),
			}
		} else if args[0] == "-" {
			if cdLastDir == "" {
				// return fmt.Errorf("cd: no previous directory")
				return Result{
					exitcode: 1,
					err:      fmt.Errorf("cd: no previous directory"),
				}
			}

			return do(cdLastDir)
		} else {
			return do(args[0])
		}

	},
}
