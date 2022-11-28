package tool

import (
	"go-cyclic/resolver"
	"go-cyclic/sprite"
	"go-cyclic/topology"
)

func CheckCycleDepend(dir string) error {
	// parse root path.
	root := resolver.ParseDir(dir)
	// parse module name.
	module, err := resolver.ParseGoModule(dir)
	if err != nil {
		return err
	}
	// parse sprite nodes by path.
	sprites, err := resolver.ParseNodeSprite(root, module, dir)
	if err != nil {
		return err
	}
	nodeSprites := new(sprite.NodeSprites)
	nodeSprites.SetNodeSprites(sprites)
	// constructor topology struct.
	topologies := topology.ConstructorTopology(nodeSprites)
	// check cycle depend.
	hasCycle := topologies.CycleDepend()
	// print cycle depend.
	topologies.PrintCycleDepend(hasCycle)
	return nil
}
