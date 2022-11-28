package tool

import (
	"github.com/elza2/go-cyclic/resolver"
	"github.com/elza2/go-cyclic/sprite"
	"github.com/elza2/go-cyclic/topology"
	"path/filepath"
)

func CheckCycleDepend(dir string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	// parse root path.
	root, err := resolver.ParseDir(abs)
	if err != nil {
		return err
	}
	// parse module name.
	module, err := resolver.ParseGoModule(abs)
	if err != nil {
		return err
	}
	// parse sprite nodes by path.
	sprites, err := resolver.ParseNodeSprite(root, module, abs)
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
