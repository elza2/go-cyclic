package main

import (
	"go-cyclic/tool"
	"testing"
)

func TestCyclic(t *testing.T) {
	tool.CheckCycleDepend("/Users/yuanyou/go/src/daji")
}
