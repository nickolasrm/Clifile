// Package parser contains the semantic rules to parse the tokens coming from the lexer
// into a structured program.
package parser

import (
	"fmt"
	"regexp"

	"github.com/nickolasrm/clifile/internal/lexer"
	"github.com/nickolasrm/clifile/internal/models/call"
	"github.com/nickolasrm/clifile/internal/models/program"
	"github.com/nickolasrm/clifile/internal/models/rule"
	"github.com/nickolasrm/clifile/internal/models/variable"
)

// Parser is a struct that represents the parser
type Parser struct {
	program *program.Program
}

// NewParser is a helper function to create a new parser
func NewParser() *Parser {
	return &Parser{
		program: program.NewProgram(),
	}
}

// parseDoc parses the documentation of the program if it exists
func (p *Parser) parseDoc(tokens []*lexer.Match) []*lexer.Match {
	if tokens[0].Type() == lexer.Docstring {
		programDoc := ""
		var i int
		var match *lexer.Match
		lines := 0
		for i, match = range tokens {
			switch match.Type() {
			case lexer.Docstring:
				if lines > 1 {
					goto Exit
				}
				programDoc += match.Value(1) + "\n"
				lines = 0
			case lexer.Line:
				lines = len(match.Value(0))
				continue
			default:
				goto Exit
			}
		}
	Exit:
		if programDoc != "" {
			p.program.SetDoc(programDoc)
			tokens = tokens[i:]
		}
	}
	return tokens
}

// parseVariable creates a variable from a variable lexer match
func parseVariable(match *lexer.Match) (*variable.Variable, error) {
	name := match.Value(1)
	value := match.Value(2)
	if value == "" {
		value = match.Value(3)
	}
	return variable.NewVariable(name, value), nil
}

// parseCall creates a call from a call lexer match
func parseCall(match *lexer.Match) (*call.Call, error) {
	name := match.Value(1)
	function := match.Value(2)
	argumentString := match.Value(3)
	argumentTokens, err := lexer.Lex(argumentString)
	if err != nil {
		return nil, err
	}
	keywords := make(map[string]string)
	for _, submatch := range argumentTokens {
		switch submatch.Type() {
		case lexer.Variable:
			variable, err := parseVariable(submatch)
			if err != nil {
				return nil, fmt.Errorf(
					"invalid parameter near '%s'",
					submatch.Value(0),
				)
			}
			keywords[variable.Name()] = variable.Value()
		case lexer.Line, lexer.Comment, lexer.Indent:
			continue
		default:
			return nil, fmt.Errorf(
				"unexpected syntax inside function call near '%s'",
				submatch.Value(0),
			)
		}
	}
	return call.NewCall(name, function, keywords), nil
}

// parseBody parses the body of the program and creates the rules
// and variables
func (p *Parser) parseBody(tokens []*lexer.Match) error {
	program := p.program
	indent := 0
	ruleBranch := make([]*rule.Rule, 0)
	ruleDoc := ""
	var parentRule *rule.Rule = nil
	var currentRule *rule.Rule = nil

	for _, match := range tokens {
		if indent > 0 {
			switch match.Type() {
			case lexer.Variable, lexer.Call:
				match = lexer.NewMatch(lexer.Action, []string{match.Value(0)})
			default:
				break
			}
		}
		switch match.Type() {
		case lexer.Line:
			indent = 0
			continue
		case lexer.Comment:
			continue
		case lexer.Indent:
			indent = len(match.Value(0))
		case lexer.Variable:
			variable, err := parseVariable(match)
			if err != nil {
				return err
			}
			program.AddVariable(variable)
		case lexer.Call:
			call, err := parseCall(match)
			if err != nil {
				return err
			}
			program.AddCall(call)
		case lexer.Docstring:
			if indent > len(ruleBranch) {
				return fmt.Errorf(
					"overly indented docstring near '%s'",
					match.Value(0),
				)
			}
			ruleDoc += match.Value(1) + "\n"
		case lexer.Rule:
			if indent > len(ruleBranch) {
				return fmt.Errorf(
					"overly indented rule near '%s'",
					match.Value(0),
				)
			}
			ruleBranch = ruleBranch[:indent]
			if indent == 0 {
				parentRule = nil
			} else {
				parentRule = ruleBranch[len(ruleBranch)-1]
			}
			name := match.Value(1)
			positionals := regexp.MustCompile(`\w+`).FindAllString(
				match.Value(2), -1,
			)
			currentRule = rule.NewRule(name, positionals, ruleDoc, "")
			ruleBranch = append(ruleBranch, currentRule)

			if parentRule != nil {
				if parentRule.Actions() != "" {
					return fmt.Errorf(
						"can't add nested rules into a rule '%s' because it has actions",
						parentRule.Name(),
					)
				}
				parentRule.AddRule(currentRule)
			} else {
				program.AddRule(currentRule)
			}
		case lexer.Action:
			if indent < len(ruleBranch) {
				return fmt.Errorf("bad indentation near '%s'", match.Value(0))
			}
			if currentRule == nil {
				return fmt.Errorf(
					"action outside of rule near '%s'",
					match.Value(0),
				)
			}
			currentRule.AppendActions(match.Value(0))
		}
	}
	return nil
}

// Parse parses a list of lexer matches into a structured program.
func (p *Parser) Parse(tokens []*lexer.Match) (*program.Program, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens to parse")
	}
	tokens = p.parseDoc(tokens)
	err := p.parseBody(tokens)
	if err != nil {
		return nil, err
	}
	return p.program, nil
}

// Parse parses a list of lexer matches into a structured program.
// it returns a pointer to a Program struct and an error if any
// semantic error is found
func Parse(tokens []*lexer.Match) (*program.Program, error) {
	parser := NewParser()
	return parser.Parse(tokens)
}
