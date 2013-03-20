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
		output, err := Parse(strings.NewReader(testcase.input)).Interpret()
		if output != testcase.output {
			t.Errorf("%v %v => %v wanted %v", i, testcase.input, output, testcase.output)
		}
		if err != nil {
			t.Errorf("%v Error %v", i, err)
		}
	}
}
