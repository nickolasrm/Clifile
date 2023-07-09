package lexer_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/nickolasrm/clifile/internal/lexer"
	. "github.com/nickolasrm/clifile/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLexer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lexer Suite")
}

// SnapshotTokens is a helper function to test the code.
// It receives a code and tries to match to the previously saved snapshots.
// If it does not match, it will fail the test.
// If an error occurs, it will fail the test.
func SnapshotTokens(code string) {
	t := GinkgoT()
	tokens, error := lexer.Lex(code)
	Expect(error).To(BeNil())
	snaps.MatchSnapshot(t, tokens)
}

var _ = Describe("Lexer", func() {
	Describe("Lex", func() {
		When("a comment is passed", func() {
			It("should return a comment token", func() {
				SnapshotTokens("# This is a comment")
			})
			When("# is present", func() {
				It("should return a single comment token", func() {
					SnapshotTokens("# This is a #comment")
				})
			})
			When("another comment is passed", func() {
				It("should return two comment tokens", func() {
					SnapshotTokens("# C1\n# C2")
				})
			})
		})
		When("a docstring is passed", func() {
			It("should return a docstring token", func() {
				SnapshotTokens("## This is a docstring")
			})
			When("## is present", func() {
				It("should return a single docstring token", func() {
					SnapshotTokens("## This is a ##docstring")
				})
			})
			When("another docstring is passed", func() {
				It("should return two docstring tokens", func() {
					SnapshotTokens("## D1\n## D2")
				})
			})
		})
		When("a variable is passed", func() {
			When("unquoted value is passed", func() {
				It("should return a variable token", func() {
					SnapshotTokens("VAR=val")
				})
				When("multiple lines are passed", func() {
					It("should return only the variable line", func() {
						SnapshotTokens("VAR=val\n#comment")
					})
				})
			})
			When("quoted value is passed", func() {
				It("should return a variable token", func() {
					SnapshotTokens("VAR=\"val\"")
				})
				When("multiple lines are between quotes", func() {
					It("should return a variable token", func() {
						SnapshotTokens("VAR=\"val\nval\"")
					})
				})
				When("space or tabs appear around equal", func() {
					It("should return a variable token", func() {
						SnapshotTokens("VAR \t=  \t\"ASD\"")
					})
				})
			})
		})
		When("a call is passed", func() {
			It("should return a call token", func() {
				SnapshotTokens("VAR=${func}")
			})
			When("multiple lines are between curly braces", func() {
				It("should return a call token", func() {
					SnapshotTokens("VAR=${func\nfunc}")
				})
			})
			When("space or tabs appear around equal", func() {
				It("should return a call token", func() {
					SnapshotTokens("VAR\t =  \t${func}")
				})
			})
		})
		When("a rule is passed", func() {
			It("should return a rule token", func() {
				SnapshotTokens("rule:")
			})
			When("arguments are passed", func() {
				It("should return a rule token", func() {
					SnapshotTokens("rule: arg1 arg2")
				})
			})
			When("another rule is passed", func() {
				It("should return two rule tokens", func() {
					SnapshotTokens("rule1: a b\nrule2: c d")
				})
				When("the second rule is nested", func() {
					It("should return two rule tokens", func() {
						SnapshotTokens("rule1: a\n\trule2:")
					})
				})
			})
		})
		When("an action is passed", func() {
			It("should return an action token", func() {
				SnapshotTokens("echo \"Hello World\"")
			})
			When("multiple lines are passed", func() {
				It("should return multiple action tokens", func() {
					SnapshotTokens("echo \"Hello World\"\necho \"Hello World\"")
				})
			})
		})
		When("a code snippet is passed", func() {
			It("should return a sequence of tokens", func() {
				SnapshotTokens(`
# This is a comment
VAR=val
VAR2="val"
VAR3=${func}
VAR4=${func
func}
rule1: a b
	rule4:
		echo "Hello World"
rule2: c d
## This is a docstring
rule3:
	echo "Hello World"
	echo "Hello World"
`)
			})
		})
	})

	Describe("Match", func() {
		When("a new match is created", func() {
			var m *lexer.Match
			BeforeEach(func() {
				m = lexer.NewMatch(lexer.Action, []string{"echo Hello World"})
			})

			When("Type() is called", func() {
				It("should return the type of the match", func() {
					MatchSnapshot(m.Type())
				})
			})

			When("Value() is called", func() {
				When("the index exists", func() {
					It("should return the first value of the match", func() {
						MatchSnapshot(m.Value(0))
					})
				})
				When("the index is out of bounds", func() {
					It("should return an empty string", func() {
						Expect(func() { m.Value(1) }).To(Panic())
					})
				})
			})
		})
	})
})
