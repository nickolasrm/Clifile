// Lexer package contains all the syntatic rules and methods for
// tokenizing a string into a stream of tokens.
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
	Call
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
	Line:      ruleRegex(`[\n ]+`),
	Indent:    ruleRegex(`\t+`),
	Docstring: ruleRegex(`##([^\n]*)`),
	Comment:   ruleRegex(`#([^\n]*)`),
	Call:      ruleRegex(`(\w+)=\${(?:\s+)?(\w+)(?:\s+)?([^}]+)?}`),
	Variable:  ruleRegex(`(\w+)=(?:"([^"]*)"|([^"\n]*))`),
	Rule:      ruleRegex(`(\w+):([\w ]*)`),
	Action:    ruleRegex(`[^\n]+`),
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
func Lex(code string) ([]*Match, error) {
	tokens := make([]*Match, 0)
	for code != "" {
		var match []string
		for token := Token(0); token < Token(Unknown); token++ {
			code = strings.Trim(code, " ")
			if match = Rules[token].FindStringSubmatch(code); match != nil {
				tokens = append(tokens, &Match{token, match})
				code = code[len(match[0]):]
				break
			}
		}
		if match == nil {
			return nil, fmt.Errorf("invalid syntax near '%s'", code)
		}
	}
	return tokens, nil
}
