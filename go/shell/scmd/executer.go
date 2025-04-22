package scmd

import (
	"github.com/mattn/go-tty"
)

type Executer struct {
	TTY *tty.TTY
}

func (e *Executer) Puts(line ...string) {
	// e.TTY.Output().Write([]byte(line))
	// fmt.Fprintln(e.TTY.Output(), line)

	for _, l := range line {
		e.TTY.Output().Write([]byte(l))
	}
}
