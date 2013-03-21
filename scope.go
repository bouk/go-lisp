package lisp

type Scope struct {
	Parent    *Scope
	Variables map[string]Value
	Functions map[string]Function
}

func (s *Scope) FindVariable(name string) (val Value) {
	val, found := s.Variables[name]
	if !found && s.Parent != nil {
		val = s.Parent.FindVariable(name)
	}
	return val
}

func (s *Scope) FindFunction(name string) (f Function) {
	f, found := s.Functions[name]
	if !found && s.Parent != nil {
		f = s.Parent.FindFunction(name)
	}
	return f
}

func (s *Scope) RegisterFunction(names []string, f Function) {
	for _, name := range names {
		s.Functions[name] = f
	}
}

func NewScope(parent *Scope) *Scope {
	return &Scope{parent, make(map[string]Value), make(map[string]Function)}
}
