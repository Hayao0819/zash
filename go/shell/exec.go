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

// IsInternalCmd は指定されたコマンドが内部コマンドかどうかを判定します。
func (s *Shell) IsInternalCmd(cmd string) bool {
	return s.Internal.Get(cmd) != nil
}

// Exec はコマンド引数を受け取り、内部または外部コマンドとして実行します。
func (s *Shell) Exec(argv []string) error {
	// 空コマンドは無視
	if len(argv) == 0 || strings.TrimSpace(strings.Join(argv, "")) == "" {
		return nil
	}

	// 内部コマンドなら内部処理へ
	if s.IsInternalCmd(argv[0]) {
		return s.execInternal(argv)
	}

	// それ以外は外部コマンドとして実行
	return s.execExternal(argv)
}

// execInternal は内部コマンドを実行し、エラーを返します。
func (s *Shell) execInternal(argv []string) error {
	return s.Internal.Run(argv[0], argv[1:]).Error()
}

// execExternal は外部コマンドを実行し、プロセスの終了状態に応じて結果を返します。
func (s *Shell) execExternal(argv []string) error {
	// 実行ファイルの絶対パスを解決
	abs, err := resolveExecPath(argv[0])
	if err != nil {
		return fmt.Errorf("exec: %s: %w", argv[0], err)
	}

	// 実行ファイルが存在しない場合はエラー
	if !futils.Exists(abs) {
		return fmt.Errorf("%s: No such file or directory", abs)
	}

	// プロセス属性の設定（標準入出力と環境変数を継承）
	attr := &os.ProcAttr{
		Files: []*os.File{
			s.TTY.Input(),  // 標準入力
			s.TTY.Output(), // 標準出力
			s.TTY.Output(), // 標準エラー出力
		},
		Env: os.Environ(),           // 現在の環境変数を使用
		Sys: &syscall.SysProcAttr{}, // OS 固有の設定（今は空）
	}

	// プロセスの開始
	process, err := os.StartProcess(abs, argv, attr)
	if err != nil {
		return err
	}

	// プロセスの終了を待機
	state, err := process.Wait()
	if err != nil {
		return err
	}

	// 終了コードをプロンプトにセット
	s.prompt.SetExitCode(state.ExitCode())

	// 正常終了かどうかで結果を分岐
	if state.Success() {
		return nil
	}
	return fmt.Errorf("%s: %s", abs, state.String())
}

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
