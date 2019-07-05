package main

import (
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			input: "GET hello",
			expected: []Token{
				Token{Type: GET, Literal: []byte("GET")},
				Token{Type: OPERAND, Literal: []byte("hello")},
			},
		},
		{
			input: "delete HELLO",
			expected: []Token{
				Token{Type: DELETE, Literal: []byte("delete")},
				Token{Type: OPERAND, Literal: []byte("HELLO")},
			},
		},
		{
			input: "put hello world",
			expected: []Token{
				Token{Type: PUT, Literal: []byte("put")},
				Token{Type: OPERAND, Literal: []byte("hello")},
				Token{Type: OPERAND, Literal: []byte("world")},
			},
		},
		{
			input: "     put  hello      world    ",
			expected: []Token{
				Token{Type: PUT, Literal: []byte("put")},
				Token{Type: OPERAND, Literal: []byte("hello")},
				Token{Type: OPERAND, Literal: []byte("world")},
			},
		},
		{
			input: "foo hello",
			expected: []Token{
				Token{Type: UNKNOWN, Literal: []byte("foo")},
				Token{Type: OPERAND, Literal: []byte("hello")},
			},
		},
	}

	for _, tt := range tests {
		tokens := tokenizer(tt.input)
		if !reflect.DeepEqual(tokens, tt.expected) {
			t.Fatalf("tokens are not match (got=%#v, expected=%#v)", tokens, tt.expected)
		}
	}
}
