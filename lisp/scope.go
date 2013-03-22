package lisp

import (
	"bufio"
	"io"
)

type Scope struct {
	Parent    *Scope
	Out       io.Writer
	In        *bufio.Reader
	Variables map[string]Value
	Functions map[string]Function
}

func (s *Scope) FindFunction(name string) (f Function) {
	f, found := s.Functions[name]
	if !found && s.Parent != nil {
		f = s.Parent.FindFunction(name)
	}
	return
}

func (s *Scope) RegisterFunctionAliases(names []string, f Function) {
	for _, name := range names {
		s.Functions[name] = f
	}
}

func (s *Scope) RegisterFunction(name string, f Function) {
	s.Functions[name] = f
}

// Gets a variable, looking in higher scopes if necessary. Returns nil if value not found
func (s *Scope) GetVariable(name string) (val Value) {
	val, found := s.Variables[name]
	if !found && s.Parent != nil {
		val = s.Parent.GetVariable(name)
	}
	return
}

// Sets a variable on the first scope that has it. Set's it on this scope if it does not exist
func (s *Scope) SetVariable(name string, value Value) {
	if !s.setVariable(name, value) {
		s.Variables[name] = value
	}
	return
}

func (s *Scope) setVariable(name string, value Value) (found bool) {
	_, found = s.Variables[name]
	if found {
		s.Variables[name] = value
	} else {
		if s.Parent != nil {
			found = s.Parent.setVariable(name, value)
		}
	}
	return
}

func NewScope(parent *Scope) (s *Scope) {
	s = &Scope{Variables: make(map[string]Value), Functions: make(map[string]Function)}
	if parent != nil {
		s.Parent = parent
		s.Out = parent.Out
		s.In = parent.In
	}
	return
}
