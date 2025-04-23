package lexer

// lexerState は字句解析器の現在の状態を表す列挙型。
type lexerState int

const (
	_               lexerState = iota
	lexText                    // 初期状態
	lexWhitespace              // 空白を連続して読み取る
	lexEscapeChar              // バックスラッシュとその次の1文字を読み取る
	lexQuotedString            // クォート内の文字列を読み取る
	lexString                  // 通常の文字列を読み取る
	lexRedirection             // リダイレクションを読み取る
	lexComment                 // コメントを読み取る
	lexPipe                    // パイプを読み取る
)

func (s lexerState) Text() string {
	return []string{
		"lexText",
		"lexWhitespace",
		"lexEscapeChar",
		"lexQuotedString",
		"lexString",
		"lexRedirection",
		"lexComment",
		"lexPipe",
	}[s]
}
