package b

import "go-cyclic/example/simple_cyclic/a"

type B struct {
	A *a.A
}
