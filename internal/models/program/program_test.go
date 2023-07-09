package program_test

import (
	"testing"

	"github.com/nickolasrm/clifile/internal/models/call"
	"github.com/nickolasrm/clifile/internal/models/program"
	"github.com/nickolasrm/clifile/internal/models/rule"
	"github.com/nickolasrm/clifile/internal/models/variable"
	. "github.com/nickolasrm/clifile/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProgram(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Program Suite")
}

var _ = Describe("Program", func() {
	When("a new program is created", func() {
		var p *program.Program

		BeforeEach(func() {
			p = program.NewProgram()
		})

		It("should return a program type", func() {
			MatchSnapshot(p)
		})

		When("Doc() is called", func() {
			It("should return the program doc", func() {
				MatchSnapshot(p.Doc())
			})
		})
		When("SetDoc() is called", func() {
			It("should replace the program doc", func() {
				p.SetDoc("new doc")
				MatchSnapshot(p)
			})
		})
		When("Call() is called", func() {
			When("the passed call exists", func() {
				It("should return the call", func() {
					c := call.NewCall("call", "fn", map[string]string{})
					p.AddCall(c)
					MatchSnapshot(p.Call("call"))
				})
			})
			When("the passed call does not exist", func() {
				It("should return nil", func() {
					Expect(p.Call("call")).To(BeNil())
				})
			})
		})
		When("AddCall() is called", func() {
			It("should add the call to the program", func() {
				c := call.NewCall("call", "fn", map[string]string{})
				p.AddCall(c)
				MatchSnapshot(p)
			})
		})
		When("Rules() is called", func() {
			When("the program has no rules", func() {
				It("should return an empty slice", func() {
					Expect(p.Rules()).To(BeEmpty())
				})
			})
			When("the program has rules", func() {
				It("should return the rules", func() {
					r := rule.NewRule("rule", []string{"pos1"}, "doc", "action")
					p.AddRule(r)
					MatchSnapshot(p.Rules())
				})
			})
		})
		When("AddRule() is called", func() {
			It("should add the rule to the program", func() {
				r := rule.NewRule("rule", []string{"pos1"}, "doc", "action")
				p.AddRule(r)
				MatchSnapshot(p)
			})
		})
		When("Rule() is called", func() {
			When("the passed rule exists", func() {
				It("should return the rule", func() {
					r := rule.NewRule("rule", []string{"pos1"}, "doc", "action")
					p.AddRule(r)
					MatchSnapshot(p.Rule("rule"))
				})
			})
			When("the passed rule does not exist", func() {
				It("should return nil", func() {
					Expect(p.Rule("rule")).To(BeNil())
				})
			})
		})
		When("AddVariable() is called", func() {
			It("should add the variable to the program", func() {
				v := variable.NewVariable("var", "value")
				p.AddVariable(v)
				MatchSnapshot(p)
			})
		})
		When("Variable() is called", func() {
			When("the passed variable exists", func() {
				It("should return the variable", func() {
					v := variable.NewVariable("var", "value")
					p.AddVariable(v)
					MatchSnapshot(p.Variable("var"))
				})
			})
			When("the passed variable does not exist", func() {
				It("should return nil", func() {
					Expect(p.Variable("var")).To(BeNil())
				})
			})
		})
		When("SetVariable() is called", func() {
			When("the passed variable exists", func() {
				It("should set the variable value", func() {
					v := variable.NewVariable("var", "value")
					p.AddVariable(v)
					p.SetVariable("var", "new value")
					MatchSnapshot(p)
				})
			})
			When("the passed variable does not exist", func() {
				It("should create the variable", func() {
					p.SetVariable("var", "value")
					MatchSnapshot(p)
				})
			})
		})
	})
})
