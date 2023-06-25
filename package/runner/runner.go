package runner

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nickolasrm/clifile/internal/interpreter"
	"github.com/nickolasrm/clifile/internal/lexer"
	"github.com/nickolasrm/clifile/internal/parser"
)

// TryThrow checks if an error is not nil and then prints it to stderr and
// finished the program execution
func TryThrow(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// ReadCode reads the Clifile reference and returns its contents or an error
func ReadCode() (string, error) {
	body, err := ioutil.ReadFile("Clifile")
	if err != nil {
		return "", fmt.Errorf("unable to read 'Clifile'")
	}
	return string(body), nil
}

// Run interprets and executes a Clifile syntax code string
func Run(code string) {
	tokens, err := lexer.Lex(code)
	TryThrow(err)
	program, err := parser.Parse(tokens)
	TryThrow(err)
	exec := interpreter.Interpret(program)
	exec.Run()
}

// RunFile interprets and executes the Clifile
func RunFile() {
	code, err := ReadCode()
	TryThrow(err)
	Run(code)
}
