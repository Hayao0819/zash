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
		ex = &executer.InternalExecuter{}
	} else {
		ex = &executer.ExternalExecuter{}
	}

	return ex
}

func (s *Shell) defaultIOContext() executer.IOContext {
	return executer.IOContext{
		Stdin:  s.TTY.Input(),
		Stdout: s.TTY.Output(),
		Stderr: s.TTY.Output(),
	}
}

func (s *Shell) Exec(cmd *ast.Command) (int, error) {
	if cmd == nil || cmd.Name == "" {
		return 0, nil
	}

	ioctx := s.defaultIOContext()

	// --- ここを追加 (リダイレクト適用) ---
	if len(cmd.Suffix.Redirections) != 0 {
		logmgr.Shell().Debug("ShellParsedCommand", "name", cmd.Name, "args", cmd.Suffix.Args)
		for _, r := range cmd.Suffix.Redirections {
			logmgr.Shell().Debug("ShellExecRedirection", "operator", r.Operator, "file", r.File)
			if err := ioctx.Redirect(r.Operator, r.File); err != nil {
				return 1, err
			}
		}
	}

	ex := s.getExecuter(cmd)
	ec, err := ex.Exec(cmd.Argv(), ioctx)
	s.prompt.SetExitCode(ec)
	return ec, err
}
