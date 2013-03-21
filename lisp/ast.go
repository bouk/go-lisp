package lisp

import (
	"fmt"
	"io"
	"os"
)

type TreeNode interface {
	Interpret(s *Scope) (v Value, err error)
}

type ValueNode struct {
	Value
}

func (node *ValueNode) Interpret(s *Scope) (v Value, err error) {
	return node.Value, nil
}

type FunctionNode struct {
	Name string
	Args []TreeNode
}

func (node *FunctionNode) Interpret(s *Scope) (v Value, err error) {
	f := s.FindFunction(node.Name)
	if f == nil {
		return nil, fmt.Errorf("function %q not found", node.Name)
	}

	return f(s, node.Args)
}

type SymbolNode struct {
	Name string
}

func (node *SymbolNode) Interpret(s *Scope) (v Value, err error) {
	return s.GetVariable(node.Name), nil
}

type RootNode struct {
	Program []TreeNode
}

func (node *RootNode) Interpret(out io.Writer) (v Value, err error) {
	if out == nil {
		out = os.Stdout
	}

	s := NewScope(nil)
	builtinFunctions(s)
	s.Out = out

	for _, node := range node.Program {
		v, err = node.Interpret(s)
		if err != nil {
			v = nil
			return
		}
	}
	return
}
