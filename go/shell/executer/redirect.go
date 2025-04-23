package executer

import (
	"fmt"
	"os"

	"github.com/Hayao0819/zash/go/ast"
)

type RedirectedExecuter struct {
	Base         Executer
	Redirections []*ast.Redirection
}

func (r *RedirectedExecuter) Exec(argv []string) (int, error) {
	// 初期値は現在の TTY
	files := []*os.File{
		os.Stdin,  // FD 0
		os.Stdout, // FD 1
		os.Stderr, // FD 2
	}

	// リダイレクト処理
	for _, redir := range r.Redirections {
		switch redir.Operator {
		case ">":
			f, err := os.Create(redir.File)
			if err != nil {
				return 127, fmt.Errorf("redirect: %s: %w", redir.File, err)
			}
			defer f.Close()
			files[1] = f

		case "<":
			f, err := os.Open(redir.File)
			if err != nil {
				return 127, fmt.Errorf("redirect: %s: %w", redir.File, err)
			}
			defer f.Close()
			files[0] = f
		}
	}

	// 外部コマンドなら `ExternalExecuter` を想定してリダイレクト渡す
	// if ext, ok := r.Base.(*ExternalExecuter); ok {
	// 	ext.Files = files
	// }else if intExt, ok := r.Base.(*InternalExecuter); ok {

	switch base := r.Base.(type) {
	case *ExternalExecuter:
		base.Files = files
	case *InternalExecuter:
		base.Files = files
	default:
		return 127, fmt.Errorf("unsupported executer type: %T", base)
	}

	return r.Base.Exec(argv)
}
