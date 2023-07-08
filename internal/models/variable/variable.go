// Package variable contains the struct that represents a key value pair
package variable

// Variable is a struct that represents a variable.
// it contains the name of the variable and the value it holds
type Variable struct {
	name  string
	value string
}

// NewVariable is a helper function to create a new variable
func NewVariable(name, value string) *Variable {
	return &Variable{name: name, value: value}
}

// Name returns the name of the variable
func (v *Variable) Name() string {
	return v.name
}

// Value returns the value of the variable
func (v *Variable) Value() string {
	return v.value
}

// SetValue sets the value of the variable
func (v *Variable) SetValue(value string) {
	v.value = value
}
