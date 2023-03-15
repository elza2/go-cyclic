package a

import (
	"go-cyclic/example/three_cyclic2/b"
)

type A struct {
	B *b.B
}

func NewA() *A {
	return &A{}
}
