// Package call contains the struct that represents a function call
package call

// Call is a struct that represents a function call.
// it contains the name of the function, the function itself
// and the arguments that the call received
type Call struct {
	name      string
	function  string
	arguments map[string]string
}

// NewCall is a helper function to create a new call
func NewCall(name, function string, arguments map[string]string) *Call {
	return &Call{
		name:      name,
		function:  function,
		arguments: arguments,
	}
}

// Name returns the name of the call
func (c *Call) Name() string {
	return c.name
}

// Function returns the name of the function of the call
func (c *Call) Function() string {
	return c.function
}

// Arguments returns the arguments of the call
func (c *Call) Arguments() map[string]string {
	return c.arguments
}

// Argument returns the value of an argument of the call
func (c *Call) Argument(arg string) string {
	return c.arguments[arg]
}
