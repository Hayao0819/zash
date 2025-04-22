package shell

import (
	"strings"

	"github.com/Hayao0819/zash/go/shell/executer"
)

// IsInternalCmd は指定されたコマンドが内部コマンドかどうかを判定します。
func (s *Shell) IsInternalCmd(cmd string) bool {
	return s.Internal.Get(cmd) != nil
}

func (s *Shell) Exec(argv []string) error {
	if len(argv) == 0 || strings.TrimSpace(strings.Join(argv, "")) == "" {
		return nil
	}

	var ex executer.Executer
	if s.IsInternalCmd(argv[0]) {
		ex = &executer.InternalExecuter{Internal: s.Internal}
	} else {
		ex = &executer.ExternalExecuter{
			TTY:    s.TTY,
			Prompt: s.prompt,
		}
	}

	return ex.Exec(argv)
}
