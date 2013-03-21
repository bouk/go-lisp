package lisp

import (
	"errors"
	"fmt"
	"strconv"
)

func invalidTypeError(expected string, actual interface{}) (err error) {
	return fmt.Errorf("Invalid type, %v expected but %T given", expected, actual)
}

func evaluateArgs(s *Scope, args []TreeNode) (result []Value, err error) {
	result = make([]Value, len(args))
	for i, node := range args {
		result[i], err = node.Interpret(s)
		if err != nil {
			return nil, err
		}
	}
	return
}

func builtinFunctions(s *Scope) {
	s.RegisterFunctionAliases([]string{"*", "mult"}, func(s *Scope, args []TreeNode) (result Value, err error) {
		if len(args) != 2 {
			return 0, errors.New("not right number of arguments for multiply, should be two")
		}

		evaluatedArgs, err := evaluateArgs(s, args)
		if err != nil {
			return nil, err
		}

		var intResult = 1
		for _, argument := range evaluatedArgs {
			switch argument.(type) {
			case int:
			default:
				return nil, invalidTypeError("int", argument)
			}
			intResult *= argument.(int)
		}
		result = intResult
		return
	})

	s.RegisterFunctionAliases([]string{"+", "add"}, func(s *Scope, args []TreeNode) (Value, error) {
		if len(args) != 2 {
			return 0, errors.New("not right number of arguments for add, should be two")
		}

		evaluatedArgs, err := evaluateArgs(s, args)
		if err != nil {
			return nil, err
		}
		switch evaluatedArgs[0].(type) {
		case int:
			var result int = evaluatedArgs[0].(int)
			for _, argument := range evaluatedArgs[1:] {
				switch argument.(type) {
				case int:
				default:
					return nil, invalidTypeError("int", argument)
				}
				result += argument.(int)
			}
			return result, nil
		case string:
			var result string = evaluatedArgs[0].(string)
			for _, argument := range evaluatedArgs[1:] {
				switch argument.(type) {
				case string:
					result += argument.(string)
				case int:
					result += strconv.Itoa(argument.(int))
				default:
					return nil, fmt.Errorf("unknown type %T", argument)
				}
			}
			return result, nil
		default:
			return nil, fmt.Errorf("unknown type %T", evaluatedArgs[0])
		}

		return nil, nil
	})

	s.RegisterFunctionAliases([]string{"-", "sub"}, func(s *Scope, args []TreeNode) (result Value, res error) {
		if len(args) != 2 {
			return 0, errors.New("not right number of arguments for subtract, should be two")
		}

		evaluatedArgs, err := evaluateArgs(s, args)
		if err != nil {
			return nil, err
		}

		switch evaluatedArgs[0].(type) {
		case int:
		default:
			return nil, invalidTypeError("int", evaluatedArgs[0])
		}
		intResult := evaluatedArgs[0].(int)
		for _, argument := range evaluatedArgs[1:] {
			switch argument.(type) {
			case int:
			default:
				return nil, invalidTypeError("int", argument)
			}
			intResult -= argument.(int)
		}
		result = intResult
		return
	})

	s.RegisterFunctionAliases([]string{"/", "div"}, func(s *Scope, args []TreeNode) (result Value, err error) {
		if len(args) != 2 {
			return 0, errors.New("invalid number of arguments for divide, has to be two")
		}

		evaluatedArgs, err := evaluateArgs(s, args)
		if err != nil {
			return
		}

		firstVal, firstOk := evaluatedArgs[0].(int)
		if !firstOk {
			return nil, invalidTypeError("int", evaluatedArgs[0])
		}
		secondVal, secondOk := evaluatedArgs[1].(int)
		if !secondOk {
			return nil, invalidTypeError("int", evaluatedArgs[1])
		}

		result = firstVal / secondVal
		return
	})

	s.RegisterFunction("set", func(s *Scope, args []TreeNode) (value Value, err error) {
		if len(args) != 2 {
			return 0, errors.New("set requires two arguments")
		}

		node, ok := args[0].(*SymbolNode)
		if !ok {
			return nil, errors.New("set requires symbol")
		}

		value, err = args[1].Interpret(s)

		s.SetVariable(node.Name, value)
		return
	})

	s.RegisterFunction("get", func(s *Scope, args []TreeNode) (value Value, err error) {
		if len(args) != 1 {
			return 0, errors.New("get requires one argument")
		}

		node, ok := args[0].(*SymbolNode)
		if !ok {
			return nil, errors.New("get requires symbol")
		}

		value = s.GetVariable(node.Name)
		return
	})

	s.RegisterFunction("scope", func(s *Scope, args []TreeNode) (value Value, err error) {
		innerScope := NewScope(s)

		for _, child := range args {
			value, err = child.Interpret(innerScope)
			if err != nil {
				return nil, err
			}
		}

		if innerScope.Variables["_return"] != nil {
			value = innerScope.Variables["_return"]
		}
		return
	})
}
