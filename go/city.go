package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mattn/go-tty"
)

type Shell struct {
	TTY     *tty.TTY
	started bool
}

func NewShell() (*Shell, error) {
	tty, err := tty.Open()
	if err != nil {
		return nil, err
	}
	return &Shell{TTY: tty}, nil
}

func (s *Shell) Puts(str string) {
	fmt.Fprintf(s.TTY.Output(), "%s", str)
}

func (s *Shell) Prompt() string {
	s.Puts("Wanya?> ")
	r, _ := s.TTY.ReadString()
	return r
}

func (s *Shell) Close() {
	s.TTY.Close()
}

func (s *Shell) Exec(entered string) {
	switch entered {
	case "exit":
		s.Close()
		os.Exit(0)
	default:
		if err := s.ExecuteCommand(exec.Command(entered)); err != nil {
			fmt.Println("Err: ", err)
		}
	}
}

func (s *Shell) ExecuteCommand(cmd *exec.Cmd) error {
	cmd.Stdin = s.TTY.Input()
	cmd.Stdout = s.TTY.Output()
	cmd.Stderr = s.TTY.Output()
	return cmd.Run()
}

func (s *Shell) StartInteractive() {
	if s.started {
		return
	}
	s.started = true

	for {
		input := s.Prompt()
		s.Exec(input)
	}
}

func main() {
	fmt.Println("Welcome to Zash!")
	shell, err := NewShell()
	handleErr(err)
	shell.StartInteractive()
}
