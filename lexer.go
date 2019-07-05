package main

import (
	"fmt"
	"strings"
)

const (
	UNKNOWN = iota
	GET
	PUT
	DELETE
	OPERAND
)

type Token struct {
	Type    int
	Literal []byte
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %d, Literal: %s}", t.Type, string(t.Literal))
}

func tokenizer(input string) []Token {
	ss := strings.Fields(input)

	if len(ss) == 0 {
		return nil
	}

	tokens := make([]Token, len(ss))

	for i, v := range ss {
		if i == 0 {
			var t int
			switch strings.ToUpper(v) {
			case "GET":
				t = GET
			case "PUT":
				t = PUT
			case "DELETE":
				t = DELETE
			}
			tokens[i] = Token{Type: t, Literal: []byte(v)}
		} else {
			tokens[i] = Token{Type: OPERAND, Literal: []byte(v)}
		}
	}

	return tokens
}
