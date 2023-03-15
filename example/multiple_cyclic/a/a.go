package a

import (
	"go-cyclic/example/multiple_cyclic/b"
)

type A struct {
	B *b.B
}

func NewA() *A {
	return &A{}
}
