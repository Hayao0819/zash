package lexer

import (
	"testing"
)

// テストケース構造体
type testCase struct {
	input    string   // 入力文字列
	expected []string // 期待されるトークンの順序
}

func TestLexer(t *testing.T) {
	tests := []testCase{
		{
			input:    "   abc \\x \"quoted text\" plain_text",
			expected: []string{"   ", "abc ", "\\x", " ", "quoted text", " ", "plain_text"},
		},
		{
			input:    "one\\two\"three four\" five",
			expected: []string{"one", "\\", "two", "\"", "three four", "\"", " ", "five"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "   ",
			expected: []string{"   "},
		},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		var tokens []string

		// 最後までトークンを取得
		for {
			token, err := lexer.NextToken()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if token == "" {
				break
			}
			tokens = append(tokens, token)
		}

		if len(tokens) != len(tt.expected) {
			t.Errorf("input: %q\nexpected %d tokens, got %d\nexpected: %#v\ngot: %#v",
				tt.input, len(tt.expected), len(tokens), tt.expected, tokens)
			continue
		}

		for i := range tokens {
			if tokens[i] != tt.expected[i] {
				t.Errorf("input: %q\ntoken %d: expected %q, got %q",
					tt.input, i, tt.expected[i], tokens[i])
			}
		}
	}
}
