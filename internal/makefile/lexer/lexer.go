package lexer

import (
	"fmt"
	"regexp"
)

type Token uint8

const (
	NewLine Token = iota
	String
	Number
	Comment
	Whitespace
	Operator
	Parenthesis
	Identifier
	Indent
	Unknown
)

var Rules = map[Token]string{
	NewLine:     `^\n$`,
	String:      `^".*"$`,
	Number:      `^[0-9]+(.[0-9]+)?$`,
	Comment:     `^#[^\n]*$`,
	Indent:      `^\t$`,
	Whitespace:  `^[\s]+$`,
	Operator:    `^[=:]+$`,
	Parenthesis: `^[(){}]$`,
	Identifier:  `^[A-Za-z_-0-9]+$`,
}

type Match struct {
	Type  Token
	Value string
}

func Lex(feed chan rune) (chan *Match, chan error) {
	tokens := make(chan *Match)
	errors := make(chan error)

	go func() {
		defer close(tokens)
		defer close(errors)

		open := true
		last_test := string(<-feed)
		test := last_test
		last_token := Unknown
		for open {
			for token, pattern := range Rules {
				matched, _ := regexp.MatchString(pattern, test)
				if !matched {
					if last_token == Unknown {
						errors <- fmt.Errorf("syntax error near '%s'", test)
						return
					} else {
						tokens <- &Match{Type: token, Value: last_test}
						test = ""
					}
				}
			}
			last_test = test
			char, ok := <-feed
			open = ok
			test += string(char)
		}
	}()

	return tokens, errors
}
