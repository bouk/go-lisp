package lisp

type Argument int
type Function func(args []int) (result int, err error)

var registeredFunctions map[string]Function

func registerFunction(names []string, f Function) {
	for _, name := range names {
		registeredFunctions[name] = f
	}
}
