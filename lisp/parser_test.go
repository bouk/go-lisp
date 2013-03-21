package lisp

import (
	"strings"
	"testing"
)

var errorTestcases = []string{
	`"`,
	"(+ 1 1",
}

func TestParserError(t *testing.T) {
	for i, testcase := range errorTestcases {
		_, err := Parse(strings.NewReader(testcase))
		if err == nil {
			t.Errorf("#%d: %#v should have failed", i, testcase)
		}
	}
}
