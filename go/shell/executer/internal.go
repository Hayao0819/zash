package executer

import (
	"log/slog"
	"os"

	"github.com/Hayao0819/zash/go/shell/builtin"
)

type InternalExecuter struct {
	Files []*os.File
	// TTY      *tty.TTY
}

func (ie *InternalExecuter) Exec(argv []string, ioctx IOContext) (int, error) {
	if len(argv) == 0 {
		return 0, nil
	}
	r := builtin.Cmds.Run(argv[0], argv[1:], ioctx.Files())
	slog.Debug("internal command", "command", argv[0], "exit code", r.ExitCode())
	return r.ExitCode(), r.Error()
}
