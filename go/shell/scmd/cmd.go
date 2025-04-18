package scmd

import "github.com/mattn/go-tty"

type InternalCmd struct {
	Name string
	Func func(t *tty.TTY) func(args []string) error
	// State map[string]any
}

