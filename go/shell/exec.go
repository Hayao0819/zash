package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Hayao0819/nahi/futils"
)

func (s *Shell) FindInternalCmd(name string) *InternalCmdFunc {
	for _, cmd := range s.Internal {
		if cmd.Name == name {
			return &cmd
		}
	}
	return nil
}

func (s *Shell) Exec(argv []string) error {
	fmt.Println("Exec: ", argv)
	if len(strings.TrimSpace(strings.Join(argv, ""))) == 0 {
		return nil
	}
	if internalCmd := s.FindInternalCmd(argv[0]); internalCmd != nil {
		return internalCmd.Func(argv[1:])
	}
	// return s.ExecuteCmd(exec.Command(cmd, args...))
	return s.ExecuteExternalCmd(argv)
}

// func (s *Shell) ExecuteCmd(cmd *exec.Cmd) error {
// 	cmd.Stdin = s.TTY.Input()
// 	cmd.Stdout = s.TTY.Output()
// 	cmd.Stderr = s.TTY.Output()
// 	return cmd.Run()
// }

func execTargetAbsPath(fp string) (string, error) {
	if strings.Contains(fp, string(os.PathSeparator)) {
		if path.IsAbs(fp) {
			return fp, nil
		}
		return filepath.Abs(fp)
	}
	return exec.LookPath(fp)
}

func (s *Shell) ExecuteExternalCmd(argv []string) error {
	abs, err := execTargetAbsPath(argv[0])
	if err != nil {
		return fmt.Errorf("exec: %s: %w", argv[0], err)
	}

	if !futils.Exists(abs) {
		return fmt.Errorf("%s: No such file or directory", abs)
	}

	attr := &os.ProcAttr{
		Files: []*os.File{
			s.TTY.Input(),
			s.TTY.Output(),
			s.TTY.Output(),
		},
		Env: os.Environ(), // 現在の環境変数を継承
		Sys: &syscall.SysProcAttr{},
	}

	// プロセスを開始
	process, err := os.StartProcess(abs, argv, attr)
	if err != nil {
		return err
	}

	// 終了を待機
	state, err := process.Wait()
	if err != nil {
		return err
	}

	if state.Success() {
		return nil
	}
	return fmt.Errorf("%s: %s", abs, state.String())
}
