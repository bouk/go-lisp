package lisp

import (
	"strings"
	"testing"
)

var functionTestcases = []struct {
	input  string
	output int
}{
	{"(+ 1 2)", 3},
	{"(- 1 2)", -1},
	{"(* 2 3)", 6},
	{"(/ 4 2)", 2},
	{"(+ (+ 1  1) (+ 1 1))", 4},
	{"1", 1},
}

func TestFunctions(t *testing.T) {
	for i, testcase := range functionTestcases {
		node, parseErr := Parse(strings.NewReader(testcase.input))
		if parseErr == nil {
			output, interpretErr := node.Interpret()
			if interpretErr == nil {
				if output != testcase.output {
					t.Errorf("%v %v => %v wanted %v", i, testcase.input, output, testcase.output)
				}
			} else {
				t.Errorf("%v Error %v", i, interpretErr)
			}
		} else {
			t.Errorf("%v Error %v", i, parseErr)
		}
	}
}
