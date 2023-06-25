package interpreter

import (
	"fmt"
	"os"
	"strings"

	"github.com/nickolasrm/clifile/internal/parser"
	"github.com/nickolasrm/clifile/package/util"
	"github.com/spf13/cobra"
)

type Flag struct {
	Name    string
	Default string
	Type    string
	Prompt  string
}

type Execution struct {
	Program *parser.Program
}

func (e *Execution) GetVariableValue(name string) (string, error) {
	variable, ok := e.Program.Variables[name]
	if !ok {
		return "", fmt.Errorf("variable '%s' not found", name)
	}
	return variable.Value, nil
}

func (e *Execution) GetFlag(name string) (*Flag, error) {
	call, ok := e.Program.Calls[name]
	if !ok {
		return nil, fmt.Errorf("flag '%s' not found", name)
	}
	if call.Function != "flag" {
		return nil, fmt.Errorf("call '%s' is not a flag", name)
	}
	flag := Flag{
		Name:    name,
		Default: call.Arguments["default"],
		Type:    call.Arguments["type"],
		Prompt:  call.Arguments["prompt"],
	}
	return &flag, nil
}

func (e *Execution) buildGroup(rule *parser.Rule) *cobra.Command {
	short := strings.SplitN(rule.Doc, "\n", 2)[0]
	return &cobra.Command{
		Use:   rule.Name,
		Short: short,
		Long:  rule.Doc,
	}
}

func (e *Execution) buildCommand(rule *parser.Rule) *cobra.Command {
	cmd := e.buildGroup(rule)
	cmd.Run = func(cmd *cobra.Command, args []string) {
		util.Shell(rule.Actions)
	}
	return cmd
}

func (e *Execution) buildMain() *cobra.Command {
	main := &cobra.Command{
		Use:  "cli",
		Long: e.Program.Doc,
	}
	var recursion func(*cobra.Command, *parser.Rule)
	recursion = func(parent *cobra.Command, rule *parser.Rule) {
		var child *cobra.Command
		if len(rule.Rules) == 0 {
			child = e.buildCommand(rule)
		} else {
			child = e.buildGroup(rule)
			for _, childRule := range rule.Rules {
				recursion(child, childRule)
			}
		}
		parent.AddCommand(child)
	}
	for _, rule := range e.Program.Rules {
		recursion(main, rule)
	}
	return main
}

func (e *Execution) Run() {
	main := e.buildMain()
	if err := main.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Interpret(program *parser.Program) *Execution {
	return &Execution{
		Program: program,
	}
}
