package lexer

// 全てのトークンを読んで返す
func (l *Lexer) ReadAll() ([]Token, error) {
	var tokens []Token
	for {
		token, err := l.NextToken()
		if err != nil {
			return nil, err
		}
		if token == nil {
			break
		}
		tokens = append(tokens, *token)
	}
	return tokens, nil
}
