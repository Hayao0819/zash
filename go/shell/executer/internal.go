package executer

import (
	"log/slog"
	"os"

	"github.com/Hayao0819/zash/go/shell/builtin"
	"github.com/mattn/go-tty"
)

type InternalExecuter struct {
	Internal *builtin.InternalCmds
	Files    []*os.File
	TTY      *tty.TTY
}

func (ie *InternalExecuter) Exec(argv []string) (int, error) {
	if len(argv) == 0 {
		return 0, nil
	}
	if ie.Files == nil {
		if ie.TTY != nil {
			ie.Files = filesFromTTY(ie.TTY)
		}
	}
	r := ie.Internal.Run(argv[0], argv[1:], ie.Files)
	// ie.Prompt.SetExitCode(r.ExitCode())
	slog.Debug("internal command", "command", argv[0], "exit code", r.ExitCode())
	return r.ExitCode(), r.Error()
}
