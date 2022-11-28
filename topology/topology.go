package topology

import (
	"fmt"
	"github.com/fatih/color"
	"go-cyclic/sprite"
	"strings"
)

var (
	Failed  = color.RedString("%v", "Failed.")
	Success = color.GreenString("%v", "Success.")
)

const (
	test = "_test"

	NotCycle  = " Not circular dependence."
	Line      = "\n\n"
	TextCycle = " %v circular dependence chains were found." + Line
)

type Topology struct {
	NodeSprites *sprite.NodeSprites
	Degrees     map[string]int
	Relations   [][]*sprite.NodeSprite
	Cycles      [][]string
	Trace       []string
	Search      map[string]int
}

func ConstructorTopology(nodeSprites *sprite.NodeSprites) (topology *Topology) {
	topology = &Topology{
		NodeSprites: nodeSprites,
		Degrees:     make(map[string]int, 0),
		Relations:   make([][]*sprite.NodeSprite, 0),
		Cycles:      make([][]string, 0),
		Trace:       make([]string, 0),
		Search:      make(map[string]int, 0),
	}
	sprites := nodeSprites.GetNodeSprites()
	for _, sp := range sprites {
		if !strings.Contains(sp.PackageName, test) {
			relation := topology.RelationTopology(nodeSprites, sp)
			topology.Degrees[sp.GetAllPath()] = len(relation)
			topology.Relations = append(topology.Relations, relation...)
		}
	}
	return topology
}

func (topology *Topology) RelationTopology(nodeSprites *sprite.NodeSprites, node *sprite.NodeSprite) [][]*sprite.NodeSprite {
	depends := make([][]*sprite.NodeSprite, 0)
	for _, in := range node.GetImportNames() {
		sprites := nodeSprites.MatchImportNodeSprite(in)
		for _, nodeSprite := range sprites {
			if strings.Contains(nodeSprite.PackageName, test) {
				continue
			}
			if node.GetAllPath() != nodeSprite.GetAllPath() {
				depends = append(depends, []*sprite.NodeSprite{node, nodeSprite})
			}
		}
	}
	return depends
}

func (topology *Topology) CycleDepend() bool {
	queue := make([]*sprite.NodeSprite, 0)
	for k, v := range topology.Degrees {
		if v == 0 {
			queue = append(queue, topology.NodeSprites.GetNodeSprite(k))
		}
	}
	for len(queue) != 0 {
		q := queue
		queue = nil
		for _, node := range q {
			for i := 0; i < len(topology.Relations); i++ {
				if topology.Relations[i][1].GetAllPath() == node.GetAllPath() {
					topology.Degrees[topology.Relations[i][0].GetAllPath()]--
					if topology.Degrees[topology.Relations[i][0].GetAllPath()] == 0 {
						queue = append(queue, topology.Relations[i][0])
					}
				}
			}
			delete(topology.Degrees, node.FilePath+"/"+node.NodeName)
		}
	}
	return len(topology.Degrees) != 0
}

func (topology *Topology) PrintCycleDepend(hasCycle bool) {
	if !hasCycle {
		fmt.Println(Success + NotCycle)
		return
	}
	depends := map[string][]*sprite.NodeSprite{}
	for k, _ := range topology.Degrees {
		for i := 0; i < len(topology.Relations); i++ {
			if _, ok := topology.Degrees[topology.Relations[i][1].GetAllPath()]; ok &&
				topology.Relations[i][0].GetAllPath() == k {
				key := topology.Relations[i][1].GetAllPath()
				sp := topology.NodeSprites.GetNodeSprite(k)
				depends[key] = append(depends[key], sp)
			}
		}
	}
	for key, _ := range depends {
		if _, ok := topology.Search[key]; ok {
			continue
		}
		sp := topology.NodeSprites.GetNodeSprite(key)
		topology.QueryCycleNodes(sp, depends)
	}
	fmt.Printf(Failed+TextCycle, len(topology.Cycles))
	for _, cycle := range topology.Cycles {
		PrintCycle(cycle)
		fmt.Print(Line)
	}
}

func (topology *Topology) QueryCycleNodes(sp *sprite.NodeSprite, depends map[string][]*sprite.NodeSprite) {
	j := -1
	for i := 0; i < len(topology.Trace); i++ {
		if topology.Trace[i] == sp.GetAllPath() {
			j = i
		}
	}
	if j != -1 {
		circle := make([]string, 0)
		for ; j < len(topology.Trace); j++ {
			circle = append(circle, topology.Trace[j])
		}
		if len(topology.Cycles) == 0 {
			topology.Cycles = append(topology.Cycles, circle)
		} else {
			flag := false
			for _, cycle := range topology.Cycles {
				if len(cycle) == len(circle) {
					flag = true
				}
			}
			if !flag {
				topology.Cycles = append(topology.Cycles, circle)
			} else {
				flag = false
				for _, cycle := range topology.Cycles {
					if len(cycle) == len(circle) {
						cycleMap := map[string]int{}
						for _, c := range cycle {
							cycleMap[c] = 1
						}
						for _, c := range circle {
							if _, ok := cycleMap[c]; !ok {
								topology.Cycles = append(topology.Cycles, circle)
								return
							}
						}
					}
				}
			}
		}
		return
	}
	topology.Trace = append(topology.Trace, sp.GetAllPath())
	for _, v := range depends[sp.GetAllPath()] {
		topology.Search[v.GetAllPath()] = 1
		topology.QueryCycleNodes(v, depends)
	}
	topology.Trace = topology.Trace[:len(topology.Trace)-1]
}

func PrintCycle(cycle []string) {
	downMark, vertical, startPre, endPre := "↓", "┆", "┌---→", "└--- "
	output := make([]string, 0)
	length := -1
	for i := 0; i < len(cycle); i++ {
		length = Max(length, len(cycle[i]))
	}
	centre := length >> 1
	blank := fillBlank(centre + 2)
	for i := len(cycle) - 1; i >= 0; i-- {
		if i == 0 {
			output = append(output, vertical+blank+downMark)
			output = append(output, endPre+fillBlank(centre-(len(cycle[i])>>1)+2)+cycle[i])
		} else if i == len(cycle)-1 {
			output = append(output, startPre+fillBlank(centre-(len(cycle[i])>>1)+2)+cycle[i])
		} else {
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
