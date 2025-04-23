package lexer

import (
	"log/slog"
)

// Lexer は字句解析を行う構造体。状態と解析対象の文字列を保持する。
type Lexer struct {
	state     lexerState
	input     string
	processed int
}

// left は未処理の残り文字列を返す。
func (l *Lexer) left() string {
	return l.input[l.processed:]
}

// GetNextState は次に遷移すべき状態を返す。
func (l *Lexer) getNextState() lexerState {
	remaining := l.left()
	// 未処理の文字列の先頭を確認
	if len(remaining) > 0 {
		switch remaining[0] {
		case ' ':
			return lexWhitespace
		case '\\':
			return lexEscapeChar
		case '"':
			return lexQuotedString
		case '>', '<':
			return lexRedirection
		case '#':
			return lexComment
		case '|':
			return lexPipe
		default:
			return lexString
		}
	}
	return lexText
}

// NextToken は現在の状態に応じて次のトークンを切り出して返す。
func (l *Lexer) NextToken() (*Token, error) {
	// すべての文字が処理されていた場合
	if len(l.left()) == 0 {
		return &Token{
			Type: TokenEOT,
			Text: "",
		}, nil
	}

	slog.Debug("LexerNextToken", "state", l.state.Text(), "remaining", l.left())

	// 状態に応じて適切な処理を実行
	switch l.state {
	case lexText:
		l.state = l.getNextState()
		return l.NextToken()

	case lexWhitespace:
		return l.lexWhitespace()
	case lexEscapeChar:
		return l.lexEscapeChar()
	case lexComment:
		return l.lexComemnt()
	case lexQuotedString:
		return l.lexQuotedString()
	case lexRedirection:
		return l.lexRedirection()
	case lexString:
		return l.lexString()
	case lexPipe:
		return l.lexPipe()
	}

	return nil, nil
}
