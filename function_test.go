package lisp

import (
	"strings"
	"testing"
)

var functionTestcases = []struct {
	input  string
	output Value
}{
	{"(+ 1 2)", 3},
	{"(- 1 2)", -1},
	{"(* 2 3)", 6},
	{"(/ 4 2)", 2},
	{"(+ (+ 1  1) (+ 1 1))", 4},
	{"(+ 1 2)", 3},
	{`
(+ 1


	(*
		1
		2)


	)
`, 3},
	{"1", 1},
	{"yoloswag", nil},
	{"-1", -1},
	{`(+ "a" "bc")`, "abc"},
	{`(+ "#" (+ "yolo" (* (* 20 7) 3)))`, "#yolo420"},
	{`"\""`, `"`},
	{`"\\"`, `\`},
	{`(scope
		(set a 1)
		a
		)`, 1},
	{`(scope
		(set lol 1)
		(set hello 4)
		(+ lol hello)
		)`, 5},
	{``, nil},
	{`(scope
		(set _return 4)
		(set hello 1337))`, 4},
	{`(set a 1)
		a`, 1},
}

func TestFunctions(t *testing.T) {
	for i, testcase := range functionTestcases {
		node, parseErr := Parse(strings.NewReader(testcase.input))
		if parseErr == nil {
			output, interpretErr := node.Interpret(nil)
			if interpretErr == nil {
				if output != testcase.output {
					t.Errorf("#%d: %#v => %#v wanted %#v", i, testcase.input, output, testcase.output)
				}
			} else {
				t.Errorf("#%v: error %v", i, interpretErr)
			}
		} else {
			t.Errorf("#%v: error %s in %#v", i, parseErr, testcase.input)
		}
	}
}
