// Package lexer contains all the syntatic rules and methods for
// tokenizing a string into a stream of tokens.
package lexer

import (
	"fmt"
	"regexp"
	"strings"
)

// Token is an enum that represents the type of token that was matched.
type Token uint8

const (
	// Line represents a line break or a space.
	Line = iota
	// Indent represents a tab.
	Indent
	// Docstring represents a documentation comment.
	Docstring
	// Comment represents an ignorable string.
	Comment
	// Call represents a function along with its arguments.
	Call
	// Variable represents an identified data.
	Variable
	// Rule represents a rule a group of rules or a sequence of commands.
	Rule
	// Action represents a command that will be executed.
	Action
	// Unknown represents an not parsed token.
	Unknown
)

// ruleRegex returns a regex that matches the pattern at the beginning of a string.
func ruleRegex(pattern string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(`^%s`, pattern))
}

// Rules is a map of tokens to regexes that identify them.
var Rules = map[Token]*regexp.Regexp{
	Line:      ruleRegex(`[\n ]+`),
	Indent:    ruleRegex(`\t+`),
	Docstring: ruleRegex(`##[ ]?([^\n]*)`),
	Comment:   ruleRegex(`#([^\n]*)`),
	Call:      ruleRegex(`(\w+)[\t ]*=[\t ]*\${(?:\s+)?(\w+)(?:\s+)?([^}]+)?}`),
	Variable:  ruleRegex(`(\w+)[\t ]*=(?:[\t ]*"([^"]*)"|([^"\n]*))`),
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
