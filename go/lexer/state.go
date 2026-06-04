package lexer

import (
	"errors"
	"log/slog"
)

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
	lexSingleQuotedStringState,
	lexCommentState,
	lexNumberState, // NumberStateはStringの直前に追加
	lexStringState, // StringStateは最後に追加
}

// シングルクォート文字列状態
var lexSingleQuotedStringState = state{
	name: "lexSingleQuotedString",
	determineFunc: func(l *Lexer) bool {
		return l.left()[0] == '\''
	},
	lexFunc: func(l *Lexer) (*Token, error) {
		remaining := l.left()
		// 開きクォートをスキップ
		i := 1
		// クォーテーションが閉じるまで読み取る
		for i < len(remaining) && remaining[i] != '\'' {
			i++
		}
		if i >= len(remaining) {
			return nil, errors.New("syntax error: unmatched single quote")
		}
		// 閉じクォートもスキップ
		i++
		// クォート内の内容（クォート自体を除く）を返す
		tok := remaining[1 : i-1]
		l.processed += i
		l.state = lexInitState
		return &Token{
			Type: TokenSingleQuotedString,
			Text: tok,
		}, nil
	},
}

// 初期状態
var lexInitState = state{
	name: "lexInit",
}

var lexWhitespaceState = state{
	name: "lexWhitespace",
	determineFunc: func(l *Lexer) bool {
		return l.left()[0] == ' ' || l.left()[0] == '\t'
	},
	lexFunc: func(l *Lexer) (*Token, error) {
		return l.lexWhile(TokenWhitespace, func(b byte) bool {
			return b == ' ' || b == '\t'
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
			return &Token{
				Type: TokenEscapeChar,
				Text: tok,
			}, nil
		}

		// 1文字しかない場合（エスケープの末尾）
		tok := remaining[:1]
		l.processed++
		l.state = lexInitState

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
		// 開きクォートをスキップ
		i := 1
		// クォーテーションが閉じるまで読み取る（エスケープ対応）
		for i < len(remaining) {
			if remaining[i] == '\\' && i+1 < len(remaining) {
				i += 2 // エスケープ文字をスキップ
				continue
			}
			if remaining[i] == '"' {
				break
			}
			i++
		}
		if i >= len(remaining) {
			return nil, errors.New("syntax error: unmatched quote")
		}
		// 閉じクォートもスキップ
		i++
		// クォート内の内容（クォート自体を除く）を返す
		tok := remaining[1 : i-1]
		l.processed += i
		l.state = lexInitState
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
		remaining := l.left()
		if len(remaining) >= 2 && remaining[1] == '&' {
			l.processed += 2
			l.state = lexInitState
			return &Token{Type: TokenAnd, Text: "&&"}, nil
		}
		l.processed++
		l.state = lexInitState
		return &Token{Type: TokenBackground, Text: "&"}, nil
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
		// コメント行（改行まで）を読む
		for i < len(remaining) && remaining[i] != '\n' {
			i++
		}

		tok := remaining[:i]
		l.processed += i

		// 状態を初期状態にリセット
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
		remaining := l.left()
		if len(remaining) >= 2 && remaining[1] == '|' {
			l.processed += 2
			l.state = lexInitState
			return &Token{Type: TokenOr, Text: "||"}, nil
		}
		l.processed++
		l.state = lexInitState
		return &Token{Type: TokenPipe, Text: "|"}, nil
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
		c := l.left()[0]
		return len(l.left()) > 0 && c != ' ' && c != '\t'
	},
	lexFunc: func(l *Lexer) (*Token, error) {
		remaining := l.left()
		if len(remaining) == 0 {
			return nil, nil
		}
		// 区切り記号を優先的にトークン化
		switch remaining[0] {
		case ';':
			l.processed++
			l.state = lexInitState
			return &Token{Type: TokenSemicolon, Text: ";"}, nil
		case '\n':
			l.processed++
			l.state = lexInitState
			return &Token{Type: TokenNewline, Text: "\n"}, nil
		case '{':
			l.processed++
			l.state = lexInitState
			return &Token{Type: TokenLBrace, Text: "{"}, nil
		case '}':
			l.processed++
			l.state = lexInitState
			return &Token{Type: TokenRBrace, Text: "}"}, nil
		case '(':
			l.processed++
			l.state = lexInitState
			return &Token{Type: TokenLParen, Text: "("}, nil
		case ')':
			l.processed++
			l.state = lexInitState
			return &Token{Type: TokenRParen, Text: ")"}, nil
		case '!':
			l.processed++
			l.state = lexInitState
			return &Token{Type: TokenBang, Text: "!"}, nil
		case '=':
			l.processed++
			l.state = lexInitState
			return &Token{Type: TokenAssign, Text: "="}, nil
		}
		// 1単語を切り出す（区切り記号・空白・バックスラッシュ・クォートで区切る）
		i := 0
		for i < len(remaining) {
			c := remaining[i]
			if c == ' ' || c == '\t' || c == '\\' || c == '"' || c == '\'' ||
				c == ';' || c == '\n' || c == '=' || c == '{' || c == '}' ||
				c == '(' || c == ')' || c == '!' || c == '|' || c == '&' ||
				c == '>' || c == '<' || c == '#' {
				break
			}
			i++
		}
		if i == 0 {
			// 1文字も進まなかった場合は1文字だけ消費して返す（無限ループ防止）
			l.processed++
			l.state = lexInitState
			return &Token{Type: TokenString, Text: string(remaining[0])}, nil
		}
		word := remaining[:i]
		l.processed += i
		l.state = lexInitState
		// キーワード判定
		switch string(word) {
		case "if":
			return &Token{Type: TokenIf, Text: string(word)}, nil
		case "then":
			return &Token{Type: TokenThen, Text: string(word)}, nil
		case "else":
			return &Token{Type: TokenElse, Text: string(word)}, nil
		case "elif":
			return &Token{Type: TokenElif, Text: string(word)}, nil
		case "fi":
			return &Token{Type: TokenFi, Text: string(word)}, nil
		case "for":
			return &Token{Type: TokenFor, Text: string(word)}, nil
		case "while":
			return &Token{Type: TokenWhile, Text: string(word)}, nil
		case "until":
			return &Token{Type: TokenUntil, Text: string(word)}, nil
		case "do":
			return &Token{Type: TokenDo, Text: string(word)}, nil
		case "done":
			return &Token{Type: TokenDone, Text: string(word)}, nil
		case "case":
			return &Token{Type: TokenCase, Text: string(word)}, nil
		case "esac":
			return &Token{Type: TokenEsac, Text: string(word)}, nil
		case "select":
			return &Token{Type: TokenSelect, Text: string(word)}, nil
		case "in":
			return &Token{Type: TokenIn, Text: string(word)}, nil
		case "function":
			return &Token{Type: TokenFunction, Text: string(word)}, nil
		default:
			return &Token{Type: TokenString, Text: string(word)}, nil
		}
	},
}
