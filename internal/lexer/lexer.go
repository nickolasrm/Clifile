package lexer

import (
	"fmt"
	"regexp"
	"strings"
)

type Token uint8

const (
	Line Token = iota
	Indent
	Docstring
	Comment
	Function
	Variable
	Rule
	Action
	Unknown
)

// ruleRegex returns a regex that matches the pattern at the beginning of a string.
func ruleRegex(pattern string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(`^%s`, pattern))
}

var Rules = map[Token]*regexp.Regexp{
	Line:      ruleRegex(`\n+`),
	Indent:    ruleRegex(`\t+`),
	Docstring: ruleRegex(`##([^\n]*)`),
	Comment:   ruleRegex(`#([^\n]*)`),
	Function:  ruleRegex(`(\w+)=\${(?:\s+)?(\w+)(?:\s+)?([^}]+)?}`),
	Variable:  ruleRegex(`(\w+)=(?:"([^"]*)"|([^"\n]*))`),
	Rule:      ruleRegex(`(\w+):([\w ]*)`),
	Action:    ruleRegex(`([^\n]+)`),
}

// Match represents a token and the value that was matched.
// The type is the token that was matched.
// The value is a slice of strings that were matched by the regex capture groups.
// The first element of the slice is the entire match.
// The remaining elements are the capture groups.
type Match struct {
	Type  Token
	Value []string
}

// Lex reads a string and tokenizes it into a stream of tokens and errors.
// This function is the first step in the interpreter pipeline. It is used
// to identify the pieces of code that contain meaningful structure for
// parsing into a syntax tree.
// The tokens channel is closed when the input string is exhausted.
// The errors channel is closed when the input string is exhausted or
// an error is encountered.
func Lex(code string) (chan *Match, chan error) {
	tokens := make(chan *Match)
	errors := make(chan error)

	go func() {
		defer close(tokens)
		defer close(errors)

		for code != "" {
			var match []string
			for token := Token(0); token < Token(Unknown); token++ {
				code = strings.Trim(code, " ")
				if match = Rules[token].FindStringSubmatch(code); match != nil {
					tokens <- &Match{token, match}
					code = code[len(match[0]):]
					break
				}
			}
			if match == nil {
				errors <- fmt.Errorf("invalid syntax near '%s'", code)
				return
			}
		}
	}()

	return tokens, errors
}
