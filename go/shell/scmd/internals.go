package scmd

import (
	"github.com/mattn/go-tty"
)

type InternalCmds struct {
	Cmds []InternalCmd
	TTY  *tty.TTY
}

func NewInternalCmds(tty *tty.TTY) *InternalCmds {
	return &InternalCmds{
		Cmds: []InternalCmd{
			cdCmd,
			exitCmd,
			pwdCmd,
			// typeCmd,
		},
		TTY: tty,
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

func (ic *InternalCmds) Run(name string, args []string) error {
	cmd := ic.Get(name)
	if cmd == nil {
		return nil
	}
	return cmd.Func(ic.TTY)(args)
}
