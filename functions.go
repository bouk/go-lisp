package lisp

import (
	"errors"
)

type Argument int
type Function func(args []int) (result int, err error)

var registeredFunctions map[string]Function

func registerFunction(names []string, f Function) {
	for _, name := range names {
		registeredFunctions[name] = f
	}
}

func init() {
	registeredFunctions = make(map[string]Function)

	registerFunction([]string{"*", "mult"}, func(args []int) (result int, err error) {
		if len(args) < 1 {
			return 0, errors.New("not enough arguments for multiply, needs at least two")
		}

		result = 1
		for _, argument := range args {
			result *= argument
		}
		return
	})

	registerFunction([]string{"+", "add"}, func(args []int) (result int, err error) {
		if len(args) < 1 {
			return 0, errors.New("not enough arguments for add, needs at least two")
		}

		result = 0
		for _, argument := range args {
			result += argument
		}
		return
	})

	registerFunction([]string{"-", "sub"}, func(args []int) (result int, err error) {
		if len(args) < 1 {
			return 0, errors.New("not enough arguments for subtract, needs at least two")
		}

		result = args[0]
		for _, argument := range args[1:] {
			result -= argument
		}
		return
	})

	registerFunction([]string{"/", "div"}, func(args []int) (result int, err error) {
		if len(args) != 2 {
			return 0, errors.New("invalid number of arguments for divide")
		}

		result = args[0] / args[1]
		return
	})
}
