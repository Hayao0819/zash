package shell

import (
	"github.com/Hayao0819/zash/go/shell/scmd"
	"github.com/mattn/go-tty"
)

type Shell struct {
	TTY      *tty.TTY
	Internal *scmd.InternalCmds
	started  bool
}

func New() (*Shell, error) {
	tty, err := tty.Open()
	if err != nil {
		return nil, err
	}
	return &Shell{
		TTY:      tty,
		started:  false,
		Internal: scmd.NewInternalCmds(tty),
	}, nil
}
