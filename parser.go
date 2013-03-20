package lisp

import (
	"fmt"
	"io"
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

func Parse(input io.Reader) TreeNode {
	return &ValueNode{0}
}
