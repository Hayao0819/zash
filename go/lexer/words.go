package lexer

// func LineToWords(line string) (word []string, err error) {
// 	l := NewLexer(strings.TrimSpace(line))
// 	tokens, err := l.ReadAll()
// 	if err != nil {
// 		return nil, err
// 	}

// 	words := []string{}
// 	for _, token := range tokens {
// 		if token.Type == TokenQuoteChar{
// 			continue
// 		}
// 		if token.Type == TokenWhitespace {
// 			// words = append(words, " ")
// 			continue
// 		}
// 		if token.Type == TokenEscapeChar {
// 			words = append(words, token.Text[1:])
// 			continue
// 		}
// 		words = append(words, token.Text)
// 	}
// 	return words, nil
// }
