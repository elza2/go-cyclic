package tool

import (
	"path/filepath"

	"github.com/elza2/go-cyclic/resolver"
	"github.com/elza2/go-cyclic/sprite"
	"github.com/elza2/go-cyclic/topology"
)

type Params struct {
	Dir     string
	Filters []string
}

func CheckCycleDepend(params *Params) error {
	abs, err := filepath.Abs(params.Dir)
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
	sprites, err := resolver.ParseNodeSprite(root, module, abs, params.Filters)
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
