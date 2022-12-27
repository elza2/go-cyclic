package resolver

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/elza2/go-cyclic/errors"
	"github.com/elza2/go-cyclic/sprite"
)

var (
	GoMod = "go.mod"

	fset = token.NewFileSet()
)

// ParseDir parse path.
func ParseDir(dir string) (path string, err error) {
	stat, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.PathNotExist(dir)
		}
		return "", err
	}
	if !stat.IsDir() {
		return "", errors.PathNotIsFile(dir)
	}
	idx := strings.LastIndex(dir, "/")
	if idx == -1 {
		return dir, nil
	}
	return dir[idx+1:], nil
}

func ParseGoModule(dir string) (module string, err error) {
	readFile, err := os.ReadFile(dir + "/" + GoMod)
	if err != nil {
		return "", errors.GoModNotExist(dir)
	}
	modFile, err := modfile.Parse(GoMod, readFile, nil)
	if err != nil {
		return "", errors.GoModParseFailed(dir)
	}
	return modFile.Module.Mod.Path, nil
}

func ParseNodeSprite(root string, module string, dir string, filters []string) (nodes []*sprite.NodeSprite, err error) {
	nodeSprites := make([]*sprite.NodeSprite, 0)
	files, err := os.ReadDir(dir)
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
			if HasFilterSprite(nodeName, filters) {
				continue
			}
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
			sprites, err := ParseNodeSprite(root, module, fmt.Sprintf("%s/%s", dir, file.Name()), filters)
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

func HasFilterSprite(name string, filters []string) bool {
	for _, p := range filters {
		if regexp.MustCompile("^" + p + "$").MatchString(name) {
			return true
		}
	}
	return false
}
