package shell

import (
	"context"

	"github.com/Hayao0819/zash/go/prompt"
	"github.com/chzyer/readline"
)

func (s *Shell) WaitInputWithPrompt(ctx context.Context) string {

	p, err := prompt.New()
	if err != nil {
		s.Println(err.Error())
		return ""
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt: p.String(),
		Stdin:  s.TTY.Input(),
		Stdout: s.TTY.Output(),
		Stderr: s.TTY.Output(),
	})
	if err != nil {
		s.Println(err.Error())
		return ""
	}
	defer rl.Close()

	r, _ := rl.Readline()
	return r
}
