package lisp

import (
	"bytes"
	"strings"
	"testing"
)

var functionTestcases = []struct {
	input  string
	output Value
}{
	{"(+ 1 2)", 3},
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
	{`
	(set a 1)
	(defun lol a)
	(set a 2)
	a
		`, 2},
	{`(if 1 1 0)`, 1},
	{`(if 1 0 1)`, 0},
	{`(if (== 1 1) 1 0)`, 1},
	{`(if (== 1 2) 0 1)`, 1},
}

func TestFunctions(t *testing.T) {
	for i, testcase := range functionTestcases {
		node := NewRootNode()
		parseErr := node.Parse(strings.NewReader(testcase.input))
		if parseErr == nil {
			output, interpretErr := node.Interpret(nil, nil)
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

var outputTestcases = []struct {
	input  string
	output string
}{
	{`(print "hello")`, `hello`},
	{`
		(defun p2 a (print a))
		(p2 "a")`, "a"},
	{`(print 123)`, "123"},
	{`
		(defun println line (print line "\n"))
		(println "yoloswag")`, "yoloswag\n"},
}

func TestOutputfunction(t *testing.T) {
	for i, testcase := range outputTestcases {
		node := NewRootNode()
		parseErr := node.Parse(strings.NewReader(testcase.input))
		if parseErr == nil {
			var buf bytes.Buffer
			_, interpretErr := node.Interpret(&buf, nil)
			if interpretErr == nil {
				if buf.String() != testcase.output {
					t.Errorf("#%d: %#v => %#v wanted %#v", i, testcase.input, buf.String(), testcase.output)
				}
			} else {
				t.Errorf("#%v: error %v", i, interpretErr)
			}
		} else {
			t.Errorf("#%v: error %s in %#v", i, parseErr, testcase.input)
		}
	}
}
