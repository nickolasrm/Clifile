package variable_test

import (
	"testing"

	"github.com/nickolasrm/clifile/internal/models/variable"
	. "github.com/nickolasrm/clifile/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVariable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Variable Suite")
}

var _ = Describe("Variable", func() {
	When("a new variable is created", func() {
		var v *variable.Variable

		BeforeEach(func() {
			v = variable.NewVariable("name", "value")
		})

		It("should return a variable type", func() {
			Expect(v).ToNot(BeNil())
		})
		When("Name() is called", func() {
			It("should return the variable name", func() {
				MatchSnapshot(v.Name())
			})
		})
		When("Value() is called", func() {
			It("should return its value", func() {
				MatchSnapshot(v.Value())
			})
		})
		When("SetValue() is called", func() {
			It("should update the variable value", func() {
				v.SetValue("custom")
				MatchSnapshot(v.Value())
			})
		})
	})
})
