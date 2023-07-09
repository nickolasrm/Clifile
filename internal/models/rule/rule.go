package rule

import "golang.org/x/exp/maps"

// Rule is a struct that represents a group of actions and its properties.
// it contains the name of the rule, the positional arguments order,
// the docstring, the actions the rule will execute and its child rules if it is a group
type Rule struct {
	name       string
	positional []string
	doc        string
	actions    string
	rules      map[string]*Rule
}

// NewRule is a helper function to create a new rule
func NewRule(name string, positional []string, doc, actions string) *Rule {
	r := &Rule{
		name:       name,
		positional: positional,
		doc:        doc,
		actions:    "",
		rules:      make(map[string]*Rule),
	}
	if actions != "" {
		r.AppendActions(actions)
	}
	return r
}

// Name returns the name of the rule
func (r *Rule) Name() string {
	return r.name
}

// Positional returns the positional arguments of the rule
func (r *Rule) Positional() []string {
	return r.positional
}

// Doc returns the docstring of the rule
func (r *Rule) Doc() string {
	return r.doc
}

// Actions returns the actions of the rule
func (r *Rule) Actions() string {
	return r.actions
}

// AppendActions appends actions to the actions of the rule
func (r *Rule) AppendActions(actions string) {
	r.actions += actions + "\n"
}

// Rules returns the child rules of the rule
func (r *Rule) Rules() []*Rule {
	return maps.Values(r.rules)
}

// Rule returns a child rule of the rule
// it returns nil if the rule does not exist
func (r *Rule) Rule(name string) *Rule {
	return r.rules[name]
}

// SetRule sets a child rule of the rule
func (r *Rule) AddRule(rule *Rule) {
	r.rules[rule.Name()] = rule
}
