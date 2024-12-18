package day18

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Map [][]byte

func loadFile(file string, xMax, yMax, size int) (Map, []point, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}
	lines := strings.Split(string(content), "\n")
	m := make(Map, yMax+3)
	for y := 0; y < yMax+3; y++ {
		m[y] = make([]byte, xMax+3)
		if y == 0 || y == yMax+2 {
			for x := range m[y] {
				m[y][x] = '#'
			}
		}
		m[y][0] = '#'
		m[y][xMax+2] = '#'
	}
	remaining := make([]point, 0, len(lines))
	for i, l := range lines {
		coords := strings.Split(strings.TrimSpace(l), ",")
		x, err := strconv.Atoi(coords[0])
		if err != nil {
			return nil, nil, err
		}
		y, err := strconv.Atoi(coords[1])
		if err != nil {
			return nil, nil, err
		}
		if i <= size-1 {
			m[y+1][x+1] = '#'
		} else {
			remaining = append(remaining, point{x + 1, y + 1})
		}
	}
	m[1][1] = 'S'
	m[yMax+1][xMax+1] = 'E'
	return m, remaining, nil
}

var debug = true

func (m Map) dump() {
	if !debug {
		return
	}
	fmt.Println("")
	var sb strings.Builder
	for y, l := range m {
		for x := range l {
			if m[y][x] != 0 {
				sb.WriteRune(rune(m[y][x]))
			} else {
				sb.WriteRune('.')
			}
		}
		sb.WriteRune('\n')
	}
	fmt.Println(sb.String())
}

type point struct {
	x, y int
}

var zero point

var directions = map[byte]point{
	'^': {0, -1},
	'<': {-1, 0},
	'>': {1, 0},
	'v': {0, 1},
}

func (p point) move(mv point) point {
	if mv == zero {
		panic("wrong move")
	}
	return point{p.x + mv.x, p.y + mv.y}
}

func explore(p point, path map[point]bool, m Map, visited map[point]int, curCost *int) {
DFS:
	for _, mv := range directions {
		n := p.move(mv)
		if m[n.y][n.x] == '#' {
			// wall
			continue
		}
		if visited[n] > 0 && visited[n] < len(path) {
			continue
		}
		visited[n] = len(path)
		if m[n.y][n.x] == 'E' {
			// wall
			path[n] = true
			c := len(path)
			if *curCost == 0 || *curCost > c {
				*curCost = c
			}
			return
		}
		// loop
		if path[n] {
			continue DFS
		}
		path[n] = true
		explore(n, path, m, visited, curCost)
		delete(path, n)
	}
}

func shortestMapPath(m Map, xMax, yMax int) int {
	area := xMax * yMax
	hash := area + 1

	visited := make([][]int, len(m))
	for y := range visited {
		visited[y] = make([]int, len(m[0]))
		for x := range visited[y] {
			if m[y][x] == '#' {
				visited[y][x] = hash // '#'
				continue
			}
			visited[y][x] = area
		}
	}
	visited[1][1] = 1
	dedup := map[point]int{{1, 1}: 1}
	for visited[yMax+1][xMax+1] == area {
		d := len(dedup)
		for p := range dedup {
			y, x := p.y, p.x
			if visited[y][x+1] != hash {
				nx := x + 1
				ny := y
				visited[ny][nx] = min(visited[ny][nx], visited[ny][nx-1]+1, visited[ny][nx+1]+1, visited[ny+1][nx]+1, visited[ny-1][nx]+1)
				dedup[point{x + 1, y}] = visited[y][x+1]
			}
			if visited[y][x-1] != hash {
				nx := x - 1
				ny := y
				visited[ny][nx] = min(visited[ny][nx], visited[ny][nx-1]+1, visited[ny][nx+1]+1, visited[ny+1][nx]+1, visited[ny-1][nx]+1)
				dedup[point{x - 1, y}] = visited[y][x-1]
			}
			if visited[y+1][x] != hash {
				nx := x
				ny := y + 1
				visited[ny][nx] = min(visited[ny][nx], visited[ny][nx-1]+1, visited[ny][nx+1]+1, visited[ny+1][nx]+1, visited[ny-1][nx]+1)
				dedup[point{x, y + 1}] = visited[y+1][x]
			}
			if visited[y-1][x] != hash {
				nx := x
				ny := y - 1
				visited[ny][nx] = min(visited[ny][nx], visited[ny][nx-1]+1, visited[ny][nx+1]+1, visited[ny+1][nx]+1, visited[ny-1][nx]+1)
				dedup[point{x, y - 1}] = visited[y][x-1]
			}
		}
		if d == len(dedup) {
			break
		}
	}
	if visited[yMax+1][xMax+1] == area {
		return -1
	}
	return visited[yMax+1][xMax+1] - 1
}

func TestDay18Part1(t *testing.T) {
	m, _, err := loadFile("input_test.txt", 6, 6, 12)
	require.NoError(t, err)
	m.dump()
	assert.Equal(t, 22, shortestMapPath(m, 6, 6))

	m, _, err = loadFile("input.txt", 70, 70, 1024)
	require.NoError(t, err)
	m.dump()
	assert.Equal(t, 408, shortestMapPath(m, 70, 70))
}

func TestDay18Part2(t *testing.T) {
	m, remaining, err := loadFile("input_test.txt", 6, 6, 12)
	require.NoError(t, err)
	for i := 0; i < len(remaining); i++ {
		b := remaining[i]
		m[b.y][b.x] = '#'
		if shortestMapPath(m, 6, 6) == -1 {
			// offset by 1 for result
			assert.Equal(t, point{7, 2}, b)
			break
		}
	}

	m, remaining, err = loadFile("input.txt", 70, 70, 1024)
	require.NoError(t, err)
	for i := 0; i < len(remaining); i++ {
		b := remaining[i]
		m[b.y][b.x] = '#'
		if shortestMapPath(m, 70, 70) == -1 {
			// offset by 1 for result
			assert.Equal(t, point{46, 17}, b)
			break
		}
		if i%100 == 0 {
			fmt.Printf(".")
		}
	}
}
