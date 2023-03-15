package core_test

import (
	"go-cyclic/core"
	"testing"
)

func TestSimpleCyclic(t *testing.T) {
	err := core.Do("../example/simple_cyclic", "")
	if err != nil {
		return
	}
}

func TestSimple2Cyclic(t *testing.T) {
	err := core.Do("../example/simple_cyclic2", "")
	if err != nil {
		return
	}
}

func TestSimple3Cyclic(t *testing.T) {
	err := core.Do("../example/simple_cyclic3", "")
	if err != nil {
		return
	}
}

func TestThreeCyclic(t *testing.T) {
	err := core.Do("../example/three_cyclic", "")
	if err != nil {
		return
	}
}

func TestThree2Cyclic(t *testing.T) {
	err := core.Do("../example/three_cyclic2", "")
	if err != nil {
		return
	}
}

func TestMultipleCyclic(t *testing.T) {
	err := core.Do("../example/multiple_cyclic", "")
	if err != nil {
		return
	}
}
