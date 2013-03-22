package lisp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

var (
	doneReading    = errors.New("done reading input")
	escapedLetters = map[byte]byte{
		'n': '\n',
		't': '\t',
	}
)

// Parses an io stream like a file
func Parse(input io.Reader) (upperNode *RootNode, err error) {
	upperNode = NewRootNode()
	err = upperNode.Parse(input)
	return
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
		if b == '\\' {
			b, err = in.ReadByte()
			actualByte, ok := escapedLetters[b]
			if ok {
				b = actualByte
			}
			if err != nil {
				return
			}
		}
		read = append(read, b)
		next, _ = in.Peek(1)
	}
	return
}

func parse(in *bufio.Reader) (result TreeNode, err error) {
	readWhile(unicode.IsSpace, in)

	next, err := in.Peek(2)
	var read []byte

	// Nothing to read, eof
	if len(next) == 0 {
		return nil, doneReading
	}

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
		for {
			readWhile(unicode.IsSpace, in)
			next, err = in.Peek(1)
			if err != nil {
				return nil, fmt.Errorf(") expected but %s was found", err)
			}
			if next[0] == ')' {
				break
			}
			var child TreeNode
			child, err = parse(in)
			if err != nil {
				return
			}

			args = append(args, child)
		}

		_, err = in.ReadByte()
		result = &FunctionNode{Name: name, Args: args}
		return
	case next[0] == '"':
		in.ReadByte()
		// read string literal
		read, err = readWhile(func(r rune) bool { return r != '"' }, in)
		if err != nil {
			return nil, fmt.Errorf(`" expected but %s`, err)
		}

		_, err = in.ReadByte()
		if err != nil {
			return nil, fmt.Errorf(`" expected but %s`, err)
		}

		result = &ValueNode{string(read)}
		return
	case len(next) > 1 && next[0] == '-' && unicode.IsDigit(rune(next[1])):
		in.ReadByte()
		// read negative number
		read, err = readWhile(unicode.IsDigit, in)
		if err != nil {
			return
		}
		var n int
		n, err = strconv.Atoi(string(read))
		result = &ValueNode{-n}
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
	case isFirstSymbolRune(rune(next[0])):
		// read symbol (which is actually just a string)
		read, err = readWhile(isSymbolRune, in)
		if err != nil {
			return
		}
		symbol := string(read)

		if symbol == "nil" {
			result = &ValueNode{nil}
		} else {
			result = &SymbolNode{symbol}
		}
		return
	default:
		err = fmt.Errorf("invalid symbol %s", next)
		return
	}
	return
}
