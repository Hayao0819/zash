package lexer

import "encoding/json"

// isWhitespace は空白かどうかを判定する関数。
func isWhitespace(b byte) bool {
	return b == ' '
}

// isNormalStringChar は通常の文字列に含まれるか（空白、\、"以外）を判定する関数。
func isNormalStringChar(b byte) bool {
	return b != ' ' && b != '\\' && b != '"'
}

func PrintJSON(ts []Token) {
	j, err := json.Marshal(ts)
	if err != nil {
		panic(err)
	}
	println(string(j))
}

func isRedirection(b byte) bool {
	return b == '>' || b == '<'
}

func isPipe(b byte) bool {
	return b == '|'
}
