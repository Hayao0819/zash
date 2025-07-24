package shell

import (
	"github.com/Hayao0819/zash/go/lexer"
	"github.com/Hayao0819/zash/go/parser"
)

func (s *Shell) Run(line string) error {

	tokens, err := lexer.NewLexer(line).ReadAll()
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return nil
	}

	st, err := parser.NewParser(tokens).Parse()
	if err != nil {
		return err
	}
	if st == nil {
		return nil
	}
	// どのASTノードでも実行できるようExecNodeを呼ぶ
	if _, err := s.ExecNode(st); err != nil {
		return err
	}
	return nil

}
