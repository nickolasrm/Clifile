// Package parser contains the semantic rules to parse the tokens coming from the lexer
// into a structured program.
package parser

import (
	"fmt"
	"regexp"

	"github.com/nickolasrm/clifile/internal/lexer"
)

// Variable is a struct that represents a variable.
// it contains the name of the variable and the value it holds
type Variable struct {
	Name  string
	Value string
}

func parseVariable(match *lexer.Match) (*Variable, error) {
	name := match.Value[1]
	value := match.Value[2]
	if value == "" {
		value = match.Value[3]
	}
	return &Variable{Name: name, Value: value}, nil
}

// Call is a struct that represents a function call.
// it contains the name of the function, the function itself
// and the arguments that the call received
type Call struct {
	Name      string
	Function  string
	Arguments map[string]string
}

func parseCall(match *lexer.Match) (*Call, error) {
	name := match.Value[1]
	function := match.Value[2]
	argumentString := match.Value[3]
	argumentTokens, err := lexer.Lex(argumentString)
	if err != nil {
		return nil, err
	}
	keywords := make(map[string]string)
	for _, submatch := range argumentTokens {
		switch submatch.Type {
		case lexer.Variable:
			variable, err := parseVariable(submatch)
			if err != nil {
				return nil, fmt.Errorf(
					"invalid parameter near '%s'",
					submatch.Value[0],
				)
			}
			keywords[variable.Name] = variable.Value
		case lexer.Line, lexer.Comment, lexer.Indent:
			continue
		default:
			return nil, fmt.Errorf(
				"unexpected syntax inside function call near '%s'",
				submatch.Value[0],
			)
		}
	}
	return &Call{
		Name:      name,
		Function:  function,
		Arguments: keywords,
	}, nil
}

// Rule is a struct that represents a group of actions and its properties.
// it contains the name of the rule, the positional arguments order,
// the docstring, the actions the rule will execute and its child rules if it is a group
type Rule struct {
	Name       string
	Positional []string
	Doc        string
	Actions    string
	Rules      map[string]*Rule
}

// Program is a struct that represents the entire program.
// it contains the variables, the function calls and the rules
type Program struct {
	Name      string
	Doc       string
	Variables map[string]*Variable
	Calls     map[string]*Call
	Rules     map[string]*Rule
}

func (p *Program) parseMetadata(matches []*lexer.Match) []*lexer.Match {
	if matches[0].Type == lexer.Docstring {
		p.Name = matches[0].Value[1]
		matches = matches[1:]
		programDoc := ""
		var i int
		var match *lexer.Match
		lines := 0
		for i, match = range matches {
			switch match.Type {
			case lexer.Docstring:
				if lines > 1 {
					goto Exit
				}
				programDoc += match.Value[1] + "\n"
				lines = 0
			case lexer.Line:
				lines = len(match.Value[0])
				continue
			default:
				goto Exit
			}
		}
	Exit:
		if programDoc != "" {
			p.Doc = programDoc
			matches = matches[i-1:]
		}
	}
	return matches
}

// Parse parses a list of lexer matches into a structured program.
// it returns a pointer to a Program struct and an error if any
// semantic error is found
func Parse(matches []*lexer.Match) (*Program, error) {
	program := &Program{
		Name:      "Software Command Line Interface (CLI)",
		Doc:       "Use this as shortcut for user-defined commands",
		Variables: make(map[string]*Variable),
		Calls:     make(map[string]*Call),
		Rules:     make(map[string]*Rule),
	}
	matches = program.parseMetadata(matches)

	indent := 0
	queue := make([]*Rule, 0)
	docstring := ""
	var parent *Rule = nil
	var rule *Rule = nil

	for _, match := range matches {
		if indent > 0 {
			switch match.Type {
			case lexer.Variable, lexer.Call:
				match = &lexer.Match{
					Type:  lexer.Action,
					Value: []string{match.Value[0]},
				}
			default:
				break
			}
		}
		switch match.Type {
		case lexer.Line:
			indent = 0
			continue
		case lexer.Comment:
			continue
		case lexer.Indent:
			indent = len(match.Value[0])
		case lexer.Variable:
			variable, err := parseVariable(match)
			if err != nil {
				return nil, err
			}
			program.Variables[variable.Name] = variable
		case lexer.Call:
			call, err := parseCall(match)
			if err != nil {
				return nil, err
			}
			program.Calls[call.Name] = call
		case lexer.Docstring:
			if indent > len(queue) {
				return nil, fmt.Errorf(
					"overly indented docstring near '%s'",
					match.Value[0],
				)
			}
			docstring += match.Value[1] + "\n"
		case lexer.Rule:
			if indent > len(queue) {
				return nil, fmt.Errorf(
					"overly indented rule near '%s'",
					match.Value[0],
				)
			}
			queue = queue[:indent]
			if indent == 0 {
				parent = nil
			} else {
				parent = queue[len(queue)-1]
			}
			rule = &Rule{
				Name:       "",
				Positional: nil,
				Doc:        docstring,
				Actions:    "",
				Rules:      make(map[string]*Rule),
			}
			queue = append(queue, rule)
			docstring = ""
			name := match.Value[1]
			rule.Name = name
			rule.Positional = regexp.MustCompile(`\w+`).FindAllString(
				match.Value[2], -1,
			)
			if parent != nil {
				if parent.Actions != "" {
					return nil, fmt.Errorf(
						"can't add nested rules into a rule '%s' because it has actions",
						parent.Name,
					)
				}
				parent.Rules[name] = rule
			} else {
				program.Rules[name] = rule
			}
		case lexer.Action:
			if indent < len(queue) {
				return nil, fmt.Errorf("bad indentation near '%s'", match.Value[0])
			}
			if rule == nil {
				return nil, fmt.Errorf(
					"action outside of rule near '%s'",
					match.Value[0],
				)
			}
			rule.Actions += match.Value[0] + "\n"
		}
	}

	return program, nil
}
