package lisp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"unicode"
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
		return 0, fmt.Errorf("function %q not found", node.Name)
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

func isFirstSymbolRune(r rune) bool {
	return (r != '(' && r != ')') && (unicode.IsPunct(r) || unicode.IsSymbol(r) || unicode.IsLetter(r))
}

func isSymbolRune(r rune) bool {
	return isFirstSymbolRune(r) || unicode.IsDigit(r)
}

func readWhile(f func(r rune) bool, in *bufio.Reader) (read []byte, err error) {
	read = make([]byte, 0, 1)
	next := make([]byte, 1)
	var b byte
	next, _ = in.Peek(1)
	for len(next) > 0 && f(rune(next[0])) {
		b, err = in.ReadByte()
		if err != nil {
			return
		}
		read = append(read, b)
		next, _ = in.Peek(1)
	}
	return
}

func parse(in *bufio.Reader) (result TreeNode, err error) {
	readWhile(unicode.IsSpace, in)

	next, err := in.Peek(1)
	var read []byte

	switch {
	case next[0] == '(':
		in.ReadByte()
		// read name
		read, err = readWhile(isSymbolRune, in)
		if err != nil {
			return
		}
		name := string(read)
		// while not ')' read spaces + argument
		args := make([]TreeNode, 0)
		next, err = in.Peek(1)
		for len(next) > 0 && next[0] != ')' {
			readWhile(unicode.IsSpace, in)
			var child TreeNode
			child, err = parse(in)
			if err != nil {
				return
			}
			args = append(args, child)
			next, err = in.Peek(1)
			if err != nil {
				return
			}
		}
		_, err = in.ReadByte()
		result = &FunctionNode{Name: name, Args: args}
		return
	case unicode.IsDigit(rune(next[0])):
		// read number
		read, err = readWhile(unicode.IsDigit, in)
		if err != nil {
			return
		}
		var n int
		n, err = strconv.Atoi(string(read))
		result = &ValueNode{n}
		return
	default:
		err = fmt.Errorf("invalid symbol %s", next)
		return
	}
	return
}
