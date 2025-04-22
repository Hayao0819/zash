package shell

import (
	"context"

	"github.com/Hayao0819/zash/go/prompt"
)

func (s *Shell) WaitInputWithPrompt(ctx context.Context) string {
	p, err := prompt.New(s.TTY)
	if err != nil {
		s.Println(err.Error())
		return ""
	}

	r, err := p.WaitInput()
	if err != nil {
		s.Println(err.Error())
		return ""
	}

	return r
}
