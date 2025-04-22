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

	cmd, err := parser.NewParser(tokens).Parse()
	if err != nil {
		return err
	}

	_argv := []string{cmd.Name}
	if cmd.CommandSuffix != nil {
		_argv = append(_argv, cmd.CommandSuffix.Args...)
	}

	if err := s.Exec(_argv); err != nil {
		return err
	}
	return nil

}
