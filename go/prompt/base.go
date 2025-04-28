package prompt

import (
	"github.com/chzyer/readline"
	"github.com/mattn/go-tty"
)

type Prompt struct {
	user        string
	currentDir  string
	exitCode    int
	tty         *tty.TTY
	historyFile string
}

func New(t *tty.TTY, historyFile string) (*Prompt, error) {
	p := Prompt{
		tty: t,
		// historyFile: "./history.txt",
	}
	if err := p.Update(); err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *Prompt) NewReadLine() (*readline.Instance, error) {
	c := readline.Config{
		Prompt: p.String(),
		Stdin:  p.tty.Input(),
		Stdout: p.tty.Output(),
		Stderr: p.tty.Output(),
	}
	if p.historyFile != "" {
		c.HistoryFile = p.historyFile
		c.HistoryLimit = 1000
		c.AutoComplete = nil
	}
	rl, err := readline.NewEx(&c)
	if err != nil {
		return nil, err
	}

	return rl, nil
}

func (p *Prompt) WaitInput() (string, error) {
	rl, err := p.NewReadLine()
	if err != nil {
		return "", err
	}
	defer rl.Close()

	line, err := rl.Readline()
	if err != nil {
		return "", err
	}
	if err := rl.SaveHistory(line); err != nil {
		return "", err
	}

	if line == "" {
		return "", nil
	}
	return line, nil
}
