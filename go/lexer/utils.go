package lexer

// isWhitespace は空白かどうかを判定する関数。
func isWhitespace(b byte) bool {
	return b == ' '
}

// isNormalStringChar は通常の文字列に含まれるか（空白、\、"以外）を判定する関数。
func isNormalStringChar(b byte) bool {
	return b != ' ' && b != '\\' && b != '"'
}
