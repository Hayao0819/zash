package shell

import (
	"github.com/Hayao0819/zash/go/lexer"
	"github.com/Hayao0819/zash/go/parser"
)

func (s *Shell) StartInteractive() {
	if s.started {
		return
	}
	s.started = true

	for {

		tokens, err := lexer.NewLexer(s.WaitInputWithPrompt()).ReadAll()
		if err != nil {
			s.Println(err.Error())
			continue
		}
		if len(tokens) == 0 {
			continue
		}

		cmd, err := parser.NewParser(tokens).Parse()
		if err != nil {
			s.Println(err.Error())
			continue
		}

		argv := []string{cmd.Name}
		if cmd.CommandSuffix != nil {
			argv = append(argv, cmd.CommandSuffix.Args...)
		}

		if err := s.Exec(argv); err != nil {
			// fmt.Println("Err: ", err)
			s.Println(err.Error())
		}
	}
}
