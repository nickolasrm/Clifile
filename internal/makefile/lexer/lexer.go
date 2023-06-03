package lexer

import (
	"fmt"
	"regexp"
)

type Token uint8

const (
	NewLine Token = iota
	Whitespace
	Parenthesis
	Operator
	Action
	Indent
	Comment
	Identifier
	Value
	UnquotedValue
	Unknown
)

var Rules = map[Token]string{
	NewLine:       `^\n$`,
	Value:         `^"[^"]*"?$`,
	UnquotedValue: `^[^\n]+$`,
	Comment:       `^#[^\n]*$`,
	Action:        `^\t[^\n]+$`,
	Indent:        `^\t$`,
	Whitespace:    `^[\s]+$`,
	Operator:      `^[=:$]+$`,
	Parenthesis:   `^[(){}]$`,
	Identifier:    `^[A-Za-z_-0-9]+$`,
	Unknown:       `.*`,
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
		last := Match{Type: Unknown, Value: ""}
		test := Match{Type: Unknown, Value: string(<-feed)}
		for open {
			for i := Token(0); i <= Unknown; i++ {
				matched, _ := regexp.MatchString(Rules[i], test.Value)
				test.Type = i
				if matched {
					break
				}
			}
			if test.Type == Unknown {
				if last.Type != Unknown {
					tokens <- &last
					last = Match{Type: Unknown, Value: ""}
					test = Match{Type: Unknown, Value: test.Value[len(test.Value)-1:]}
					continue
				} else {
					errors <- fmt.Errorf("syntax error near '%s'", test.Value)
					return
				}
			} else {
				var char rune
				char, open = <-feed
				if !open {
					tokens <- &test
					return
				}
				last.Type = test.Type
				last.Value = test.Value
				test.Value += string(char)
			}
		}
	}()

	return tokens, errors
}
