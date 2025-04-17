package lexer

// Implement the io.Reader interface
func (l *Lexer) Read(p []byte) (n int, err error) {
	if l.state == lexText {
		return 0, nil
	}
	token, err := l.NextToken()
	if err != nil {
		return 0, err
	}
	n = copy(p, token)
	l.input = l.input[n:]
	return n, nil
}
