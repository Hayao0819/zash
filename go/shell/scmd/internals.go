package scmd

import (
	"fmt"

	"github.com/mattn/go-tty"
)

type InternalCmds struct {
	Cmds     []InternalCmd
	Executer Executer
}

var internalCmds = []InternalCmd{}

func NewInternalCmds(tty *tty.TTY) *InternalCmds {
	return &InternalCmds{
		// Cmds: []InternalCmd{
		// 	cdCmd,
		// 	exitCmd,
		// 	pwdCmd,
		// 	// typeCmd,
		// },
		Cmds: internalCmds,
		Executer: Executer{
			TTY: tty,
		},
	}
}

func (ic *InternalCmds) Get(name string) *InternalCmd {
	for _, cmd := range ic.Cmds {
		if cmd.Name == name {
			return &cmd
		}
	}
	return nil
}

func (ic *InternalCmds) Run(name string, args []string) Result {
	cmd := ic.Get(name)
	if cmd == nil {
		return Result{
			err:      fmt.Errorf("%s: command not found", name),
			exitcode: 127,
		}
	}
	return cmd.Func(ic.Executer, args)

}
