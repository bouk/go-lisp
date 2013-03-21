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

func evaluateToString(s *Scope, node TreeNode) (str string, err error) {
	funNameValue, err := node.Interpret(s)
	if err != nil {
		return
	}
	str, ok := funNameValue.(string)
	if !ok {
		err = errors.New("should be string")
	}
	return
}

func builtinFunctions(defaultScope *Scope) {
	defaultScope.RegisterFunctionAliases([]string{"*", "mult"}, func(s *Scope, args []TreeNode) (result Value, err error) {
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

	defaultScope.RegisterFunctionAliases([]string{"+", "add"}, func(s *Scope, args []TreeNode) (Value, error) {
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

	defaultScope.RegisterFunctionAliases([]string{"-", "sub"}, func(s *Scope, args []TreeNode) (result Value, res error) {
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

	defaultScope.RegisterFunctionAliases([]string{"/", "div"}, func(s *Scope, args []TreeNode) (result Value, err error) {
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

	defaultScope.RegisterFunction("set", func(s *Scope, args []TreeNode) (value Value, err error) {
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

	defaultScope.RegisterFunction("get", func(s *Scope, args []TreeNode) (value Value, err error) {
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

	defaultScope.RegisterFunction("scope", func(s *Scope, args []TreeNode) (value Value, err error) {
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

	defaultScope.RegisterFunction("defun", func(s *Scope, args []TreeNode) (value Value, err error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("defun requires at least 2 arguments, %d given", len(args))
		}

		symbolNode, ok := args[0].(*SymbolNode)
		if !ok {
			return nil, errors.New("function name should be a symbol")
		}
		name := symbolNode.Name

		funArgs := make([]string, len(args)-2)
		for i, argument := range args[1 : len(args)-1] {
			symbolNode, ok = argument.(*SymbolNode)
			if !ok {
				return nil, errors.New("function argument name should be a symbol")
			}
			funArgs[i] = symbolNode.Name
		}

		program := args[len(args)-1]

		s.RegisterFunction(name, func(funScope *Scope, args []TreeNode) (Value, error) {
			if len(args) != len(funArgs) {
				return nil, fmt.Errorf("function %s needs %d arguments", name, len(funArgs))
			}

			innerScope := NewScope(funScope)
			for i, arg := range funArgs {
				innerScope.Variables[arg], err = args[i].Interpret(funScope)
				if err != nil {
					return nil, err
				}
			}
			return program.Interpret(innerScope)
		})
		return
	})

	defaultScope.RegisterFunction("print", func(scope *Scope, args []TreeNode) (v Value, err error) {
		evaluatedArgs, err := evaluateArgs(scope, args)
		if err != nil {
			return nil, err
		}

		for _, val := range evaluatedArgs {
			switch val.(type) {
			case int:
				fmt.Fprint(scope.Out, val.(int))
			case string:
				fmt.Fprint(scope.Out, val.(string))
			default:
				fmt.Fprintf(scope.Out, "%T", val)
			}
		}

		return
	})
}
