package lexer

import "encoding/json"

type TokenType int

const (
	_                 TokenType = iota
	TokenWhitespace             // " " character
	TokenEscapeChar             // \ character
	TokenQuoteChar              // " character
	TokenQuotedString           // "string"
	TokenString                 // string
	TokenRedirection            // > or <
	TokenUnknown                // unknown token
	TokenEOT                    // End of Text
)

func (t TokenType) String() string {
	switch t {
	case TokenWhitespace:
		return "Whitespace"
	case TokenEscapeChar:
		return "EscapeChar"
	case TokenQuoteChar:
		return "QuoteChar"
	case TokenQuotedString:
		return "QuotedString"
	case TokenString:
		return "String"
	case TokenRedirection:
		return "Redirection"
	case TokenUnknown:
		return "Unknown"
	case TokenEOT:
		return "EOT"
	default:
		return "Unknown"
	}
}

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

func (t *Token) JSON() []byte {
	tj := struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
	}{
		Type: t.Type.String(),
		Text: t.Text,
	}
	j, _ := json.Marshal(tj)
	return j
}
