package flag_test

import (
	"regexp"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/nickolasrm/clifile/internal/interpreter/flag"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

func TestFlag(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flag Suite")
}

// TestPrompt is the fake Prompt for the survey ask function
type TestPrompt struct{}

func (t *TestPrompt) Prompt(config *survey.PromptConfig) (interface{}, error) {
	return "test", nil
}

func (t *TestPrompt) Cleanup(cfg *survey.PromptConfig, val interface{}) error {
	return nil
}

func (t *TestPrompt) Error(cfg *survey.PromptConfig, err error) error {
	return nil
}

// TestLogic is the fake logic for running the flag
type TestLogic struct{}

func (*TestLogic) Prompt(data *flag.FlagData) survey.Prompt {
	return &TestPrompt{}
}

func (*TestLogic) Validator(data *flag.FlagData) *regexp.Regexp {
	return regexp.MustCompile(".*")
}

func (*TestLogic) Parse(data *flag.FlagData, value interface{}) string {
	return value.(string)
}

// CreateFlag create a Flag object from a map of arguments
func CreateFlag(args map[string]string) (*flag.Flag, error) {
	name := args["name"]
	if name == "" {
		name = "flag"
	}
	return flag.CreateFlag(
		name,
		args["doc"],
		args["default"],
		args["type"],
		args["question"],
		args["validation"],
		args["init"],
		&cobra.Command{},
	)
}

// SnapshotFlag is captures the current structure of the Flag and its child objects
// alongside the parsed output
func SnapshotFlag(args map[string]string, output interface{}) {
	f, err := CreateFlag(args)
	t := GinkgoT()
	snaps.MatchSnapshot(t, map[string]interface{}{
		"flagData":           f.Data,
		"flagLogic":          f.Logic,
		"flagLogicPrompt":    f.Logic.Prompt(f.Data),
		"flagLogicValidator": f.Logic.Validator(f.Data).String(),
		"flagLogicParse":     f.Logic.Parse(f.Data, output),
		"err":                err,
	})
}

var _ = Describe("Flag", func() {
	When("no value is provided", func() {
		When("default value is provided", func() {
			It("should return the default", func() {
				f, _ := CreateFlag(map[string]string{
					"type":    "input",
					"default": "test",
				})
				Expect(f.Get()).To(Equal("test"))
			})
		})
		When("no default value is provided", func() {
			It("should prompt the user", func() {
				f, _ := CreateFlag(map[string]string{"validation": ".+"})
				f.Logic = &TestLogic{}
				o, _ := f.Get()
				Expect(o).To(Equal("test"))
			})
		})
	})
	When("value is provided", func() {
		It("should validate it", func() {
			f, _ := CreateFlag(map[string]string{})
			f.Command.Flags().Lookup(f.Name()).Value.Set("test")
			f.Command.Flags().Lookup(f.Name()).Changed = true
			Expect(f.Get()).To(Equal("test"))
		})
	})
	When("calling doc", func() {
		It("should return the doc", func() {
			f, _ := CreateFlag(map[string]string{
				"doc": "Doc",
			})
			Expect(f.Data.Doc()).To(Equal("Doc"))
		})
		When("default value is provided", func() {
			f, _ := CreateFlag(map[string]string{
				"doc":     "Doc",
				"default": "def",
			})
			Expect(f.Data.Doc()).To(Equal("Doc (default: def)"))
		})
	})
	When("calling question", func() {
		When("question is not provided", func() {
			It("should return the flag name", func() {
				f, _ := CreateFlag(map[string]string{})
				Expect(f.Data.Question()).To(Equal("flag"))
			})
		})
		When("question is provided", func() {
			It("should return the question", func() {
				f, _ := CreateFlag(map[string]string{
					"question": "test",
				})
				Expect(f.Data.Question()).To(Equal("test"))
			})
		})
	})
	When("type is defined", func() {
		When("type is input", func() {
			It("should return a flag with an input survey", func() {
				SnapshotFlag(map[string]string{
					"type": "input",
				}, "test")
			})
		})
		When("type is multiline", func() {
			It("should return a flag with a multiline survey", func() {
				SnapshotFlag(map[string]string{
					"type": "multiline",
				}, "test")
			})
		})
		When("type is confirm", func() {
			It("should return a confirm survey", func() {
				SnapshotFlag(map[string]string{
					"type": "confirm",
				}, false)
			})
			When("output is false", func() {
				It("should parse to 'n'", func() {
					f, _ := CreateFlag(map[string]string{"type": "confirm"})
					Expect(f.Logic.Parse(f.Data, false)).To(Equal("n"))
				})
			})
			When("output is true", func() {
				It("should return a confirm survey and parse to 'y'", func() {
					f, _ := CreateFlag(map[string]string{"type": "confirm"})
					Expect(f.Logic.Parse(f.Data, true)).To(Equal("y"))
				})
			})
		})
		When("type is password", func() {
			It("should return a flag with a password survey", func() {
				SnapshotFlag(map[string]string{
					"type": "password",
				}, "test")
			})
		})
		When("type is select", func() {
			It("should return a flag with a select survey", func() {
				SnapshotFlag(map[string]string{
					"type": "select",
				}, survey.OptionAnswer{Value: "test", Index: 0})
			})
		})
		When("type is editor", func() {
			It("should return a flag with an editor survey", func() {
				SnapshotFlag(map[string]string{
					"type": "editor",
				}, "test")
			})
		})
		When("type is multiselect", func() {
			It("should return a flag with a multiselect survey", func() {
				SnapshotFlag(map[string]string{
					"type": "multiselect",
				}, []survey.OptionAnswer{
					{Index: 0, Value: "a"},
					{Index: 1, Value: "b"},
				})
			})
		})
		When("type is unknown", func() {
			It("should return an error", func() {
				_, err := CreateFlag(map[string]string{
					"type": "unknown",
				})
				Expect(err).ToNot(BeNil())
			})
		})
	})
	When("validator doesn't compile", func() {
		It("should error", func() {
			_, err := CreateFlag(map[string]string{
				"validation": "?<!error",
			})
			Expect(err).ToNot(BeNil())
		})
	})
	When("name is not standardized", func() {
		It("name should be standardized", func() {
			flag, _ := CreateFlag(map[string]string{"name": "NoN_STANDARD"})
			Expect(flag.Name()).To(Equal("non-standard"))
		})
		It("replacement should stay the same", func() {
			flag, _ := CreateFlag(map[string]string{"name": "NoN_Standard"})
			Expect(flag.Replacement()).To(Equal("NoN_Standard"))
		})
	})
})
