package shell

import (
	"github.com/Hayao0819/zash/go/ast"
	"github.com/Hayao0819/zash/go/internal/logmgr"
	"github.com/Hayao0819/zash/go/shell/builtin"
	"github.com/Hayao0819/zash/go/shell/executer"
)

// 指定されたコマンドが内部コマンドかどうかを判定
func (s *Shell) IsInternalCmd(cmd string) bool {
	return builtin.Cmds.Get(cmd) != nil
}

// 指定されたコマンドに基づいて適切なExecuterを取得
func (sn *Shell) getExecuter(c *ast.Command) executer.Executer {
	var ex executer.Executer
	if sn.IsInternalCmd(c.Name) {
		ex = &executer.InternalExecuter{
			Internal: &builtin.Cmds,
			TTY:      sn.TTY,
		}
	} else {
		ex = &executer.ExternalExecuter{
			TTY: sn.TTY,
		}
	}

	if len(c.Suffix.Redirections) != 0 {
		ex = &executer.RedirectedExecuter{
			Base:         ex,
			Redirections: c.Suffix.Redirections,
		}
	}
	return ex
}

func (s *Shell) Exec(cmd *ast.Command) (int, error) {
	if cmd == nil || cmd.Name == "" {
		return 0, nil
	}
	ex := s.getExecuter(cmd)
	if len(cmd.Suffix.Redirections) != 0 {
		{
			logmgr.Shell().Debug("ShellParsedCommand", "name", cmd.Name, "args", cmd.Suffix.Args)
			for _, r := range cmd.Suffix.Redirections {
				logmgr.Shell().Debug("ShellExecRedirection", "operator", r.Operator, "file", r.File)
			}
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
