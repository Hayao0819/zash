package shell

import (
	"fmt"
	"os"
)

type InternalCmdFunc struct {
	Name string
	Func func(args []string) error
}

var cdCmd = InternalCmdFunc{
	Name: "cd",
	Func: func(args []string) error {
		if len(args) == 0 {
			return os.Chdir(os.Getenv("HOME"))
		}
		if len(args) > 1 {
			return fmt.Errorf("cd: too many arguments")
		}
		if err := os.Chdir(args[0]); err != nil {
			return fmt.Errorf("cd: %s", err)
		}
		return nil
	},
}

var exitCmd = InternalCmdFunc{
	Name: "exit",
	Func: func(args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("exit: too many arguments")
		}
		os.Exit(0)
		return nil
	},
}
