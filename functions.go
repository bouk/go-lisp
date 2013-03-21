package lisp

type Value interface{}

type Function func(s *Scope, args []TreeNode) (result Value, err error)
