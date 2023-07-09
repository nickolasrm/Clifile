// Package program contains the structs that represent the program
// components
package program

import (
	"github.com/nickolasrm/clifile/internal/models/call"
	"github.com/nickolasrm/clifile/internal/models/rule"
	"github.com/nickolasrm/clifile/internal/models/variable"
	"golang.org/x/exp/maps"
)

// Program is a struct that represents the entire program.
// it contains the variables, the function calls and the rules
type Program struct {
	doc       string
	variables map[string]*variable.Variable
	calls     map[string]*call.Call
	rules     map[string]*rule.Rule
}

// NewProgram is a helper function to create a new program
func NewProgram() *Program {
	return &Program{
		doc: `Software Command Line Interface (CLI)
Use this as shortcut for user-defined commands`,
		variables: make(map[string]*variable.Variable),
		calls:     make(map[string]*call.Call),
		rules:     make(map[string]*rule.Rule),
	}
}

// Doc returns the documentation of the program
func (p *Program) Doc() string {
	return p.doc
}

// SetDoc sets the documentation of the program
func (p *Program) SetDoc(doc string) {
	p.doc = doc
}

// Call returns a call from the program
// it returns nil if the call does not exist
func (p *Program) Call(name string) *call.Call {
	return p.calls[name]
}

// AddCall adds a call to the program
func (p *Program) AddCall(c *call.Call) {
	p.calls[c.Name()] = c
}

// Rules returns the rules of the program
func (p *Program) Rules() []*rule.Rule {
	return maps.Values(p.rules)
}

// Rule returns a rule from the program
// it returns nil if the rule does not exist
func (p *Program) Rule(name string) *rule.Rule {
	return p.rules[name]
}

// AddRule sets the rule of a call or creates a new one if it does not exist
func (p *Program) AddRule(rule *rule.Rule) {
	p.rules[rule.Name()] = rule
}

// Variable returns a variable from the program
// it returns nil if the variable does not exist
func (p *Program) Variable(name string) *variable.Variable {
	return p.variables[name]
}

// SetVariable sets the value of a variable or creates a new one if it does not exist
func (p *Program) SetVariable(name string, value string) {
	if v := p.Variable(name); v != nil {
		v.SetValue(value)
	} else {
		p.variables[name] = variable.NewVariable(name, value)
	}
}

// AddVariable adds a variable to the program
func (p *Program) AddVariable(v *variable.Variable) {
	p.variables[v.Name()] = v
}
