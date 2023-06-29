package interpreter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nickolasrm/clifile/internal/interpreter/flag"
	"github.com/nickolasrm/clifile/internal/parser"
	"github.com/nickolasrm/clifile/pkg/util"
	"github.com/spf13/cobra"
)

type Execution struct {
	Program *parser.Program
}

func (e *Execution) buildGroup(rule *parser.Rule) *cobra.Command {
	short := strings.SplitN(rule.Doc, "\n", 2)[0]
	return &cobra.Command{
		Use:   rule.Name,
		Short: short,
		Long:  rule.Doc,
	}
}

func (e *Execution) findReplacements(code string) []string {
	pattern := regexp.MustCompile(`(\$+)\{(\w+)\}`)
	matches := pattern.FindAllStringSubmatch(code, -1)
	unique := make(map[string]bool)
	replacements := make([]string, 0)
	for _, match := range matches {
		if len(match[1])%2 == 1 && !unique[match[2]] {
			unique[match[2]] = true
			replacements = append(replacements, match[2])
		}
	}
	return replacements
}

func (e *Execution) replaceVariables(actions string, variables []string) (string, error) {
	for _, name := range variables {
		if variable, ok := e.Program.Variables[name]; ok {
			actions = strings.ReplaceAll(actions, "${"+name+"}", variable.Value)
		} else {
			return "", fmt.Errorf("could not replacement for '%s'", name)
		}
	}
	actions = strings.ReplaceAll(actions, "$$", "$")
	return actions, nil
}

func (e *Execution) buildFlags(cmd *cobra.Command, flags []string) ([]*flag.Flag, error) {
	replacements := make([]*flag.Flag, 0)
	for _, name := range flags {
		if call, ok := e.Program.Calls[name]; ok && call.Function == "flag" {
			flag, err := flag.CreateFlag(
				call.Name,
				call.Arguments["doc"],
				call.Arguments["default"],
				call.Arguments["type"],
				call.Arguments["question"],
				call.Arguments["validation"],
				call.Arguments["init"],
				cmd,
			)
			if err != nil {
				return nil, err
			}
			replacements = append(replacements, flag)
		}
	}
	return replacements, nil
}

func (e *Execution) buildCommand(rule *parser.Rule) (*cobra.Command, error) {
	cmd := e.buildGroup(rule)
	replacements := e.findReplacements(rule.Actions)
	flags, err := e.buildFlags(cmd, replacements)
	if err != nil {
		return nil, err
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var value string
		var err error
		for _, flag := range flags {
			value, err = flag.Get()
			if err != nil {
				return err
			}
			e.Program.Variables[flag.Replacement()] = &parser.Variable{
				Name:  flag.Replacement(),
				Value: value,
			}
		}
		actions := rule.Actions
		if actions, err = e.replaceVariables(actions, replacements); err != nil {
			return err
		}
		if err = util.Shell(actions); err != nil {
			return err
		}
		return nil
	}
	return cmd, nil
}

func (e *Execution) buildMain() (*cobra.Command, error) {
	main := &cobra.Command{
		Use:  "cli",
		Long: e.Program.Doc,
	}
	var recursion func(*cobra.Command, *parser.Rule) error
	recursion = func(parent *cobra.Command, rule *parser.Rule) error {
		var err error
		var child *cobra.Command
		if len(rule.Rules) == 0 {
			if child, err = e.buildCommand(rule); err != nil {
				return err
			}
		} else {
			child = e.buildGroup(rule)
			for _, childRule := range rule.Rules {
				if err = recursion(child, childRule); err != nil {
					return err
				}
			}
		}
		parent.AddCommand(child)
		return nil
	}
	for _, rule := range e.Program.Rules {
		if err := recursion(main, rule); err != nil {
			return nil, err
		}
	}
	return main, nil
}

func (e *Execution) Run() {
	main, err := e.buildMain()
	util.TryThrow(err)
	err = main.Execute()
	util.TryThrow(err)
}

func Interpret(program *parser.Program) *Execution {
	return &Execution{
		Program: program,
	}
}
