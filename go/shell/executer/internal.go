package executer

import (
	"os"

	"github.com/Hayao0819/zash/go/shell/builtin"
	"github.com/mattn/go-tty"
)

type InternalExecuter struct {
	Internal *builtin.InternalCmds
	Files    []*os.File
	TTY      *tty.TTY
}

func (ie *InternalExecuter) Exec(argv []string) error {
	if len(argv) == 0 {
		return nil
	}
	if ie.Files == nil {
		if ie.TTY != nil {
			ie.Files = filesFromTTY(ie.TTY)
		}
	}
	return ie.Internal.Run(argv[0], argv[1:], ie.Files).Error()
}
