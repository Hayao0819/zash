package shell

import (
	"context"

	"github.com/Hayao0819/zash/go/shell/scmd"
	"github.com/mattn/go-tty"
)

type Shell struct {
	TTY          *tty.TTY
	Internal     *scmd.InternalCmds
	started      bool
	lastExitCode int
}

func (s *Shell) Context() context.Context {
	return context.TODO()
}

func New() (*Shell, error) {
	t, err := tty.Open()
	if err != nil {
		return nil, err
	}
	return &Shell{
		TTY:      t,
		started:  false,
		Internal: scmd.NewInternalCmds(t),
	}, nil
}
