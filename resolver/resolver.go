package resolver

import (
	"fmt"
	"go-cyclic/errors"
	"go-cyclic/sprite"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"strings"
)

var (
	GoMod = "go.mod"

	fset = token.NewFileSet()
)

// ParseDir parse path.
func ParseDir(dir string) string {
	idx := strings.LastIndex(dir, "/")
	if idx == -1 {
		return dir
	}
	return dir[idx+1:]
}

func ParseGoModule(dir string) (moduleName string, err error) {
	readFile, err := ioutil.ReadFile(dir + "/" + GoMod)
	if err != nil {
		return "", errors.GoModNotExist
	}
	modFile, err := modfile.Parse(GoMod, readFile, nil)
	if err != nil {
		return "", errors.GoModParseFailed
	}
	return modFile.Module.Mod.Path, nil
}

func ParseNodeSprite(root string, module string, dir string) (nodes []*sprite.NodeSprite, err error) {
	nodeSprites := make([]*sprite.NodeSprite, 0)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	parseDir, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return nodeSprites, err
	}
	for k, v := range parseDir {
		nodeImports := map[string][]string{}
		for key, value := range v.Files {
			for _, port := range value.Imports {
				nodeImports[key] = append(nodeImports[key], port.Path.Value[1:len(port.Path.Value)-1])
			}
			filePath, nodeName := GetNodeSpritePathName(key)
			nodeSprites = append(nodeSprites, &sprite.NodeSprite{
				FilePath:    filePath,
				RootName:    root,
				ModuleName:  module,
				PackageName: k,
				NodeName:    nodeName,
				ImportNames: nodeImports[key],
			})
		}
	}
	for _, file := range files {
		if file.IsDir() {
			sprites, err := ParseNodeSprite(root, module, fmt.Sprintf("%s/%s", dir, file.Name()))
			if err != nil {
				return sprites, err
			}
			nodeSprites = append(nodeSprites, sprites...)
		}
	}
	return nodeSprites, nil
}

func GetNodeSpritePathName(content string) (filePath, nodeName string) {
	contains := strings.Contains(content, "/")
	if !contains {
		return content, content
	}
	index := strings.LastIndex(content, "/")
	return content[:index], content[index+1:]
}
