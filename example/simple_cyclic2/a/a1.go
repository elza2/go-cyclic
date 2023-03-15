package a

import "go-cyclic/example/simple_cyclic2/b"

type A1 struct {
	B *b.B
}

func NewA1() *A1 {
	return &A1{}
}
