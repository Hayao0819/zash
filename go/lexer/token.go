package lexer

type TokenType int

const (
	_ TokenType = iota
	TokenWhitespace
	TokenEscapeChar
	TokenQuoteChar
	TokenQuotedString
	TokenString
	TokenUnknown
	TokenEOT // End of Text
)

type Token struct {
	Type TokenType
	Text string
}

func (t Token) String() string {
	if t.Type == TokenWhitespace {
		return " "
	}
	if t.Type == TokenEOT {
		return ""
	}
	if t.Type == TokenEscapeChar {
		return t.Text[1:]
	}
	return t.Text
}
