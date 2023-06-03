package lexer_test

import (
	"testing"

	"github.com/nickolasrm/clifile/internal/makefile/lexer"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

func TestLexer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lexer Suite")
}

// ToChan is a helper function to convert a string to a channel
// It returns a channel of runes
func ToChan(str string) chan rune {
	ch := make(chan rune)
	go func() {
		for _, char := range str {
			ch <- rune(char)
		}
		close(ch)
	}()
	return ch
}

// MatchTokens is a helper function to test channels
// It receives a struct and checks if it matches the fields passed
// If it does not match, it will fail the test
func MatchTokens(code string, fields ...Fields) {
	ch := ToChan(code)
	s, _ := lexer.Lex(ch)
	Expect(s).ToNot(BeNil())
	for _, field := range fields {
		Eventually(s).Should(Receive(PointTo(MatchAllFields(field))))
	}
}

var _ = Describe("Lexer", func() {
	When("a comment is passed", func() {
		When("one line", func() {
			It("should return a comment token", func() {
				MatchTokens("# This is a comment", Fields{
					"Type":  Equal(lexer.Comment),
					"Value": Equal("# This is a comment"),
				})
			})
		})
		When("# is present", func() {
			It("should return a single comment token", func() {
				MatchTokens("# This is a #comment", Fields{
					"Type":  Equal(lexer.Comment),
					"Value": Equal("# This is a #comment"),
				})
			})
		})
		When("another comment is passed", func() {
			It("should return two comment tokens", func() {
				MatchTokens("# C1\n# C2", Fields{
					"Type":  Equal(lexer.Comment),
					"Value": Equal("# C1"),
				}, Fields{
					"Type":  Equal(lexer.NewLine),
					"Value": Equal("\n"),
				}, Fields{
					"Type":  Equal(lexer.Comment),
					"Value": Equal("# C2"),
				})
			})
		})
	})
	When("a value is passed", func() {
		When("one line", func() {
			It("should return a value token", func() {
				MatchTokens("\"This is a value\"", Fields{
					"Type":  Equal(lexer.Value),
					"Value": Equal("\"This is a value\""),
				})
			})
		})
		When("multiline", func() {
			It("should return a value token", func() {
				MatchTokens("\"This is a\nvalue\"", Fields{
					"Type":  Equal(lexer.Value),
					"Value": Equal("\"This is a\nvalue\""),
				})
			})
		})
		When("another value is passed", func() {
			It("should return two value tokens", func() {
				MatchTokens("\"S1\"\n\"S2\"", Fields{
					"Type":  Equal(lexer.Value),
					"Value": Equal("\"S1\""),
				}, Fields{
					"Type":  Equal(lexer.NewLine),
					"Value": Equal("\n"),
				}, Fields{
					"Type":  Equal(lexer.Value),
					"Value": Equal("\"S2\""),
				})
			})
		})
		When("a value is not closed", func() {
			It("should return a value token", func() {
				MatchTokens("\"This is a value", Fields{
					"Type":  Equal(lexer.Value),
					"Value": Equal("\"This is a value"),
				})
			})
		})
	})
	When("unquoted value is passed", func() {
		It("should return an unquoted value token", func() {
			MatchTokens("This is a value", Fields{
				"Type":  Equal(lexer.UnquotedValue),
				"Value": Equal("This is a value"),
			})
		})
		When("\n is present", func() {
			It("should return an unquoted value token", func() {
				MatchTokens("This is a value\n", Fields{
					"Type":  Equal(lexer.UnquotedValue),
					"Value": Equal("This is a value"),
				}, Fields{
					"Type":  Equal(lexer.NewLine),
					"Value": Equal("\n"),
				})
			})
		})
	})
	When("an action is passed", func() {
		It("should return an action token", func() {
			MatchTokens("\tThis is an action", Fields{
				"Type":  Equal(lexer.Action),
				"Value": Equal("\tThis is an action"),
			})
		})
		When("\\n is present", func() {
			It("should return an action token", func() {
				MatchTokens("\tThis is an action\n", Fields{
					"Type":  Equal(lexer.Action),
					"Value": Equal("\tThis is an action"),
				}, Fields{
					"Type":  Equal(lexer.NewLine),
					"Value": Equal("\n"),
				})
			})
			When("another action is passed", func() {
				It("should return two action tokens", func() {
					MatchTokens("\tA1\n\tA2", Fields{
						"Type":  Equal(lexer.Action),
						"Value": Equal("\tA1"),
					}, Fields{
						"Type":  Equal(lexer.NewLine),
						"Value": Equal("\n"),
					}, Fields{
						"Type":  Equal(lexer.Action),
						"Value": Equal("\tA2"),
					})
				})
			})
		})
	})
})
