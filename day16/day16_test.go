package day16

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type edge struct {
	target point
	cost   int
	dir    byte
}
type node struct {
	pos     point
	targets []edge
}
type graph map[point]*node

type Map [][]byte

func loadFile(file string) (Map, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	m := make(Map, len(lines))
	for i, l := range lines {
		m[i] = []byte(strings.TrimSpace(l))
	}
	return m, nil
}

type point struct {
	x, y int
}
type position struct {
	point
	d byte
}

var zero point

func (p point) move(mv point) point {
	if mv == zero {
		panic("wrong move")
	}
	return point{p.x + mv.x, p.y + mv.y}
}

var directions = map[byte]point{
	'^': {0, -1},
	'<': {-1, 0},
	'>': {1, 0},
	'v': {0, 1},
}

func findPosition(m Map, c byte) point {
	xMax, yMax := len(m[0]), len(m)
	for y := 0; y < yMax; y++ {
		for x := 0; x < xMax; x++ {
			if m[y][x] == c {
				return point{x, y}
			}
		}
	}
	return point{}
}

func cost(path map[point][]byte) int {
	c := 0
	for _, d := range path {
		if len(d) == 2 {
			c += 1001
		} else {
			c++
		}
	}
	return c
}

func explore(p position, path map[point][]byte, m Map, curCost *int) {
	curMv := directions[p.d]
DFS:
	for d, mv := range directions {
		n := p.move(mv)
		if m[n.y][n.x] == '#' {
			// wall
			continue
		}
		if curMv.move(mv) == zero {
			// U turn
			continue
		}
		if m[n.y][n.x] == 'E' {
			// wall
			path[n] = []byte{d}
			c := cost(path)
			if *curCost == 0 || *curCost > c {
				*curCost = c
			}
			return
		}
		// loop
		if path[n] != nil {
			continue DFS
		}
		if d != p.d {
			path[p.point] = append(path[p.point], d)
		}
		path[n] = []byte{d}
		explore(position{n, d}, path, m, curCost)
		delete(path, n)
		if d != p.d {
			path[p.point] = path[p.point][:1]
		}
	}
}

func abs(i int) int {
	if i >= 0 {
		return i
	}
	return -i
}

func shortestMapPath(m Map) int {
	s := findPosition(m, 'S')
	path := map[point][]byte{s: {'>'}}
	var c int
	explore(position{s, '>'}, path, m, &c)
	return c - 1
}

func digest(m Map) graph {
	s := findPosition(m, 'S')
	nodes := make([]*node, 0, 1000)
	nodes = append(nodes, &node{pos: s})
	g := graph{s: nodes[0]}
	for i := 0; i < len(nodes); i++ {
		nd := nodes[i]
		p := nd.pos
		for dir, d := range directions {
			n := p.move(d)
			if m[n.y][n.x] == '#' {
				continue
			}

			for {
				up := n.move(point{d.y, d.x})
				down := n.move(point{-d.y, -d.x})
				if m[n.y][n.x] == '#' {
					break
				}
				if m[up.y][up.x] != '#' || m[down.y][down.x] != '#' {
					break
				}
				n = n.move(d)
			}
			if existing := g[n]; existing != nil {
				newEdge := true
				// for _, t := range existing.targets {
				// 	if t.target == nd.pos {
				// 		newEdge = false
				// 	}
				// }
				if newEdge {
					nd.targets = append(nd.targets, edge{
						target: n,
						cost:   abs(n.x - nd.pos.x + n.y - nd.pos.y),
						dir:    dir,
					})
				}
				continue
			}

			if m[n.y][n.x] == '#' {
				newNode := &node{pos: point{n.x - d.x, n.y - d.y}}
				nd.targets = append(nd.targets, edge{
					target: newNode.pos,
					cost:   abs(newNode.pos.x - nd.pos.x + newNode.pos.y - nd.pos.y),
					dir:    dir,
				})
				nodes = append(nodes, newNode)
				g[newNode.pos] = newNode
				continue
			}
			up := n.move(point{d.y, d.x})
			down := n.move(point{-d.y, -d.x})
			if m[up.y][up.x] != '#' || m[down.y][down.x] != '#' {
				newNode := &node{pos: n}
				nd.targets = append(nd.targets, edge{
					target: n,
					cost:   abs(n.x - nd.pos.x + n.y - nd.pos.y),
					dir:    dir,
				})
				nodes = append(nodes, newNode)
				g[n] = newNode
			}
		}
	}

	return g
}

func display(path map[point]edge) string {
	var disp string
	cur := path[point{-1, -1}]
	for i := 0; i < len(path)-1; i++ {
		cur = path[cur.target]
		disp += fmt.Sprintf(" %c(%d) %v", cur.dir, cur.cost, cur.target)
	}
	return disp
}

func shortestGraphPath(g graph, m Map) int {
	var res int
	pos := findPosition(m, 'S')
	path := map[point]edge{
		{-1, -1}: {pos, 0, '>'},
	}
	var curDir byte = '>'
	cost := &res
	visited := map[edge]int{{pos, 0, '>'}: 0}
	return internalShortestTravel(pos, curDir, 0, visited, path, g, m, cost, nil)
}

func internalShortestTravel(pos point, curDir byte, curCost int,
	visited map[edge]int, path map[point]edge,
	g graph, m Map, cost *int, usedEdges map[edge]int) int {
	curNd := g[pos]
	for _, t := range curNd.targets {
		if _, exists := path[t.target]; exists {
			continue
		}
		delta := t.cost
		if curDir != t.dir {
			delta += 1000
		}
		if visited[t] != 0 && visited[t] < curCost+delta {
			continue
		}
		visited[t] = curCost + delta
		path[curNd.pos] = t
		if m[t.target.y][t.target.x] == 'E' {
			if *cost == 0 || curCost+delta <= *cost {
				if usedEdges != nil {
					for _, e := range path {
						usedEdges[e] = usedEdges[e] + 1
					}
				}
				fmt.Printf("%d : %s\n", curCost+delta, display(path))
				*cost = curCost + delta
			}
		}
		if *cost == 0 || curCost+delta < *cost {
			internalShortestTravel(t.target, t.dir, curCost+delta, visited, path, g, m, cost, usedEdges)
		}
		delete(path, curNd.pos)
	}
	return *cost
}

func TestDay16Part1(t *testing.T) {
	m, err := loadFile("input_test1.txt")
	require.NoError(t, err)
	assert.Equal(t, 7036, shortestMapPath(m))

	g := digest(m)
	st := shortestGraphPath(g, m)
	assert.Equal(t, 7036, st)

	m, err = loadFile("input_test2.txt")
	require.NoError(t, err)
	assert.Equal(t, 11048, shortestMapPath(m))

	g = digest(m)
	st = shortestGraphPath(g, m)
	assert.Equal(t, 11048, st)

	m, err = loadFile("input.txt")
	require.NoError(t, err)
	g = digest(m)
	st = shortestGraphPath(g, m)
	assert.Equal(t, 107512, st)
}

func countvisitedSpots(g graph, m Map, minLength int) int {
	usedEdges := make(map[edge]int)
	pos := findPosition(m, 'S')
	path := map[point]edge{
		{-1, -1}: {pos, 0, '>'},
	}
	var curDir byte = '>'
	cost := &minLength
	visited := map[edge]int{{pos, 0, '>'}: 0}

	internalShortestTravel(pos, curDir, 0, visited, path, g, m, cost, usedEdges)
	res := map[point]bool{}
	for e := range usedEdges {
		res[e.target] = true
		d := directions[e.dir]
		for i := 0; i < e.cost; i++ {
			res[point{e.target.x - i*d.x, e.target.y - i*d.y}] = true
		}
	}
	return len(res)
}

func TestDay16Part2(t *testing.T) {
	m, err := loadFile("input_test1.txt")
	require.NoError(t, err)
	g := digest(m)
	s := countvisitedSpots(g, m, 7036)
	assert.Equal(t, 45, s)

	m, err = loadFile("input_test2.txt")
	require.NoError(t, err)
	g = digest(m)
	s = countvisitedSpots(g, m, 11048)
	assert.Equal(t, 64, s)

	m, err = loadFile("input.txt")
	require.NoError(t, err)
	g = digest(m)
	s = countvisitedSpots(g, m, 107512)
	assert.Equal(t, 561, s)
}
