package shell

import (
	"fmt"
	"os"
	"strings"
)

func (s *Shell) Print(str string) {
	fmt.Fprintf(s.TTY.Output(), "%s", str)
}
func (s *Shell) Println(str string) {
	fmt.Fprintln(s.TTY.Output(), str)
}

func (s *Shell) PromptStr() string {
	promptStr := "Wanya?"

	if s.lastExitCode != 0 {
		promptStr = fmt.Sprintf("%s [exit %d]", promptStr, s.lastExitCode)
	}

	currentDir, err := os.Getwd()
	if err == nil {
		currentDir = strings.Replace(currentDir, os.Getenv("HOME"), "~", 1)
		promptStr = fmt.Sprintf("%s | %s", promptStr, currentDir)
	}
	promptStr = fmt.Sprintf("%s > ", promptStr)
	return promptStr
}

func (s *Shell) WaitInputWithPrompt() string {
	s.Print(s.PromptStr())
	r, _ := s.TTY.ReadString()

	return r
}

func (s *Shell) Close() {
	s.TTY.Close()
}
