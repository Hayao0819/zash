package lexer

// 全てのトークンを読んで返す
func (l *Lexer) ReadAll() ([]string, error) {
	var tokens []string
	for {
		token, err := l.NextToken()
		if err != nil {
			return nil, err
		}
		if token == "" {
			break
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}
