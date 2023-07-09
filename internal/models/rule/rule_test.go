package rule_test

import (
	"testing"

	"github.com/nickolasrm/clifile/internal/models/rule"
	. "github.com/nickolasrm/clifile/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rule Suite")
}

var _ = Describe("Rule", func() {
	When("a new rule is created", func() {
		var r *rule.Rule

		BeforeEach(func() {
			r = rule.NewRule("name", []string{"arg1", "arg2", "arg3"}, "doc", "action")
		})

		It("should return a variable type", func() {
			Expect(r).ToNot(BeNil())
		})
		When("Name() is called", func() {
			It("should return the rule name", func() {
				MatchSnapshot(r.Name())
			})
		})
		When("Positional() is called", func() {
			It("it should return the positional args", func() {
				MatchSnapshot(r.Positional())
			})
		})
		When("Doc() is called", func() {
			It("should return the rule doc", func() {
				MatchSnapshot(r.Doc())
			})
		})
		When("Actions() is called", func() {
			It("should return the rule actions", func() {
				MatchSnapshot(r.Actions())
			})
		})
		When("AppendActions() is called", func() {
			It("should append the text in actions as a newline", func() {
				r.AppendActions("newline\nnewline2")
				MatchSnapshot(r.Actions())
			})
		})
		When("AddRule() is called", func() {
			It("should add a new child rule", func() {
				r.AddRule(r)
				MatchSnapshot(r.Rules())
			})
		})
		When("Rules() is called", func() {
			When("no child rules exist", func() {
				It("should return an empty array", func() {
					MatchSnapshot(r.Rules())
				})
			})
			When("child rules exist", func() {
				It("should return the array of child rules", func() {
					r.AddRule(r)
					MatchSnapshot(r)
				})
			})
		})
		When("Rule() is called", func() {
			When("an existent child rule name is passed", func() {
				It("should return the rule", func() {
					r.AddRule(r)
					MatchSnapshot(r.Rule("name"))
				})
			})
			When("a non-existent child rule name is passed", func() {
				It("should return nil", func() {
					Expect(r.Rule("asd")).To(BeNil())
				})
			})
		})
	})
})
