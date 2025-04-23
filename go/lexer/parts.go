package lexer

import "log/slog"

func (l *Lexer) lexWhitespace() (*Token, error) {
	return l.lexWhile(TokenWhitespace, func(b byte) bool {
		return b == ' '
	})
}

func (l *Lexer) lexRedirection() (*Token, error) {
	return l.lexWhile(TokenRedirection, func(b byte) bool {
		return b == '>' || b == '<'
	})
}

func (l *Lexer) lexString() (*Token, error) {
	return l.lexWhile(TokenString, func(b byte) bool {
		return b != ' ' && b != '\\' && b != '"'
	})
}

func (l *Lexer) lexPipe() (*Token, error) {
	return l.lexWhile(TokenPipe, func(b byte) bool {
		return b == '|'
	})
}

func (l *Lexer) lexComemnt() (*Token, error) {
	remaining := l.left()
	i := 0
	for i < len(remaining) && remaining[i] != '\n' {
		i++
	}

	tok := remaining[:i]
	l.processed += i

	// 状態を更新
	l.state = lexText
	// return tok, nil
	return &Token{
		Type: TokenComment,
		Text: tok,
	}, nil
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
	l.state = lexText
	// return token, nil
	return &Token{
		Type: t,
		Text: token,
	}, nil
}

// lexEscapeChar はバックスラッシュ文字とその次の文字を読み進める。
func (l *Lexer) lexEscapeChar() (*Token, error) {
	remaining := l.left()
	// バックスラッシュとその次の1文字を取得
	if len(remaining) > 1 {
		tok := remaining[:2]
		slog.Debug("lexEscapeChar", "tok", tok)
		l.processed += 2
		l.state = lexText
		// return tok, nil
		return &Token{
			Type: TokenEscapeChar,
			Text: tok,
		}, nil
	}

	// 1文字しかない場合（エスケープの末尾）
	tok := remaining[:1]
	l.processed++
	l.state = lexText

	// return tok, nil
	return &Token{
		Type: TokenEscapeChar,
		Text: tok,
	}, nil
}

// lexQuotedString はダブルクォート内の文字列を切り出す
func (l *Lexer) lexQuotedString() (*Token, error) {
	remaining := l.left()

	if remaining[0] == '"' {
		// 先頭のクォートだけ返す
		l.processed++
		l.state = lexQuotedString
		// return `"`, nil
		return &Token{
			Type: TokenQuoteChar,
			Text: `"`,
		}, nil
	}

	// クォーテーションが閉じるまで読み取る
	i := 0
	for i < len(remaining) && remaining[i] != '"' {
		i++
	}

	tok := remaining[:i]
	l.processed += i

	// 次は閉じクォートを処理する
	l.state = lexQuotedString
	// return tok, nil
	return &Token{
		Type: TokenQuotedString,
		Text: tok,
	}, nil
}
