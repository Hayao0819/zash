package shell

import (
	"fmt"
	"os"
	"strings"
)

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
