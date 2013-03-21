package lisp

import (
	"errors"
	"fmt"
	"strconv"
)

func invalidTypeError(expected string, actual interface{}) (err error) {
	return fmt.Errorf("Invalid type, %v expected but %T given", expected, actual)
}

func init() {
	registerFunction([]string{"*", "mult"}, func(args []Value) (Value, error) {
		if len(args) < 2 {
			return 0, errors.New("not enough arguments for multiply, needs at least two")
		}
		var result int = 1
		for _, argument := range args {
			switch argument.(type) {
			case int:
			default:
				return nil, invalidTypeError("int", argument)
			}
			result *= argument.(int)
		}
		return result, nil
	})

	registerFunction([]string{"+", "add"}, func(args []Value) (Value, error) {
		if len(args) < 2 {
			return 0, errors.New("not enough arguments for add, needs at least two")
		}

		switch args[0].(type) {
		case int:
			var result int = args[0].(int)
			for _, argument := range args[1:] {
				switch argument.(type) {
				case int:
				default:
					return nil, invalidTypeError("int", argument)
				}
				result += argument.(int)
			}
			return result, nil
		case string:
			var result string = args[0].(string)
			for _, argument := range args[1:] {
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
			return nil, fmt.Errorf("unknown type %T", args[0])
		}

		return nil, nil
	})

	registerFunction([]string{"-", "sub"}, func(args []Value) (Value, error) {
		if len(args) < 2 {
			return 0, errors.New("not enough arguments for subtract, needs at least two")
		}
		switch args[0].(type) {
		case int:
		default:
			return nil, invalidTypeError("int", args[0])
		}
		result := args[0].(int)
		for _, argument := range args[1:] {
			switch argument.(type) {
			case int:
			default:
				return nil, invalidTypeError("int", argument)
			}
			result -= argument.(int)
		}
		return result, nil
	})

	registerFunction([]string{"/", "div"}, func(args []Value) (Value, error) {
		if len(args) != 2 {
			return 0, errors.New("invalid number of arguments for divide, has to be two")
		}

		switch args[0].(type) {
		case int:
			switch args[1].(type) {
			case int:
			default:
				return nil, invalidTypeError("int", args[1])
			}
		default:
			return nil, invalidTypeError("int", args[0])
		}

		result := args[0].(int) / args[1].(int)
		return result, nil
	})
}
