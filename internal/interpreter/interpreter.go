package interpreter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nickolasrm/clifile/internal/parser"
	"github.com/nickolasrm/clifile/package/util"
	"github.com/spf13/cobra"
)

type Flag struct {
	Name       string
	Doc        string
	Default    string
	Type       string
	Prompt     string
	Validation string
	Init       string
}

func (f *Flag) Run(value string, exist bool) string {
	if !exist {

	}
	return "TODO"
}

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
	replacements := make([]string, 0)
	for _, match := range matches {
		if len(match[1]) % 2 == 1 {
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

func (e *Execution) buildFlags(cmd *cobra.Command, flags []string) []*Flag {
	replacements := make([]*Flag, 0)
	for _, name := range flags {
		if call, ok := e.Program.Calls[name]; ok && call.Function == "flag" {
			flag := &Flag{
				Name:       call.Name,
				Doc:        call.Arguments["doc"],
				Default:    call.Arguments["default"],
				Type:       call.Arguments["type"],
				Prompt:     call.Arguments["prompt"],
				Validation: call.Arguments["validation"],
			}
			replacements = append(replacements, flag)
			cmd.Flags().String(flag.Name, flag.Default, flag.Doc)
		}
	}
	return replacements
}

func (e *Execution) buildCommand(rule *parser.Rule) *cobra.Command {
	cmd := e.buildGroup(rule)
	replacements := e.findReplacements(rule.Actions)
	flags := e.buildFlags(cmd, replacements)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var err error
		for _, flag := range flags {
			value, parseError := cmd.Flags().GetString(flag.Name)
			exist := parseError == nil
			e.Program.Variables[flag.Name] = &parser.Variable{
				Name:  flag.Name,
				Value: flag.Run(value, exist),
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
	return cmd
}

func (e *Execution) buildMain() (*cobra.Command, error) {
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
