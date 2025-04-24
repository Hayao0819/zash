package lexer

import (
	"log/slog"
)


// GetNextState は次に遷移すべき状態を返す。
func (l *Lexer) getNextState() state {
	// 現在の状態に応じて次の状態を決定する
	for _, s := range states {
		if s.determineFunc != nil && s.determineFunc(l) {
			l.state = s
			return s
		}
	}
	return lexInitState
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

	slog.Debug("LexerNextToken", "state", l.state.name, "remaining", l.left())

	// 状態に応じて適切な処理を実行
	switch l.state.name {
	case lexInitState.name:
		l.state = l.getNextState()
		return l.NextToken()
	default:
		return l.state.lexFunc(l)
	}
}

// lexWhile は matchFn が true を返す限り文字を読み進め、トークンを切り出す共通処理。
func (l *Lexer) lexWhile(t TokenType, matchFn func(byte) bool) (*Token, error) {
	remaining := l.left()
	i := 0
	for i < len(remaining) && matchFn(remaining[i]) {
		i++
	}

	// トークンを切り出し、残りの入力に更新
	token := remaining[:i]
	l.processed += i

	// 状態を更新
	l.state = lexInitState
	// return token, nil
	return &Token{
		Type: t,
		Text: token,
	}, nil
}
