package lexer

import "log/slog"

// lexerState は字句解析器の現在の状態を表す列挙型。
type state struct {
	name          string                       // 状態名
	determineFunc func(*Lexer) bool            // 状態遷移を決定する関数
	lexFunc       func(*Lexer) (*Token, error) // 状態に応じた処理関数
}

var states = []state{
	lexInitState,
	lexWhitespaceState,
	lexEscapeCharState,
	lexRedirectionState,
	lexAndState,
	lexPipeState,
	lexQuotedStringState,
	lexCommentState,
	lexNumberState, // NumberStateはStringの直前に追加
	lexStringState, // StringStateは最後に追加
}

// 初期状態
var lexInitState = state{
	name: "lexInit",
}

var lexWhitespaceState = state{
	name: "lexWhitespace",
	determineFunc: func(l *Lexer) bool {
		return l.left()[0] == ' '
	},
	lexFunc: func(l *Lexer) (*Token, error) {
		return l.lexWhile(TokenWhitespace, func(b byte) bool {
			return b == ' '
		})
	},
}

var lexEscapeCharState = state{
	name: "lexEscapeChar",
	determineFunc: func(l *Lexer) bool {
		return l.left()[0] == '\\'
	},
	lexFunc: func(l *Lexer) (*Token, error) {
		remaining := l.left()
		// バックスラッシュとその次の1文字を取得
		if len(remaining) > 1 {
			tok := remaining[:2]
			slog.Debug("lexEscapeChar", "tok", tok)
			l.processed += 2
			l.state = lexInitState
			// return tok, nil
			return &Token{
				Type: TokenEscapeChar,
				Text: tok,
			}, nil
		}

		// 1文字しかない場合（エスケープの末尾）
		tok := remaining[:1]
		l.processed++
		l.state = lexInitState

		// return tok, nil
		return &Token{
			Type: TokenEscapeChar,
			Text: tok,
		}, nil
	},
}

var lexQuotedStringState = state{
	name: "lexQuotedString",
	determineFunc: func(l *Lexer) bool {
		return l.left()[0] == '"'
	},
	lexFunc: func(l *Lexer) (*Token, error) {
		remaining := l.left()

		if remaining[0] == '"' {
			// 先頭のクォートだけ返す
			l.processed++
			// l.state = lexQuotedStringState
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
		return &Token{
			Type: TokenQuotedString,
			Text: tok,
		}, nil
	},
}

var lexRedirectionState = state{
	name: "lexRedirection",
	determineFunc: func(l *Lexer) bool {
		return l.left()[0] == '>' || l.left()[0] == '<'
	},
	lexFunc: func(l *Lexer) (*Token, error) {
		return l.lexWhile(TokenRedirection, func(b byte) bool {
			return b == '>' || b == '<'
		})
	},
}

var lexAndState = state{
	name: "lexAnd",
	determineFunc: func(l *Lexer) bool {
		return l.left()[0] == '&'
	},
	lexFunc: func(l *Lexer) (*Token, error) {
		return l.lexWhile(TokenAnd, func(b byte) bool {
			return b == '&'
		})
	},
}

var lexCommentState = state{
	name: "lexComment",
	determineFunc: func(l *Lexer) bool {
		return l.left()[0] == '#'
	},
	lexFunc: func(l *Lexer) (*Token, error) {
		remaining := l.left()
		i := 0
		for i < len(remaining) && remaining[i] != '\n' {
			i++
		}

		tok := remaining[:i]
		l.processed += i

		// 状態を更新
		l.state = lexInitState
		return &Token{
			Type: TokenComment,
			Text: tok,
		}, nil
	},
}

var lexPipeState = state{
	name: "lexPipe",
	determineFunc: func(l *Lexer) bool {
		return l.left()[0] == '|'
	},
	lexFunc: func(l *Lexer) (*Token, error) {
		return l.lexWhile(TokenPipe, func(b byte) bool {
			return b == '|'
		})
	},
}

var lexNumberState = state{
	name: "lexNumber",
	determineFunc: func(l *Lexer) bool {
		return l.left()[0] >= '0' && l.left()[0] <= '9'
	},
	lexFunc: func(l *Lexer) (*Token, error) {
		return l.lexWhile(TokenNumber, func(b byte) bool {
			return b >= '0' && b <= '9'
		})
	},
}

var lexStringState = state{
	name: "lexString",
	determineFunc: func(l *Lexer) bool {
		// return l.left()[0] != ' ' && l.left()[0] != '\\' && l.left()[0] != '"'
		return true
	},
	lexFunc: func(l *Lexer) (*Token, error) {
		return l.lexWhile(TokenString, func(b byte) bool {
			return b != ' ' && b != '\\' && b != '"'
		})
	},
}
