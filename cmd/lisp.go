package main

import (
	"fmt"
	"lisp"
	"os"
)

func main() {
	root, err := lisp.Parse(os.Stdin)
	if err != nil {
		fmt.Println(err)
	}
	root.Interpret(os.Stdout)
}
