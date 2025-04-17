package lexer

// NewLexer は新しい Lexer を初期化して返す。
func NewLexer(input string) *Lexer {
	return &Lexer{
		state:     lexText,
		input:     input,
		processed: 0,
	}
}
