package runner

import (
	"fmt"
	"io/ioutil"

	"github.com/nickolasrm/clifile/internal/interpreter"
	"github.com/nickolasrm/clifile/internal/lexer"
	"github.com/nickolasrm/clifile/internal/parser"
	"github.com/nickolasrm/clifile/pkg/util"
)

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
	util.TryThrow(err)
	program, err := parser.Parse(tokens)
	util.TryThrow(err)
	exec := interpreter.Interpret(program)
	exec.Run()
}

// RunFile interprets and executes the Clifile
func RunFile() {
	code, err := ReadCode()
	util.TryThrow(err)
	Run(code)
}
