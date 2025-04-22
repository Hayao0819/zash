package prompt

import (
	"github.com/chzyer/readline"
	"github.com/mattn/go-tty"
)

type Prompt struct {
	user       string
	currentDir string
	exitCode   int
	tty        *tty.TTY
}

func New(t *tty.TTY) (*Prompt, error) {
	p := Prompt{
		tty: t,
	}
	if err := p.Update(); err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *Prompt) WaitInput() (string, error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt: p.String(),
		Stdin:  p.tty.Input(),
		Stdout: p.tty.Output(),
		Stderr: p.tty.Output(),
	})
	if err != nil {
		return "", err
	}
	defer rl.Close()

	line, err := rl.Readline()
	if err != nil {
		return "", err
	}
	if line == "" {
		return "", nil
	}
	return line, nil
}
