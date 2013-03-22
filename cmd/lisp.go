package main

import (
	"fmt"
	"github.com/boukevanderbijl/go-lisp/lisp"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Please give at least one input file")
		os.Exit(0)
	}

	root := lisp.NewRootNode()
	var err error

	for _, name := range os.Args[1:] {
		var input io.Reader

		if name == "-" {
			input = os.Stdin
		} else {
			file, err := os.Open(name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(0)
			} else {
				input = file
			}
		}

		err = root.Parse(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "PARSING ERROR: ", err)
			os.Exit(0)
		}
	}

	_, err = root.Interpret(os.Stdout, os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: ", err)
	}
}
