package runner

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nickolasrm/clifile/internal/interpreter"
	"github.com/nickolasrm/clifile/internal/lexer"
	"github.com/nickolasrm/clifile/internal/parser"
)

func TryThrow(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ReadCode() (string, error) {
	body, err := ioutil.ReadFile("Clifile")
	if err != nil {
		return "", fmt.Errorf("unable to read 'Clifile'")
	}
	return string(body), nil
}

func Run(code string) {
	tokens, err := lexer.Lex(code)
	TryThrow(err)
	program, err := parser.Parse(tokens)
	TryThrow(err)
	exec := interpreter.Interpret(program)
	exec.Run()
}

func RunFile() {
	code, err := ReadCode()
	TryThrow(err)
	Run(code)
}
