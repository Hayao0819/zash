package executer

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

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

func filesFromTTY(tty *tty.TTY) []*os.File {
	files := make([]*os.File, 0, 3)
	files = append(files, tty.Input(), tty.Output(), tty.Output())
	return files
}
