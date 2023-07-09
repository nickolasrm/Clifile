package call_test

import (
	"testing"

	"github.com/nickolasrm/clifile/internal/models/call"
	. "github.com/nickolasrm/clifile/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCall(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Call Suite")
}

var _ = Describe("Call", func() {
	When("a newcall is created", func() {
		var c *call.Call

		BeforeEach(func() {
			c = call.NewCall("name", "func", map[string]string{
				"arg1": "val",
				"arg2": "val2",
			})
		})

		It("should return a call type", func() {
			Expect(c).ToNot(BeNil())
		})
		When("Function() is called", func() {
			It("should return the function name", func() {
				MatchSnapshot(c.Name())
			})
		})
		When("Name() is called", func() {
			It("should return the call name", func() {
				MatchSnapshot(c.Function())
			})
		})
		When("Arguments() is called", func() {
			It("should return the arguments map", func() {
				MatchSnapshot(c.Arguments())
			})
		})
		When("Argument() is called", func() {
			When("an existing argument is passed", func() {
				It("should return its value", func() {
					MatchSnapshot(c.Argument("arg1"))
				})
			})
			When("an non-existing argument is passed", func() {
				It("should be empty", func() {
					Expect(c.Argument("random")).To(BeEmpty())
				})
			})
		})
	})
})
