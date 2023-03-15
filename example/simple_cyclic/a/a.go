package a

import (
	"go-cyclic/example/simple_cyclic/b"
)

type A struct {
	B *b.B
}

func NewA() *A {
	return &A{}
}
