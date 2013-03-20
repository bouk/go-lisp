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
}

func TestFunctions(t *testing.T) {
	for i, testcase := range functionTestcases {
		output := Parse(strings.NewReader(testcase.input)).Interpret()
		if output != testcase.output {
			t.Errorf("%d %q => %q wanted %q", i, testcase.input, output, testcase.output)
		}
	}
}
