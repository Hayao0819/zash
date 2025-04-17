package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Hayao0819/zash/go/lexer"
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

func (s *Shell) FindInternalCmd(name string) *InternalCmdFunc {
	for _, cmd := range s.Internal {
		if cmd.Name == name {
			return &cmd
		}
	}
	return nil
}

func (s *Shell) Puts(str string) {
	fmt.Fprintf(s.TTY.Output(), "%s", str)
}

func (s *Shell) PromptStr() string {
	promptStr := "Wanya?"

	currentDir, err := os.Getwd()
	if err == nil {
		currentDir = strings.Replace(currentDir, os.Getenv("HOME"), "~", 1)
		promptStr = fmt.Sprintf("%s | %s", promptStr, currentDir)
	}
	promptStr = fmt.Sprintf("%s > ", promptStr)
	return promptStr
}

func (s *Shell) WaitInputWithPrompt() string {
	s.Puts(s.PromptStr())
	r, _ := s.TTY.ReadString()

	return r
}

func (s *Shell) Close() {
	s.TTY.Close()
}

func (s *Shell) Exec(cmd string, args []string) {
	fmt.Println("Exec: ", cmd, args)
	if cmd == "" {
		return
	}
	if internalCmd := s.FindInternalCmd(cmd); internalCmd != nil {
		if err := internalCmd.Func(args); err != nil {
			fmt.Println("Err: ", err)
		}
		return
	}
	if err := s.ExecuteCmd(exec.Command(cmd, args...)); err != nil {
		fmt.Println("Err: ", err)
	}

}

func (s *Shell) ExecuteCmd(cmd *exec.Cmd) error {
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
		input := s.WaitInputWithPrompt()


		tokens, err := lexer.NewLexer(input).ReadAll()
		if err != nil {
			fmt.Println("Err: ", err)
			continue
		}
		if len(tokens) == 0 {
			continue
		}
		
		args:= strings.Split(strings.Join(tokens[1:], ""), " ")

		s.Exec(tokens[0], args)
	}
}
