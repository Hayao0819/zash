package executer

import (
	"fmt"
	"os"
	"syscall"

	"github.com/Hayao0819/nahi/futils"
	"github.com/Hayao0819/zash/go/prompt"
	"github.com/mattn/go-tty"
)

type ExternalExecuter struct {
	TTY    *tty.TTY
	Prompt *prompt.Prompt
	Files  []*os.File
}

func (ee *ExternalExecuter) Exec(argv []string) (int, error) {
	if len(argv) == 0 {
		return 0, nil
	}

	abs, err := resolveExecPath(argv[0])
	if err != nil {
		return 127, fmt.Errorf("exec: %s: %w", argv[0], err)
	}

	if !futils.Exists(abs) {
		return 127, fmt.Errorf("%s: No such file or directory", abs)
	}

	files := ee.Files
	if len(files) == 0 {
		files = []*os.File{
			ee.TTY.Input(),
			ee.TTY.Output(),
			ee.TTY.Output(),
		}
	}
	attr := &os.ProcAttr{
		Files: files,
		Env:   os.Environ(),
		Sys:   &syscall.SysProcAttr{},
	}

	process, err := os.StartProcess(abs, argv, attr)
	if err != nil {
		return 127, err
	}

	state, err := process.Wait()
	if err != nil {
		return 127, err
	}

	ee.Prompt.SetExitCode(state.ExitCode())
	if state.Success() {
		return 0, nil
	}

	return state.ExitCode(), fmt.Errorf("%s: %s", abs, state.String())
}
