package shell

import (
	"context"

	"github.com/Hayao0819/zash/go/prompt"
	"github.com/Hayao0819/zash/go/shell/scmd"
	"github.com/mattn/go-tty"
)

type Shell struct {
	TTY      *tty.TTY
	Internal *scmd.InternalCmds
	started  bool
	prompt   *prompt.Prompt
}

func (s *Shell) Context() context.Context {
	return context.TODO()
}

func New() (*Shell, error) {
	t, err := tty.Open()
	if err != nil {
		return nil, err
	}

	p, err := prompt.New(t)
	if err != nil {
		return nil, err
	}

	s := Shell{
		TTY:      t,
		prompt:   p,
		started:  false,
		Internal: scmd.NewInternalCmds(t),
	}

	return &s, nil
}
