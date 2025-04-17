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

// lexerState は字句解析器の現在の状態を表す列挙型。
type lexerState int

const (
	_               lexerState = iota
	lexText                    // 初期状態
	lexWhitespace              // 空白を連続して読み取る
	lexEscapeChar              // バックスラッシュとその次の1文字を読み取る
	lexQuotedString            // クォート内の文字列を読み取る
	lexString                  // 通常の文字列を読み取る
)

func lexerStateText(s lexerState) string {
	switch s {
	case lexText:
		return "lexText"
	case lexWhitespace:
		return "lexWhitespace"
	case lexEscapeChar:
		return "lexEscapeChar"
	case lexQuotedString:
		return "lexQuotedString"
	case lexString:
		return "lexString"
	default:
		return "unknown state"
	}
}

// NewLexer は新しい Lexer を初期化して返す。
func NewLexer(input string) *Lexer {
	return &Lexer{
		state:     lexText,
		input:     input,
		processed: 0,
	}
}

// left は未処理の残り文字列を返す。
func (l *Lexer) left() string {
	return l.input[l.processed:]
}

// GetNextState は次に遷移すべき状態を返す。
func (l *Lexer) GetNextState() lexerState {
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
	slog.Info("NextToken", "state", lexerStateText(l.state),"processed", l.processed, "left", l.left())

	// すべての文字が処理されていた場合
	if len(l.left()) == 0 {
		return "", nil
	}

	// 状態に応じて適切な処理を実行
	switch l.state {
	case lexText:
		l.state = l.GetNextState()
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
	// 最初のクォーテーションを切り出す
	if remaining[0] == '"' {
		l.processed++
		// クォート自体をトークンとして返す
		return `"`, nil
	}

	// クォーテーションが閉じるまで文字を読み進める
	i := 0
	for i < len(remaining) && remaining[i] != '"' {
		i++
	}

	tok := remaining[:i]
	l.processed += i
	
	l.state = lexText
	return tok, nil
}

// isWhitespace は空白かどうかを判定する関数。
func isWhitespace(b byte) bool {
	return b == ' '
}

// isNormalStringChar は通常の文字列に含まれるか（空白、\、"以外）を判定する関数。
func isNormalStringChar(b byte) bool {
	return b != ' ' && b != '\\' && b != '"'
}
