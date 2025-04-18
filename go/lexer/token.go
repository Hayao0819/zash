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
)

type Token struct {
	Type TokenType
	Text string
}
