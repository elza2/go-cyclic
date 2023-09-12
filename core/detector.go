package core

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/mod/modfile"
)

var (
	fset = token.NewFileSet()

	Failed  = color.RedString("%v", "Failed.")
	Success = color.GreenString("%v", "Success.")
)

// NodeResolver used to parse node.
type NodeResolver struct {
	// full file path.
	FullPath string
	// absolute file path.
	AbsFilePath string
	// relative file path.
	RelFilePath string
	// package name.
	PackageName string
	// package ref path.
	PackagePath string
	// node name.
	NodeName string
	// import values.
	ImportValues []string
}

// CyclicDetector used to cyclic node detection.
type CyclicDetector struct {
	// detection dir param.
	Dir string
	// source filter param.
	SourceFilter string
	// detection ignored.
	Filters []string
	// dir absolute file path.
	AbsPath string
	// Dir the root node of the project.
	RootPath string
	// parse module name.
	ModuleName string
	// parse nodes.
	NodeResolvers []*NodeResolver
	// package ref path to nodes.
	NodePackageMap map[string][]*NodeResolver
	// full path map nodes collections.
	NodeResolversMap map[string]*NodeResolver
}

func Do(dir, filter string) (err error) {
	detector, err := NewCyclicDetector(dir, filter)
	if err != nil {
		return err
	}
	relations, degrees := detector.buildNodeRelation()
	queue := make([]*NodeResolver, 0)
	for full, degree := range degrees {
		if degree == 0 {
			queue = append(queue, detector.NodeResolversMap[full])
		}
	}
	for len(queue) != 0 {
		q := queue
		queue = nil
		for _, node := range q {
			for i := 0; i < len(relations); i++ {
				if relations[i][1].FullPath == node.FullPath {
					degrees[relations[i][0].FullPath]--
					if degrees[relations[i][0].FullPath] == 0 {
						queue = append(queue, relations[i][0])
					}
				}
			}
			delete(degrees, node.FullPath)
		}
	}
	if len(degrees) == 0 {
		fmt.Println(Success + " Not circular dependence.")
		return
	}
	detector.PrintCycleDepend(relations, degrees)
	return nil
}

func NewCyclicDetector(dir string, filter string) (detector *CyclicDetector, err error) {
	detector = &CyclicDetector{
		Dir:              dir,
		SourceFilter:     filter,
		NodePackageMap:   make(map[string][]*NodeResolver, 0),
		NodeResolversMap: make(map[string]*NodeResolver, 0),
	}
	if err = detector.Parse(); err != nil {
		return nil, err
	}
	return detector, err
}

func (c *CyclicDetector) Parse() error {
	abs, err := filepath.Abs(c.Dir)
	if err != nil {
		return err
	}
	stat, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("not exist path: %v", abs))
		}
		return err
	}
	if !stat.IsDir() {
		return errors.New(fmt.Sprintf("not a folder path: %v", abs))
	}
	readFile, err := os.ReadFile(abs + "/go.mod")
	if err != nil {
		return errors.New("not find go.mod file, path: " + abs)
	}
	modFile, err := modfile.Parse("go.mod", readFile, nil)
	if err != nil {
		return errors.New("go.mod file parse failed, path: " + abs + "/go.mod")
	}

	c.AbsPath = abs
	c.RootPath = abs[strings.LastIndex(abs, "/")+1:]
	c.ModuleName = modFile.Module.Mod.Path

	if err = c.ParseFilters(c.SourceFilter); err != nil {
		return err
	}
	if err = c.ParseResolvers(); err != nil {
		return err
	}
	return nil
}

func (c *CyclicDetector) ParseFilters(filter string) (err error) {
	if filter == "" {
		return nil
	}
	if strings.Contains(filter, "，") {
		return errors.New("not supported `，` symbol, please use `,` split")
	}
	filters := strings.Split(filter, ",")
	for i, f := range filters {
		if strings.Contains(f, ".go") {
			continue
		}
		filters[i] += ".go"
	}

	c.Filters = filters

	return nil
}

func (c *CyclicDetector) ParseResolvers() error {
	nodeResolvers, err := c.ParseNodeResolver(c.AbsPath)
	if err != nil {
		return err
	}

	c.NodeResolvers = nodeResolvers
	for _, node := range c.NodeResolvers {
		c.NodePackageMap[node.PackagePath] = append(c.NodePackageMap[node.PackagePath], node)
		c.NodeResolversMap[node.FullPath] = node
	}

	return nil
}

func (c *CyclicDetector) ParseNodeResolver(path string) ([]*NodeResolver, error) {
	nodeResolvers := make([]*NodeResolver, 0)
	parseFiles, err := parser.ParseDir(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}
	for pkg, file := range parseFiles {
		importValues := map[string][]string{}
		for key, value := range file.Files {
			filePath, nodeName := GetNodePathName(key)
			if c.HasFilter(nodeName) {
				continue
			}
			for _, imp := range value.Imports {
				importValues[key] = append(importValues[key], imp.Path.Value[1:len(imp.Path.Value)-1])
			}
			relFilePath := strings.ReplaceAll(filePath, c.AbsPath, "")
			nodeResolvers = append(nodeResolvers, &NodeResolver{
				FullPath:     key,
				AbsFilePath:  filePath,
				RelFilePath:  relFilePath,
				PackageName:  pkg,
				PackagePath:  c.ModuleName + relFilePath,
				NodeName:     nodeName,
				ImportValues: importValues[key],
			})
		}
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			nodes, err := c.ParseNodeResolver(fmt.Sprintf("%s/%s", path, file.Name()))
			if err != nil {
				return nil, err
			}
			nodeResolvers = append(nodeResolvers, nodes...)
		}
	}
	return nodeResolvers, nil
}

func GetNodePathName(content string) (filePath, nodeName string) {
	contains := strings.Contains(content, "/")
	if !contains {
		return content, content
	}
	index := strings.LastIndex(content, "/")
	return content[:index], content[index+1:]
}

func (c *CyclicDetector) HasFilter(name string) bool {
	for _, p := range c.Filters {
		if regexp.MustCompile("^" + p + "$").MatchString(name) {
			return true
		}
	}
	return false
}

func (c *CyclicDetector) buildNodeRelation() ([][]*NodeResolver, map[string]int) {
	degrees := make(map[string]int, 0)
	relations := make([][]*NodeResolver, 0)
	for _, node := range c.NodeResolvers {
		depends := make([][]*NodeResolver, 0)
		for _, in := range node.ImportValues {
			subPkg := in[strings.LastIndex(in, "/")+1:]
			pkgs := c.NodePackageMap[in]
			pkgMap := make(map[string]int, 0)
			simple := true
			for _, pkg := range pkgs {
				pkgMap[pkg.PackageName] = 1
				if len(pkgMap) > 1 {
					simple = false
					break
				}
			}
			for _, pkg := range pkgs {
				if (simple || subPkg == pkg.PackageName) && node.FullPath != pkg.FullPath {
					depends = append(depends, []*NodeResolver{node, pkg})
				}
			}
		}
		degrees[node.FullPath] = len(depends)
		relations = append(relations, depends...)
	}
	return relations, degrees
}

func (c *CyclicDetector) PrintCycleDepend(relations [][]*NodeResolver, degrees map[string]int) {
	depends := map[string][]*NodeResolver{}
	for full := range degrees {
		for i := 0; i < len(relations); i++ {
			if _, ok := degrees[relations[i][1].FullPath]; ok && relations[i][0].FullPath == full {
				key := relations[i][1].FullPath
				depends[key] = append(depends[key], c.NodeResolversMap[full])
			}
		}
	}

	search := make(map[string]int, 0)
	traces := make([]string, 0)
	cycles := make([][]string, 0)
	for key := range depends {
		if _, ok := search[key]; ok {
			continue
		}
		cycles = c.QueryCycleNodes(c.NodeResolversMap[key], depends, search, traces, cycles)
	}
	fmt.Printf(Failed+" %v circular dependence chains were found. \n\n", len(cycles))
	for _, cycle := range cycles {
		PrintCycle(cycle)
		fmt.Print("\n\n")
	}
}

func (c *CyclicDetector) QueryCycleNodes(
	node *NodeResolver,
	depends map[string][]*NodeResolver,
	search map[string]int,
	traces []string,
	cycles [][]string,
) [][]string {
	j := -1
	for i := 0; i < len(traces); i++ {
		if traces[i] == node.FullPath {
			j = i
		}
	}
	if j != -1 {
		circle := make([]string, 0)
		for ; j < len(traces); j++ {
			circle = append(circle, traces[j])
		}
		if len(cycles) == 0 {
			cycles = append(cycles, circle)
		} else {
			flag := false
			for _, cycle := range cycles {
				if len(cycle) == len(circle) {
					flag = true
				}
			}
			if !flag {
				cycles = append(cycles, circle)
			} else {
				flag = false
				for _, cycle := range cycles {
					if len(cycle) == len(circle) {
						cycleMap := map[string]int{}
						for _, c := range cycle {
							cycleMap[c] = 1
						}
						for _, c := range circle {
							if _, ok := cycleMap[c]; !ok {
								cycles = append(cycles, circle)
								return cycles
							}
						}
					}
				}
			}
		}
		return cycles
	}
	traces = append(traces, node.FullPath)
	for _, v := range depends[node.FullPath] {
		search[v.FullPath] = 1
		cycles = c.QueryCycleNodes(v, depends, search, traces, cycles)
	}
	traces = traces[:len(traces)-1]
	return cycles
}

func PrintCycle(cycle []string) {
	sort.Slice(cycle, func(i, j int) bool {
		return cycle[i] > cycle[j]
	})
	downMark, vertical, startPre, endPre := "↓", "┆", "┌---→", "└--- "
	output := make([]string, 0)
	length := -1
	for i := 0; i < len(cycle); i++ {
		length = Max(length, len(cycle[i]))
	}
	centre := length >> 1
	blank := fillBlank(centre + 2)
	for i := len(cycle) - 1; i >= 0; i-- {
		switch i {
		case 0:
			output = append(output, vertical+blank+downMark)
			output = append(output, endPre+fillBlank(centre-(len(cycle[i])>>1)+2)+cycle[i])
		case len(cycle) - 1:
			output = append(output, startPre+fillBlank(centre-(len(cycle[i])>>1)+2)+cycle[i])
		default:
			output = append(output, vertical+blank+downMark)
			output = append(output, vertical+fillBlank(centre-(len(cycle[i])>>1)+len(startPre)-3)+cycle[i])
		}
	}
	for _, out := range output {
		fmt.Println(out)
	}
}

func fillBlank(n int) string {
	blank := ""
	for i := 0; i < n; i++ {
		blank += " "
	}
	return blank
}

func Max(m, n int) int {
	if m > n {
		return m
	}
	return n
}
