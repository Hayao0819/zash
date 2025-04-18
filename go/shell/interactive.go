package shell

import (
	"github.com/Hayao0819/zash/go/lexer"
)

func (s *Shell) StartInteractive() {
	if s.started {
		return
	}
	s.started = true

	for {
		input := s.WaitInputWithPrompt()

		words, err := lexer.LineToWords(input)
		if err != nil {
			// fmt.Println("Err: ", err)
			s.Println(err.Error())
			continue
		}
		if len(words) == 0 {
			continue
		}

		// for _, word := range words {
		// 	fmt.Println(word)
		// }

		if err := s.Exec(words); err != nil {
			// fmt.Println("Err: ", err)
			s.Println(err.Error())
		}
	}
}
