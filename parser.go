package lisp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type TreeNode interface {
	Interpret() (n int, err error)
}

type ValueNode struct {
	Value int
}

func (node *ValueNode) Interpret() (n int, err error) {
	return node.Value, nil
}

type FunctionNode struct {
	Name string
	Args []TreeNode
}

func (node *FunctionNode) Interpret() (n int, err error) {
	f, exists := registeredFunctions[node.Name]
	if !exists {
		return 0, fmt.Errorf("Function %q not found", node.Name)
	}

	args := make([]int, len(node.Args))
	for index, child := range node.Args {
		value, err := child.Interpret()
		if err != nil {
			return 0, err
		}
		args[index] = value
	}

	return f(args)
}

func Parse(input io.Reader) (result TreeNode, err error) {
	return parse(bufio.NewReader(input))
}

func isNumber(b byte) bool {
	return b >= '0' && b <= '9'
}

func parse(in *bufio.Reader) (result TreeNode, err error) {
	next, err := in.Peek(1)
	switch {
	case next[0] == '(':
		// read name

		// while not ')' read space + argument
		err = fmt.Errorf("not implemented")
		return
	case isNumber(next[0]):
		// read number
		read := make([]byte, 0, 1)

		for len(next) > 0 && isNumber(next[0]) {
			var b byte
			b, err = in.ReadByte()
			if err != nil {
				return
			}
			read = append(read, b)
			next, _ = in.Peek(1)
		}

		n, err := strconv.Atoi(string(read))
		return &ValueNode{n}, err
	default:
		err = fmt.Errorf("invalid symbol")
		return
	}
	return
}
