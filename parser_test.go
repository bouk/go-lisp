package lisp

import (
	"strings"
	"testing"
)

var parserTestcases = []struct {
	input  string
	output TreeNode
}{
	{
		"(+ 1 2)",
		&FunctionNode{
			Name: "+",
			Args: []TreeNode{
				&ValueNode{1},
				&ValueNode{2},
			},
		},
	},
	{
		"(+ (- 2 2) (* 3 4))",
		&FunctionNode{
			Name: "+",
			Args: []TreeNode{
				&FunctionNode{
					Name: "-",
					Args: []TreeNode{
						&ValueNode{2},
						&ValueNode{2},
					},
				},
				&FunctionNode{
					Name: "*",
					Args: []TreeNode{
						&ValueNode{3},
						&ValueNode{4},
					},
				},
			},
		},
	},
}

func TestParser(t *testing.T) {
	for i, testcase := range parserTestcases {
		output := Parse(strings.NewReader(testcase.input))
		if output != testcase.output {
			t.Errorf("%v %v => %#v wanted %#v", i, testcase.input, output, testcase.output)
		}
	}
}
