package shell

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/Hayao0819/zash/go/ast"
	"github.com/Hayao0819/zash/go/internal/logmgr"
	"github.com/Hayao0819/zash/go/shell/builtin"
	"github.com/Hayao0819/zash/go/shell/executer"
)

// ASTノード種別ごとに再帰的に実行する
func (s *Shell) ExecNode(node ast.Node) (int, error) {
	switch n := node.(type) {
	case *ast.Script:
		return s.ExecNode(n.List)
	case *ast.List:
		var last int
		for _, item := range n.Items {
			ec, err := s.ExecNode(item)
			if err != nil {
				return ec, err
			}
			last = ec
		}
		return last, nil
	case *ast.Pipeline:
		return s.execPipeline(n)
	case *ast.SimpleCommand:
		return s.Exec(n)
	case *ast.ShellCommand:
		logmgr.Shell().Debug("ShellExecShellCommand", "kind", n.Kind)
		switch n.Kind {
		case "if":
			ifc := n.Node.(*ast.IfCommand)
			cond, err := s.ExecNode(ifc.Cond)
			if err != nil {
				return cond, err
			}
			if cond == 0 {
				return s.ExecNode(ifc.Then)
			} else if ifc.Else != nil {
				return s.ExecNode(ifc.Else)
			}
			return 0, nil
		case "for":
			fc := n.Node.(*ast.ForCommand)
			var last int
			for range fc.Words {
				// 変数展開は未実装: 環境変数にセットする場合はここで
				var err error
				last, err = s.ExecNode(fc.Body)
				if err != nil {
					return last, err
				}
			}
			return last, nil
		case "while":
			wc := n.Node.(*ast.WhileCommand)
			var last int
			for {
				logmgr.Shell().Debug("ShellExecWhileCommand", "condition", wc.Cond)
				cond, err := s.ExecNode(wc.Cond)
				if err != nil {
					return cond, err
				}
				if cond != 0 { // 条件が偽（非0）の場合にループを終了
					break
				}
				last, err = s.ExecNode(wc.Body)
				if err != nil {
					return last, err
				}
			}
			return last, nil
		case "until":
			uc := n.Node.(*ast.UntilCommand)
			var last int
			for {
				cond, err := s.ExecNode(uc.Cond)
				if err != nil {
					return cond, err
				}
				if cond == 0 {
					break
				}
				last, err = s.ExecNode(uc.Body)
				if err != nil {
					return last, err
				}
			}
			return last, nil
		case "select":
			// select文はforと同様に処理（簡易）
			sc := n.Node.(*ast.SelectCommand)
			var last int
			for range sc.Words {
				var err error
				last, err = s.ExecNode(sc.Body)
				if err != nil {
					return last, err
				}
			}
			return last, nil
		case "case":
			// case文は未実装: 必要に応じてパターンマッチを追加
			return 0, nil
		case "function":
			// function定義は未実装: 必要に応じて関数テーブルに登録
			return 0, nil
		}
		return 0, nil
	case *ast.Subshell:
		// サブシェルは新しいShellインスタンスで実行するのが理想だが、ここでは再帰で代用
		return s.ExecNode(n.Body)
	case *ast.GroupCommand:
		return s.ExecNode(n.Body)
	case *ast.CompoundList:
		return s.ExecNode(n.List)
	default:
		return 0, fmt.Errorf("unsupported AST node: %v", reflect.TypeOf(node))
	}
}

// execPipeline はパイプライン（複数コマンドのパイプ接続）を実行する
func (s *Shell) execPipeline(p *ast.Pipeline) (int, error) {
	if len(p.Commands) == 0 {
		return 0, nil
	}

	// 単一コマンドの場合はそのまま実行
	if len(p.Commands) == 1 {
		ec, err := s.ExecNode(p.Commands[0].Cmd)
		if p.Bang {
			if ec == 0 {
				ec = 1
			} else {
				ec = 0
			}
		}
		return ec, err
	}

	// 複数コマンドのパイプライン
	n := len(p.Commands)
	ioctxs := make([]*executer.IOContext, n)

	// 各コマンドにIOContextを作成
	for i := 0; i < n; i++ {
		ctx := s.defaultIOContext()
		ioctxs[i] = &ctx
	}

	// PipeToで隣接コマンド間を接続
	for i := 0; i < n-1; i++ {
		if err := ioctxs[i].PipeTo(ioctxs[i+1]); err != nil {
			// エラー時はすでに作成したIOContextをクローズ
			for j := 0; j <= i; j++ {
				ioctxs[j].Close()
			}
			return 1, fmt.Errorf("pipe error: %w", err)
		}
	}

	// goroutineで並行実行
	type result struct {
		ec  int
		err error
	}
	results := make([]result, n)
	var wg sync.WaitGroup

	for i, pc := range p.Commands {
		wg.Add(1)
		go func(idx int, cmd ast.Node, ioctx *executer.IOContext) {
			defer wg.Done()
			ec, err := s.execSimpleCommandWithIO(cmd, ioctx)
			// コマンド完了後にパイプfdを閉じる
			// これにより下流コマンドがEOFを受け取れる
			ioctx.Close()
			results[idx] = result{ec, err}
		}(i, pc.Cmd, ioctxs[i])
	}

	wg.Wait()

	// 最後のコマンドの終了コードを返す
	last := results[n-1]
	ec := last.ec
	if p.Bang {
		if ec == 0 {
			ec = 1
		} else {
			ec = 0
		}
	}
	return ec, last.err
}

// execSimpleCommandWithIO は指定されたIOContextでコマンドを実行する
func (s *Shell) execSimpleCommandWithIO(node ast.Node, ioctx *executer.IOContext) (int, error) {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		// ShellCommand等の場合は通常のExecNodeで実行
		return s.ExecNode(node)
	}

	if cmd == nil || len(cmd.Elements) == 0 {
		return 0, nil
	}

	// コマンド名・引数・リダイレクト抽出
	var name string
	var args []string
	var redirs []*ast.Redirection
	for _, el := range cmd.Elements {
		switch v := el.(type) {
		case *ast.Word:
			if name == "" {
				name = v.Value
			} else {
				args = append(args, v.Value)
			}
		case *ast.Redirection:
			redirs = append(redirs, v)
		}
	}

	// リダイレクト適用
	for _, r := range redirs {
		logmgr.Shell().Debug("ShellExecRedirection", "operator", r.Operator, "target", r.Target)
		if r.Target != nil {
			if err := ioctx.Redirect(r.Operator, r.Target.Value); err != nil {
				return 1, err
			}
		}
	}

	ex := s.getExecuter(cmd)
	argv := append([]string{name}, args...)
	return ex.Exec(argv, *ioctx)
}

// 指定されたコマンドが内部コマンドかどうかを判定
func (s *Shell) IsInternalCmd(cmd string) bool {
	return builtin.Cmds.Get(cmd) != nil
}

// 指定されたコマンドに基づいて適切なExecuterを取得
func (sn *Shell) getExecuter(cmd *ast.SimpleCommand) executer.Executer {
	name := ""
	for _, el := range cmd.Elements {
		if w, ok := el.(*ast.Word); ok {
			name = w.Value
			break
		}
	}
	var ex executer.Executer
	if sn.IsInternalCmd(name) {
		ex = &executer.InternalExecuter{}
	} else {
		ex = &executer.ExternalExecuter{}
	}
	return ex
}

func (s *Shell) defaultIOContext() executer.IOContext {
	return executer.IOContext{
		Stdin:  s.TTY.Input(),
		Stdout: s.TTY.Output(),
		Stderr: s.TTY.Output(),
	}
}

// SimpleCommand専用: ExecNodeから呼ばれる
func (s *Shell) Exec(cmd *ast.SimpleCommand) (int, error) {
	if cmd == nil || len(cmd.Elements) == 0 {
		return 0, nil
	}

	ioctx := s.defaultIOContext()

	// コマンド名・引数・リダイレクト抽出
	var name string
	var args []string
	var redirs []*ast.Redirection
	for _, el := range cmd.Elements {
		switch v := el.(type) {
		case *ast.Word:
			if name == "" {
				name = v.Value
			} else {
				args = append(args, v.Value)
			}
		case *ast.Redirection:
			redirs = append(redirs, v)
		}
	}

	if len(redirs) != 0 {
		logmgr.Shell().Debug("ShellParsedCommand", "name", name, "args", args)
		for _, r := range redirs {
			logmgr.Shell().Debug("ShellExecRedirection", "operator", r.Operator, "target", r.Target)
			if r.Target != nil {
				if err := ioctx.Redirect(r.Operator, r.Target.Value); err != nil {
					return 1, err
				}
			}
		}
	}

	defer ioctx.Close()

	ex := s.getExecuter(cmd)
	argv := append([]string{name}, args...)
	ec, err := ex.Exec(argv, ioctx)
	s.prompt.SetExitCode(ec)
	return ec, err
}
