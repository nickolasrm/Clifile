package lexer_test

import (
	"testing"

	"github.com/nickolasrm/clifile/internal/makefile/lexer"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLexer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lexer Suite")
}

var _ = Describe("Lexer", func() {
	Context("when a valid code is passed", func() {
		code := func() chan rune {
			ch := make(chan rune)
			go func() {
				code := `# This is a commend
				floatvar=1.77
				stringvar="string"
				intvar=3
				rule_1:
					echo "hello world"
				
				# This is another comment
				Rule-2:
					echo "Test"

			`
				for char := range code {
					ch <- rune(char)
				}
			}()
			return ch
		}

		It("should return a channel of tokens and errors", func() {
			tokens, errors := lexer.Lex(code())
			Expect(tokens).ToNot(BeNil())
			Expect(errors).ToNot(BeNil())
		})

		It("should return tokens", func() {
			tokens, errors := lexer.Lex(code())
			expected_tokens := []lexer.Match{
				{Type: lexer.Comment, Value: "# This is a commend"},
				{Type: lexer.NewLine, Value: "\n"},
				{Type: lexer.Indent, Value: "\t"},
				{Type: lexer.Identifier, Value: "floatvar"},
				{Type: lexer.Operator, Value: "="},
				{Type: lexer.Number, Value: "1.77"},
				{Type: lexer.NewLine, Value: "\n"},
				{Type: lexer.Indent, Value: "\t"},
				{Type: lexer.Identifier, Value: "stringvar"},
				{Type: lexer.Operator, Value: "="},
				{Type: lexer.String, Value: "\"string\""},
				{Type: lexer.NewLine, Value: "\n"},
				{Type: lexer.Indent, Value: "\t"},
				{Type: lexer.Identifier, Value: "intvar"},
				{Type: lexer.Operator, Value: "="},
				{Type: lexer.Number, Value: "3"},
				{Type: lexer.NewLine, Value: "\n"},
				{Type: lexer.Indent, Value: "\t"},
				{Type: lexer.Identifier, Value: "rule_1"},
				{Type: lexer.Operator, Value: ":"},
				{Type: lexer.NewLine, Value: "\n"},
				{Type: lexer.Indent, Value: "\t\t"},
				{Type: lexer.Identifier, Value: "echo"},
				{Type: lexer.String, Value: "\"hello world\""},
				{Type: lexer.NewLine, Value: "\n"},
				{Type: lexer.NewLine, Value: "\n"},
				{Type: lexer.Comment, Value: "# This is another comment"},
				{Type: lexer.NewLine, Value: "\n"},
				{Type: lexer.Indent, Value: "\t"},
				{Type: lexer.Identifier, Value: "Rule-2"},
				{Type: lexer.Operator, Value: ":"},
				{Type: lexer.NewLine, Value: "\n"},
				{Type: lexer.Indent, Value: "\t\t"},
				{Type: lexer.Identifier, Value: "echo"},
				{Type: lexer.String, Value: "\"Test\""},
				{Type: lexer.NewLine, Value: "\n"},
				{Type: lexer.NewLine, Value: "\n"},
			}
			i := 0
			for token := range tokens {
				Expect(token.Type).To(Equal(expected_tokens[i].Type))
				Expect(token.Value).To(Equal(expected_tokens[i].Value))
				i++
			}
			Expect(errors).To(BeClosed())
		})
	})
})
