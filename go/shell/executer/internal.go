package executer

import "github.com/Hayao0819/zash/go/shell/scmd"

type InternalExecuter struct {
	Internal *scmd.InternalCmds
}

func (ie *InternalExecuter) Exec(argv []string) error {
	if len(argv) == 0 {
		return nil
	}
	return ie.Internal.Run(argv[0], argv[1:]).Error()
}
