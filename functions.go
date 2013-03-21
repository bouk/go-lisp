package lisp

type Value interface{}

type Function func(args []Value) (result Value, err error)

var (
	registeredFunctions map[string]Function = make(map[string]Function)
)

func registerFunction(names []string, f Function) {
	for _, name := range names {
		registeredFunctions[name] = f
	}
}
