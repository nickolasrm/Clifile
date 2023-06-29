// Package flag contains all the syntatic rules and methods for building a flag
// from the function call arguments.
package flag

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nickolasrm/clifile/pkg/util"
	"github.com/spf13/cobra"
)

// FlagData represents the data needed to build a flag.
type FlagData struct {
	name         string
	defaultValue string
	Type         string
	doc          string
	question     string
	validation   string
	init         string
}

// Name returns the parsed name of the flag.
func (f *FlagData) Name() string {
	return strings.ToLower(strings.ReplaceAll(f.name, "_", "-"))
}

// Doc returns the flag documentation.
func (f *FlagData) Doc() string {
	if f.defaultValue != "" {
		return fmt.Sprintf("%s (default: %s)", f.doc, f.defaultValue)
	}
	return f.doc
}

// Question returns the text to show to the user when prompting for the flag value.
func (f *FlagData) Question() string {
	if f.question != "" {
		return f.question
	}
	return f.Name()
}

// FlagLogic is an interface that represents the logic needed to build a flag.
type FlagLogic interface {
	// Prompt defines what kind of prompt will be used an how
	Prompt(*FlagData) survey.Prompt
	// Validator defines the default regexp for validation
	Validator(*FlagData) *regexp.Regexp
	// Parse defines how to process the output of the survey
	Parse(*FlagData, interface{}) string
}

// FlagLogicInput creates an string input prompt.
type FlagLogicInput struct{}

func (f *FlagLogicInput) Prompt(data *FlagData) survey.Prompt {
	return &survey.Input{
		Message: data.Question(),
		Help:    data.Doc(),
	}
}

func (f *FlagLogicInput) Validator(data *FlagData) *regexp.Regexp {
	return regexp.MustCompile(`.*`)
}

func (f *FlagLogicInput) Parse(data *FlagData, value interface{}) string {
	return value.(string)
}

// FlagLogicConfirm creates a confirmation prompt
type FlagLogicConfirm struct{}

func (f *FlagLogicConfirm) Prompt(data *FlagData) survey.Prompt {
	return &survey.Confirm{
		Message: data.Question(),
		Help:    data.Doc(),
	}
}

func (f *FlagLogicConfirm) Validator(data *FlagData) *regexp.Regexp {
	return regexp.MustCompile(".*")
}

func (f *FlagLogicConfirm) Parse(data *FlagData, value interface{}) string {
	if value.(bool) {
		return "y"
	}
	return "n"
}

// FlagLogicMultiline creates a multiline string input
type FlagLogicMultiLine struct{}

func (f *FlagLogicMultiLine) Prompt(data *FlagData) survey.Prompt {
	return &survey.Multiline{
		Message: data.Question(),
		Help:    data.Doc(),
	}
}

func (f *FlagLogicMultiLine) Validator(data *FlagData) *regexp.Regexp {
	return regexp.MustCompile(`.*`)
}

func (f *FlagLogicMultiLine) Parse(data *FlagData, value interface{}) string {
	return value.(string)
}

// FlagLogicPasssword creates a hidden string input
type FlagLogicPassword struct{}

func (f *FlagLogicPassword) Prompt(data *FlagData) survey.Prompt {
	return &survey.Password{
		Message: data.Question(),
		Help:    data.Doc(),
	}
}

func (f *FlagLogicPassword) Validator(data *FlagData) *regexp.Regexp {
	return regexp.MustCompile(`.*`)
}

func (f *FlagLogicPassword) Parse(data *FlagData, value interface{}) string {
	return value.(string)
}

// FlagLogicSelect creates a selectable list of predefined options
type FlagLogicSelect struct{}

func (f *FlagLogicSelect) Prompt(data *FlagData) survey.Prompt {
	opts := strings.Split(data.init, ", ")
	return &survey.Select{
		Message: data.Question(),
		Help:    data.Doc(),
		Options: opts,
	}
}

func (f *FlagLogicSelect) Validator(data *FlagData) *regexp.Regexp {
	return regexp.MustCompile(`.*`)
}

func (f *FlagLogicSelect) Parse(data *FlagData, value interface{}) string {
	return value.(survey.OptionAnswer).Value
}

// FlagLogicEditor uses an editor to input a large string
type FlagLogicEditor struct{}

func (f *FlagLogicEditor) Prompt(data *FlagData) survey.Prompt {
	return &survey.Editor{
		Message:  data.Question(),
		Help:     data.Doc(),
		FileName: data.init,
	}
}

func (f *FlagLogicEditor) Validator(data *FlagData) *regexp.Regexp {
	return regexp.MustCompile(`.*`)
}

func (f *FlagLogicEditor) Parse(data *FlagData, value interface{}) string {
	return value.(string)
}

// FlagLogicMultiSelect creates a multi selectable list of options
type FlagLogicMultiSelect struct{}

func (f *FlagLogicMultiSelect) Prompt(data *FlagData) survey.Prompt {
	opts := strings.Split(data.init, ", ")
	return &survey.MultiSelect{
		Message: data.Question(),
		Help:    data.Doc(),
		Options: opts,
	}
}

func (f *FlagLogicMultiSelect) Validator(data *FlagData) *regexp.Regexp {
	return regexp.MustCompile(`.*`)
}

func (f *FlagLogicMultiSelect) Parse(data *FlagData, value interface{}) string {
	answer := ""
	for _, el := range value.([]survey.OptionAnswer) {
		answer += fmt.Sprintf("\"%s\" ", el.Value)
	}
	return strings.TrimRight(answer, " ")
}

// Flag is a struct that controls the logic for building a flag.
type Flag struct {
	Data    *FlagData
	Logic   FlagLogic
	Command *cobra.Command
}

// CreateFlag creates a new Flag from the function call arguments.
func CreateFlag(
	name string,
	doc string,
	defaultValue string,
	flagType string,
	question string,
	validation string,
	init string,
	command *cobra.Command,
) (*Flag, error) {
	data := &FlagData{
		name:         name,
		defaultValue: defaultValue,
		Type:         flagType,
		doc:          doc,
		question:     question,
		validation:   validation,
		init:         init,
	}
	var logic FlagLogic
	if flagType == "" {
		flagType = "input"
	}
	switch flagType {
	case "input":
		logic = &FlagLogicInput{}
	case "multiline":
		logic = &FlagLogicMultiLine{}
	case "confirm":
		logic = &FlagLogicConfirm{}
	case "password":
		logic = &FlagLogicPassword{}
	case "select":
		logic = &FlagLogicSelect{}
	case "editor":
		logic = &FlagLogicEditor{}
	case "multiselect":
		logic = &FlagLogicMultiSelect{}
	default:
		return nil, fmt.Errorf("invalid flag type '%s'", flagType)
	}
	flag := &Flag{
		Data:    data,
		Logic:   logic,
		Command: command,
	}
	if _, err := flag.Validator(); err != nil {
		return nil, err
	}
	command.Flags().String(flag.Data.Name(), "", flag.Data.Doc())
	return flag, nil
}

// Answer is the container for prompted values
type Answer struct {
	Value interface{}
}

// WriteAnswer convert any output from a survey prompt to string
func (a *Answer) WriteAnswer(name string, value interface{}) error {
	a.Value = value
	return nil
}

// Ask prompts the user for the flag value.
// The return value is the user input.
func (f *Flag) Ask() string {
	result := Answer{}
	prompt := f.Logic.Prompt(f.Data)
	if exited := survey.AskOne(prompt, &result); exited != nil {
		util.TryThrow(exited)
	}
	return f.Logic.Parse(f.Data, result.Value)
}

// Validator returns the regexp pattern for the flag value
func (f *Flag) Validator() (*regexp.Regexp, error) {
	if f.Data.validation == "" {
		return f.Logic.Validator(f.Data), nil
	}
	return regexp.Compile(f.Data.validation)
}

// Validate validates the flag value against the flag validation pattern.
// If the flag has no validation pattern, it will use the default validator.
// The first return value is the validation result, the second is the error
// returned by the regexp.Compile function.
func (f *Flag) Validate(value string) error {
	pattern, _ := f.Validator()
	if ok := pattern.MatchString(value); !ok {
		return fmt.Errorf("invalid value '%s' for flag '%s'", value, f.Data.Name())
	}
	return nil
}

// Name returns the standardized name of the flag for the cobra command.
func (f *Flag) Name() string {
	return f.Data.Name()
}

// Replacement returns the flag name to be replaced in the template.
func (f *Flag) Replacement() string {
	return f.Data.name
}

// Get executes the flag logic.
// If the flag value exists, it will validate it and return it.
// If the flag value does not exist, it will prompt the user for it
// until it is valid or CTRL+C is pressed.
func (f *Flag) Get() (string, error) {
	value, _ := f.Command.Flags().GetString(f.Name())
	exist := f.Command.Flags().Lookup(f.Name()).Changed
	if exist {
		return value, f.Validate(value)
	} else {
		value = f.Data.defaultValue
		for f.Validate(value) != nil {
			value = f.Ask()
		}
	}
	return value, nil
}
