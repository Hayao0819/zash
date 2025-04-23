package shell

import (
	"log/slog"

	"github.com/Hayao0819/zash/go/lexer"
	"github.com/Hayao0819/zash/go/parser"
	"github.com/samber/lo"
)

func (s *Shell) Run(line string) error {

	tokens, err := lexer.NewLexer(line).ReadAll()
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return nil
	}
	{
		tj := lo.Map(tokens, func(t lexer.Token, i int) string {
			return string(t.JSON())
		})
		slog.Debug("ShellGotTokens", "tokens", tj)
	}

	st, err := parser.NewParser(tokens).Parse()
	if err != nil {
		return err
	}
	if st == nil {
		return nil
	}
	if st.Cmd == nil {
		return nil
	}

	if _, err := s.Exec(st.Cmd); err != nil {
		return err
	}
	return nil

}
