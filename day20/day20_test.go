package day20

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Map [][]byte

func (m Map) val(p point) byte {
	return m[p.y][p.x]
}

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

var directions = map[byte]point{
	'^': {0, -1},
	'<': {-1, 0},
	'>': {1, 0},
	'v': {0, 1},
}

var zero point

func (p point) move(mv point) point {
	if mv == zero {
		panic("wrong move")
	}
	return point{p.x + mv.x, p.y + mv.y}
}

func findPosition(m Map, c byte) point {
	xMax, yMax := len(m[0]), len(m)
	for y := 0; y < yMax; y++ {
		for x := 0; x < xMax; x++ {
			p := point{x, y}
			if m.val(p) == c {
				return p
			}
		}
	}
	return point{}
}

func findPath(m Map) map[point]int {
	s := findPosition(m, 'S')
	path := map[point]int{s: 0}
	cur := s
	l := 0
	for m[cur.y][cur.x] != 'E' {
		for _, d := range directions {
			n := cur.move(d)
			if (m.val(n) == '.' || m.val(n) == 'E') && path[n] == 0 {
				cur = n
				l++
				path[n] = l
				break
			}
		}
	}
	return path
}

type position struct {
	point
	d byte
}

func abs(i int) int {
	if i >= 0 {
		return i
	}
	return -i
}

func explore(path map[point]int, warpSize int) map[int]int {
	cheats := map[int]int{}
	for s, l := range path {
		for t, n := range path {
			if n <= l {
				continue
			}
			c := abs(t.x-s.x) + abs(t.y-s.y)
			if c > warpSize {
				continue
			}
			if (n - l) > c {
				cheats[n-l-c]++
			}
		}
	}
	return cheats
}

func TestDay20(t *testing.T) {
	m, err := loadFile("input_test.txt")
	require.NoError(t, err)

	path := findPath(m)
	cheats := explore(path, 2)
	assert.Equal(t, 11, len(cheats))
	assert.Equal(t, 14, cheats[2])
	assert.Equal(t, 3, cheats[12])
	assert.Equal(t, 1, cheats[64])
	cheats = explore(path, 20)
	assert.Equal(t, 39, cheats[56])
	assert.Equal(t, 19, cheats[64])
	assert.Equal(t, 22, cheats[72])
	assert.Equal(t, 3, cheats[76])

	m, err = loadFile("input.txt")
	require.NoError(t, err)

	path = findPath(m)
	cheats = explore(path, 2)
	res := 0
	for g, cheat := range cheats {
		if g >= 100 {
			res += cheat
		}
	}
	assert.Equal(t, 1351, res)

	cheats = explore(path, 20)
	res = 0
	for g, cheat := range cheats {
		if g >= 100 {
			res += cheat
		}
	}
	assert.Equal(t, 966130, res)
}
