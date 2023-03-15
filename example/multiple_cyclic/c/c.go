package c

import (
	"go-cyclic/example/multiple_cyclic/a"
	"go-cyclic/example/multiple_cyclic/d"
)

type C struct {
	A *a.A
	D *d.D
}
