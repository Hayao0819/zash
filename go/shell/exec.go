package shell

import (
	"log/slog"

	"github.com/Hayao0819/zash/go/ast"
	"github.com/Hayao0819/zash/go/shell/builtin"
	"github.com/Hayao0819/zash/go/shell/executer"
)

// IsInternalCmd は指定されたコマンドが内部コマンドかどうかを判定します。
func (s *Shell) IsInternalCmd(cmd string) bool {
	return builtin.Cmds.Get(cmd) != nil
}

func (s *Shell) getExecuter(name string) executer.Executer {
	if s.IsInternalCmd(name) {
		return &executer.InternalExecuter{
			Internal: &builtin.Cmds,
			TTY:      s.TTY,
		}
	} else {
		return &executer.ExternalExecuter{
			TTY: s.TTY,
		}
	}
}

func (s *Shell) Exec(cmd *ast.Command) (int, error) {
	if cmd == nil || cmd.Name == "" {
		return 0, nil
	}
	ex := s.getExecuter(cmd.Name)
	if len(cmd.Suffix.Redirections) != 0 {
		{
			slog.Debug("ShellParsedCommand", "name", cmd.Name, "args", cmd.Suffix.Args)
			for _, r := range cmd.Suffix.Redirections {
				slog.Debug("redirection", "operator", r.Operator, "file", r.File)
			}
		}
		ex = &executer.RedirectedExecuter{
			Base: ex,
		}
	}

	argv := []string{cmd.Name}
	if len(cmd.Suffix.Args) != 0 {
		argv = append(argv, cmd.Suffix.Args...)
	}

	ec, err := ex.Exec(argv)
	s.prompt.SetExitCode(ec)
	return ec, err
}
