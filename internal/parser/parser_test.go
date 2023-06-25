package parser_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/nickolasrm/clifile/internal/lexer"
	"github.com/nickolasrm/clifile/internal/parser"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parser Suite")
}

// SnapshotProgram is a helper function to test the code
// It receives a code and tries to tokenize it and parse the tokens.
// If an lexical error is thrown, it will fail the test.
// If a semantic error is thrown, it will fail the test.
func SnapshotProgram(code string) {
	t := GinkgoT()
	tokens, error := lexer.Lex(code)
	Expect(error).To(BeNil())
	program, error := parser.Parse(tokens)
	Expect(error).To(BeNil())
	snaps.MatchSnapshot(t, program)
}

// SnapshotProgramError is a helper function to test the code
// It receives a code and tries to tokenize it and parse the tokens.
// If an lexical error is thrown, it will fail the test.
// If a semantic error is not thrown, it will succeed the test.
func SnapshotProgramError(code string) {
	t := GinkgoT()
	tokens, error := lexer.Lex(code)
	Expect(error).To(BeNil())
	_, error = parser.Parse(tokens)
	Expect(error).NotTo(BeNil())
	snaps.MatchSnapshot(t, error)
}

var _ = Describe("Parser", func() {
	When("a code is passed", func() {
		When("first line is a docstring", func() {
			It("should parse it as doc", func() {
				SnapshotProgram(`## Name`)
			})
			When("contains a rule", func() {
				It("should not parse it as rule docstring", func() {
					SnapshotProgram(`## Name
## Descr

## Rule
rule:
	echo "test"
`)
				})
			})
		})
		When("contains a comment", func() {
			It("should ignore a comment", func() {
				SnapshotProgram(`
# This is a comment
`)
			})
		})
		When("contains a variable", func() {
			It("should return a variable", func() {
				SnapshotProgram(`
Var=1
`)
			})
		})
		When("contains a function call", func() {
			It("should return a function call", func() {
				SnapshotProgram(`
Rec=${Fn}
`)
			})
			When("contains arguments", func() {
				It("should return a function call with arguments", func() {
					SnapshotProgram(`
Rec=${Fn arg="1" arg2="2"}
`)
				})
				When("contains an invalid token", func() {
					It("should return an error", func() {
						SnapshotProgramError(`
Rec=${Fn arg="1" arg2="2" arg3}
`)
					})
				})
				When("contains a comment", func() {
					It("should ignore the comment", func() {
						SnapshotProgram(`
Rec=${Fn arg="1" arg2="2" # arg3}
`)
					})
				})
			})
		})
		When("doesn't contain a rule", func() {
			When("contains an action", func() {
				It("should return an error", func() {
					SnapshotProgramError(`
echo "Test"
`)
				})
			})
		})
		When("contains a rule", func() {
			When("contains positional parameters", func() {
				It("should return rule with positional params", func() {
					SnapshotProgram(`
rule: param1 param2
	echo "test"
`)
				})
			})
			When("contains an action", func() {
				It("should return a rule with an action", func() {
					SnapshotProgram(`
rule:
	echo "Test"
`)
				})
				When("contains an indented variable", func() {
					It("should return a rule with an action", func() {
						SnapshotProgram(`
rule:
	Var=1
	echo "Test"
`)
					})
				})
				When("contains an indented function call", func() {
					It("should return a rule with an action", func() {
						SnapshotProgram(`
rule:
	Rec=${Fn}
	echo "Test"
`)
					})
				})
				When("contains multiple actions", func() {
					It("should return a rule with multiple actions", func() {
						SnapshotProgram(`
rule:
	echo "Test"
	echo "Test2"
`)
					})
				})
				When("preceded by a docstring", func() {
					It("should return a rule with an action and a docstring", func() {
						SnapshotProgram(`
## This is a docstring
rule:
	echo "Test"
`)
					})
					When("more than one docstring", func() {
						It("should concatenate the docstrings into one", func() {
							SnapshotProgram(`
## This is a docstring
## This is another docstring
rule:
	echo "Test"
`)
						})
					})
					When("docstring is overly indented", func() {
						It("should return an error", func() {
							SnapshotProgramError(`
	## This is a docstring
rule:
	echo "Test"
`)
						})
					})
				})
				When("contains multiple rules", func() {
					It("should return multiple rules", func() {
						SnapshotProgram(`
rule:
	echo "Test"

rule2:
	echo "Test2"
`)
					})
					When("a rule is nested", func() {
						It("should return a rule with a nested rule", func() {
							SnapshotProgram(`
rule:
	rule2:
		echo "Test2"
`)
						})
						When("another rule is nested and a group appears after", func() {
							It("should return a rule with a nested rule and a group", func() {
								SnapshotProgram(`
rule:
	rule2:
		echo "Test2"
	rule3:
		echo "Test3"
rule4:
	echo "Test4"
`)
							})
						})
						When("deeply nested", func() {
							It("should return a rule with a nested rule and a group", func() {
								SnapshotProgram(`
rule:
	rule2:
		rule3:
			rule4:
				echo "Test4"
`)
							})
						})
						When("more than two tabs are used", func() {
							It("should return an error", func() {
								SnapshotProgramError(`
rule:
	rule2:
			rule3:
				echo "Test4"
`)
							})
						})
						When("an action appears in the group rule", func() {
							When("the action appears before the nested rule", func() {
								It("should return an error", func() {
									SnapshotProgramError(`
rule:
	echo "Test"
	rule2:
		echo "Test2"
`)
								})
							})
							When("the action appears after the nested rule", func() {
								It("should return an error", func() {
									SnapshotProgramError(`
rule:
	rule2:
		echo "Test2"
	echo "Test"
`)
								})
							})
						})
					})
				})
			})
		})
		When("a full example appeared", func() {
			It("should parse multiple matches", func() {
				SnapshotProgram(`## Name
## Descr
## iption

VAR=abd
VAR2="def
dfg"
FLAG1=${fn1 x=y
	z="d"
}
FLAG2=${fn2}

# This is a comment
grp1: FLAG2 FLAG1
	cmd2:
		VARSHELL=2
		echo "test"
	## Documentation
	cmd3:
		sed ${FLAG1}
## Doc
grp2:
	cmd1:
		hey ${VAR1}
`)
			})
		})
	})
})
