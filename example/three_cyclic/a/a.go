package a

import (
	"go-cyclic/example/three_cyclic/b"
)

type A struct {
	B *b.B
}

func NewA() *A {
	return &A{}
}
