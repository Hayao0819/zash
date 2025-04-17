package shell

import (
	"github.com/mattn/go-tty"
)

type Shell struct {
	TTY      *tty.TTY
	Internal []InternalCmdFunc
	started  bool
}

func New() (*Shell, error) {
	tty, err := tty.Open()
	if err != nil {
		return nil, err
	}
	return &Shell{
		TTY: tty,
		Internal: []InternalCmdFunc{
			cdCmd,
			exitCmd,
		},
	}, nil
}
