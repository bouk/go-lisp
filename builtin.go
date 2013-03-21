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
	s.RegisterFunction([]string{"*", "mult"}, func(s *Scope, args []TreeNode) (Value, error) {
		if len(args) != 2 {
			return 0, errors.New("not right number of arguments for multiply, should be two")
		}

		evaluatedArgs, err := evaluateArgs(s, args)
		if err != nil {
			return nil, err
		}

		var result int = 1
		for _, argument := range evaluatedArgs {
			switch argument.(type) {
			case int:
			default:
				return nil, invalidTypeError("int", argument)
			}
			result *= argument.(int)
		}
		return result, nil
	})

	s.RegisterFunction([]string{"+", "add"}, func(s *Scope, args []TreeNode) (Value, error) {
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

	s.RegisterFunction([]string{"-", "sub"}, func(s *Scope, args []TreeNode) (Value, error) {
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
		result := evaluatedArgs[0].(int)
		for _, argument := range evaluatedArgs[1:] {
			switch argument.(type) {
			case int:
			default:
				return nil, invalidTypeError("int", argument)
			}
			result -= argument.(int)
		}
		return result, nil
	})

	s.RegisterFunction([]string{"/", "div"}, func(s *Scope, args []TreeNode) (Value, error) {
		if len(args) != 2 {
			return 0, errors.New("invalid number of arguments for divide, has to be two")
		}

		evaluatedArgs, err := evaluateArgs(s, args)
		if err != nil {
			return nil, err
		}
		switch evaluatedArgs[0].(type) {
		case int:
			switch evaluatedArgs[1].(type) {
			case int:
			default:
				return nil, invalidTypeError("int", evaluatedArgs[1])
			}
		default:
			return nil, invalidTypeError("int", evaluatedArgs[0])
		}

		result := evaluatedArgs[0].(int) / evaluatedArgs[1].(int)
		return result, nil
	})
}
