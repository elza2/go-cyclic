package sprite

import (
	"fmt"
	"strings"
)

type NodeSprites struct {
	nodeSprites []*NodeSprite
}

func (nodes *NodeSprites) SetNodeSprites(sprites []*NodeSprite) {
	nodes.nodeSprites = sprites
}

func (nodes *NodeSprites) GetNodeSprites() []*NodeSprite {
	return nodes.nodeSprites
}

type NodeSprite struct {
	// full file path.
	FilePath string
	// root name.
	RootName string
	// module name.
	ModuleName string
	// package name.
	PackageName string
	// node name.
	NodeName string
	// import names.
	ImportNames []string
}

// MatchImportNodeSprite matches the specified sprite.
func (nodes *NodeSprites) MatchImportNodeSprite(match string) []*NodeSprite {
	nodeSprites := make([]*NodeSprite, 0)
	for _, sprite := range nodes.nodeSprites {
		path := sprite.GetRootPath()
		if path == match {
			nodeSprites = append(nodeSprites, sprite)
		}
	}
	return nodeSprites
}

func (nodes *NodeSprites) GetNodeSprite(allPath string) *NodeSprite {
	for _, sprite := range nodes.nodeSprites {
		if sprite.GetAllPath() == allPath {
			return sprite
		}
	}
	return nil
}

func (node *NodeSprite) GetFilePath() string {
	return node.FilePath
}

func (node *NodeSprite) GetAllPath() string {
	return fmt.Sprintf("%s/%s", node.FilePath, node.NodeName)
}

func (node *NodeSprite) GetRootName() string {
	return node.RootName
}

func (node *NodeSprite) GetRootPath() string {
	idx := strings.LastIndex(node.FilePath, "/"+node.RootName)
	if idx == -1 {
		return node.FilePath
	}
	return node.ModuleName + node.FilePath[idx+len(node.RootName)+1:]
}

func (node *NodeSprite) GetImportNames() []string {
	return node.ImportNames
}
