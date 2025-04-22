package shell

import (
	"context"
)

func (s *Shell) WaitInputWithPrompt(ctx context.Context) string {
	r, err := s.prompt.WaitInput()
	if err != nil {
		s.Println(err.Error())
		return ""
	}

	return r
}
