package interpreter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nickolasrm/clifile/internal/parser"
	"github.com/nickolasrm/clifile/pkg/util"
	"github.com/spf13/cobra"
)

type FlagType uint8

const (
	FlagString = FlagType(iota)
	FlagText
	FlagConfirm
	FlagSelect
)

var FlagTypeMap = map[string]FlagType{
	"string":    FlagString,
	"multiline": FlagText,
	"confirm":   FlagConfirm,
	"select":    FlagSelect,
}

// TODO: move to another file
type Flag struct {
	Name       string
	Doc        string
	Default    string
	Type       FlagType
	Question   string
	Validation *regexp.Regexp
	Init       string
}

func createValidator(validation string, flagType FlagType) (*regexp.Regexp, error) {
	switch flagType {
	case FlagConfirm:
		return regexp.Compile(fmt.Sprintf("%v|%v", true, false))
	default:
		if validation == "" {
			return regexp.Compile(`.*`)
		}
	}
	return regexp.Compile(validation)
}

func createFlag(
	name string,
	doc string,
	defaultValue string,
	flagType string,
	question string,
	validation string,
	init string,
) (*Flag, error) {
	if flagType == "" {
		flagType = "string"
	}
	parsedType, ok := FlagTypeMap[flagType]
	if !ok {
		return nil, fmt.Errorf("invalid flag type '%s'", flagType)
	}
	validator, err := createValidator(validation, parsedType)
	if err != nil {
		return nil, err
	}
	return &Flag{
		Name:       name,
		Doc:        doc,
		Default:    defaultValue,
		Type:       parsedType,
		Question:   question,
		Validation: validator,
		Init:       init,
	}, nil
}

func (f *Flag) docstring() string {
	if f.Default != "" {
		return fmt.Sprintf("%s (default: %s)", f.Doc, f.Default)
	}
	return f.Doc
}

func (f *Flag) validate(value string) error {
	if !f.Validation.MatchString(value) {
		return fmt.Errorf("invalid value '%s' for flag '%s'", value, f.Name)
	}
	return nil
}

func (f *Flag) prompt() string {
	if f.Question != "" {
		return f.Question
	}
	return f.Name
}

func (f *Flag) ask() string {
	var result interface{}
	var prompt survey.Prompt
	switch f.Type {
	case FlagString:
		prompt = &survey.Input{
			Help:    f.Doc,
			Message: f.prompt(),
		}
	case FlagText:
		prompt = &survey.Multiline{
			Help:    f.Doc,
			Message: f.prompt(),
		}
	case FlagConfirm:
		prompt = &survey.Confirm{
			Help:    f.Doc,
			Message: f.prompt(),
		}
	}
	if exited := survey.AskOne(prompt, &result); exited != nil {
		util.TryThrow(exited)
	}
	return fmt.Sprintf("%v", result)
}

func (f *Flag) Run(value string, exist bool) (string, error) {
	if exist {
		return value, f.validate(value)
	} else {
		value = f.Default
		for f.validate(value) != nil {
			value = f.ask()
		}
	}
	return value, nil
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

// TODO: Return set
func (e *Execution) findReplacements(code string) []string {
	pattern := regexp.MustCompile(`(\$+)\{(\w+)\}`)
	matches := pattern.FindAllStringSubmatch(code, -1)
	replacements := make([]string, 0)
	for _, match := range matches {
		if len(match[1])%2 == 1 {
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

func (e *Execution) buildFlags(cmd *cobra.Command, flags []string) ([]*Flag, error) {
	replacements := make([]*Flag, 0)
	for _, name := range flags {
		if call, ok := e.Program.Calls[name]; ok && call.Function == "flag" {
			flag, err := createFlag(
				call.Name,
				call.Arguments["doc"],
				call.Arguments["default"],
				call.Arguments["type"],
				call.Arguments["question"],
				call.Arguments["validation"],
				call.Arguments["init"],
			)
			if err != nil {
				return nil, err
			}
			replacements = append(replacements, flag)
			cmd.Flags().String(flag.Name, "", flag.docstring())
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
			value, _ = cmd.Flags().GetString(flag.Name)
			exist := cmd.Flags().Lookup(flag.Name).Changed
			value, err = flag.Run(value, exist)
			if err != nil {
				return err
			}
			e.Program.Variables[flag.Name] = &parser.Variable{
				Name:  flag.Name,
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
