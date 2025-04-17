package lexer

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
		default:
			return lexString
		}
	}
	return lexText
}

// NextToken は現在の状態に応じて次のトークンを切り出して返す。
func (l *Lexer) NextToken() (string, error) {
	// すべての文字が処理されていた場合
	if len(l.left()) == 0 {
		return "", nil
	}

	// 状態に応じて適切な処理を実行
	switch l.state {
	case lexText:
		l.state = l.getNextState()
		return l.NextToken()

	case lexWhitespace:
		return l.lexWhile(isWhitespace)

	case lexEscapeChar:
		return l.lexEscapeChar()

	case lexQuotedString:
		return l.lexQuotedString()

	case lexString:
		return l.lexWhile(isNormalStringChar)
	}

	return "", nil
}

// lexWhile は matchFn が true を返す限り文字を読み進め、トークンを切り出す共通処理。
func (l *Lexer) lexWhile(matchFn func(byte) bool) (string, error) {
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
	return token, nil
}

// lexEscapeChar はバックスラッシュ文字とその次の文字を読み進める。
func (l *Lexer) lexEscapeChar() (string, error) {
	remaining := l.left()
	// バックスラッシュとその次の1文字を取得
	if len(remaining) > 1 {
		tok := remaining[:2]
		// slog.Info("lexEscapeChar", "tok", tok)
		l.processed += 2
		l.state = lexText
		return tok, nil
	}

	// 1文字しかない場合（エスケープの末尾）
	tok := remaining[:1]
	l.processed++
	l.state = lexText
	return tok, nil
}

// lexQuotedString はダブルクォート内の文字列を切り出す
func (l *Lexer) lexQuotedString() (string, error) {
	remaining := l.left()

	if remaining[0] == '"' {
		// 先頭のクォートだけ返す
		l.processed++
		l.state = lexQuotedString
		return `"`, nil
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
	return tok, nil
}
