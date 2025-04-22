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
	"github.com/Hayao0819/zash/go/prompt"
	"github.com/Hayao0819/zash/go/shell/scmd"
	"github.com/mattn/go-tty"
)

// resolveExecPath は与えられたコマンド名から実行可能な絶対パスを解決します。
func resolveExecPath(cmd string) (string, error) {
	if strings.Contains(cmd, string(os.PathSeparator)) {
		// パス区切り文字を含んでいる場合（例: ./foo、/bin/ls）
		if path.IsAbs(cmd) {
			return cmd, nil
		}
		return filepath.Abs(cmd)
	}
	// 環境変数 PATH を使って検索
	return exec.LookPath(cmd)
}

type Executer interface {
	Exec(argv []string) error
}

type InternalExecuter struct {
	Internal *scmd.InternalCmds
}

// IsInternalCmd は指定されたコマンドが内部コマンドかどうかを判定します。
func (s *Shell) IsInternalCmd(cmd string) bool {
	return s.Internal.Get(cmd) != nil
}

func (ie *InternalExecuter) Exec(argv []string) error {
	if len(argv) == 0 {
		return nil
	}
	return ie.Internal.Run(argv[0], argv[1:]).Error()
}

type ExternalExecuter struct {
	TTY    *tty.TTY
	Prompt *prompt.Prompt
}

func (ee *ExternalExecuter) Exec(argv []string) error {
	if len(argv) == 0 {
		return nil
	}

	abs, err := resolveExecPath(argv[0])
	if err != nil {
		return fmt.Errorf("exec: %s: %w", argv[0], err)
	}

	if !futils.Exists(abs) {
		return fmt.Errorf("%s: No such file or directory", abs)
	}

	attr := &os.ProcAttr{
		Files: []*os.File{
			ee.TTY.Input(),
			ee.TTY.Output(),
			ee.TTY.Output(),
		},
		Env: os.Environ(),
		Sys: &syscall.SysProcAttr{},
	}

	process, err := os.StartProcess(abs, argv, attr)
	if err != nil {
		return err
	}

	state, err := process.Wait()
	if err != nil {
		return err
	}

	ee.Prompt.SetExitCode(state.ExitCode())
	if state.Success() {
		return nil
	}

	return fmt.Errorf("%s: %s", abs, state.String())
}

func (s *Shell) Exec(argv []string) error {
	if len(argv) == 0 || strings.TrimSpace(strings.Join(argv, "")) == "" {
		return nil
	}

	var executer Executer
	if s.IsInternalCmd(argv[0]) {
		executer = &InternalExecuter{Internal: s.Internal}
	} else {
		executer = &ExternalExecuter{
			TTY:    s.TTY,
			Prompt: s.prompt,
		}
	}

	return executer.Exec(argv)
}
