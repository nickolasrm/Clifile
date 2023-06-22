package parser

import (
	"fmt"
	"regexp"

	"github.com/nickolasrm/clifile/internal/lexer"
)

// Variable is a struct that represents a variable
// It contains the name of the variable and the value it holds
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

// Call is a struct that represents a function call
// It contains the name of the function, the function itself
// and the parameters that the call received
type Call struct {
	Name       string
	Function   string
	Parameters map[string]string
}

func parseCall(match *lexer.Match) (*Call, error) {
	name := match.Value[1]
	function := match.Value[2]
	parametersString := match.Value[3]
	paramTokens, err := lexer.Lex(parametersString)
	if err != nil {
		return nil, err
	}
	keywords := make(map[string]string)
	for _, submatch := range paramTokens {
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
		Name:       name,
		Function:   function,
		Parameters: keywords,
	}, nil
}

// Rule is a struct that represents a group of actions and its properties
// It contains the name of the rule, the positional arguments order
// the docstring, the actions the rule will execute and its child rules if it is a group
type Rule struct {
	Name       string
	Positional []string
	Docstring  string
	Actions    string
	Rules      map[string]*Rule
}

// Program is a struct that represents the entire program
// It contains the variables, the function calls and the rules
type Program struct {
	Variables map[string]*Variable
	Calls     map[string]*Call
	Rules     map[string]*Rule
}

// Parse parses a list of lexer matches into a structured program
// It returns a pointer to a Program struct and an error if any
// semantic error is found
func Parse(matches []*lexer.Match) (*Program, error) {
	program := &Program{
		Variables: make(map[string]*Variable),
		Calls:     make(map[string]*Call),
		Rules:     make(map[string]*Rule),
	}

	indent := 0
	queue := make([]*Rule, 0)
	docstring := ""
	var parent *Rule = nil
	var rule *Rule = nil

	for _, match := range matches {
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
				return nil, fmt.Errorf("overly indented docstring near '%s'", match.Value[0])
			}
			docstring += match.Value[1] + "\n"
		case lexer.Rule:
			if indent > len(queue) {
				return nil, fmt.Errorf("overly indented rule near '%s'", match.Value[0])
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
				Docstring:  docstring,
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
					return nil, fmt.Errorf("can't add nested rules into a rule '%s' because it has actions", parent.Name)
				}
				parent.Rules[name] = rule
			} else {
				program.Rules[name] = rule
			}
		case lexer.Action:
			if indent < len(queue) {
				return nil, fmt.Errorf("bad identation near '%s'", match.Value[0])
			}
			if rule == nil {
				return nil, fmt.Errorf("action outside of rule near '%s'", match.Value[0])
			}
			rule.Actions += match.Value[1] + "\n"
		}
	}

	return program, nil
}
