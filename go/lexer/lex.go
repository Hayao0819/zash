package lexer

// Lexer は字句解析を行う構造体。状態と解析対象の文字列を保持する。
type Lexer struct {
	state lexerState
	input string
}

// lexerState は字句解析器の現在の状態を表す列挙型。
type lexerState int

const (
	_               lexerState = iota
	lexText                    // 空白を読み飛ばす状態
	lexEscapeChar              // バックスラッシュの直前まで読み進める状態
	lexQuotedString            // ダブルクォートの直前まで読み進める状態
	lexString                  // 特殊文字（空白、\、"）まで読み進める状態
)

// NewLexer は新しい Lexer を初期化して返す。
func NewLexer(input string) *Lexer {
	return &Lexer{
		state: lexText,
		input: input,
	}
}

// NextToken は現在の状態に応じて次のトークンを切り出して返す。
func (l *Lexer) NextToken() (string, error) {
	switch l.state {
	case lexText:
		// 空白文字の間をトークンとして切り出す
		return l.lexWhile(func(b byte) bool { return b == ' ' }, lexEscapeChar)
	case lexEscapeChar:
		// バックスラッシュが出るまでの文字列を切り出す
		return l.lexWhile(func(b byte) bool { return b != '\\' }, lexText)
	case lexQuotedString:
		// ダブルクォートが出るまでの文字列を切り出す
		return l.lexWhile(func(b byte) bool { return b != '"' }, lexText)
	case lexString:
		// 空白、バックスラッシュ、ダブルクォートのいずれかが出るまでの文字列を切り出す
		return l.lexWhile(func(b byte) bool {
			return b != '\\' && b != '"' && b != ' '
		}, lexText)
	default:
		// 未知の状態（安全のため空文字を返す）
		return "", nil
	}
}

// lexWhile は matchFn が true を返す限り文字を読み進め、トークンを切り出す共通処理。
// nextState に指定された状態に遷移する。
func (l *Lexer) lexWhile(matchFn func(byte) bool, nextState lexerState) (string, error) {
	i := 0
	for i < len(l.input) && matchFn(l.input[i]) {
		i++
	}

	// トークンを切り出し、残りの入力に更新
	token := l.input[:i]
	l.input = l.input[i:]

	// 入力が残っていれば次の状態へ、なければ初期状態に戻す
	if len(l.input) > 0 {
		l.state = nextState
	} else {
		l.state = lexText
	}
	return token, nil
}

// 全てのトークンを読んで返す
func (l *Lexer) ReadAll() ([]string, error) {
	var tokens []string
	for {
		token, err := l.NextToken()
		if err != nil {
			return nil, err
		}
		if token == "" {
			break
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}
