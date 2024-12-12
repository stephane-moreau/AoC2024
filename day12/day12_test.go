package day12

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Map [][]byte

const (
	North = 1
	East  = 2
	South = 4
	West  = 8
)

var directions = map[byte]point{
	North: {0, -1},
	East:  {1, 0},
	South: {0, 1},
	West:  {-1, 0},
}

func (r Map) expand(start point, m Map, used map[point]bool) {
	xMax, yMax := len(m[0]), len(m)
	c := r[start.y][start.x]
	for _, d := range directions {
		p := start.move(d)
		if p.y < 0 || p.x < 0 || p.y == yMax || p.x == xMax {
			continue
		}
		if r[p.y][p.x] != 0 {
			continue
		}
		if m[p.y][p.x] == c {
			used[point{p.x, p.y}] = true
			r[p.y][p.x] = c
			r.expand(p, m, used)
		}
	}
}

func (m Map) region(start point, used map[point]bool) Map {
	xMax, yMax := len(m[0]), len(m)
	res := make(Map, len(m))
	for y := 0; y < yMax; y++ {
		res[y] = make([]byte, xMax)
	}

	c := m[start.y][start.x]
	res[start.y][start.x] = c
	for y := 0; y < yMax; y++ {
		for x := 0; x < xMax; x++ {
			if res[y][x] == c {
				used[point{x, y}] = true
				res.expand(point{x, y}, m, used)
			}
		}
	}
	return res
}

func loadFile(file string) (Map, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	fieldMap := make(Map, len(lines))
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		fieldMap[i] = []byte(line)
	}
	return fieldMap, nil
}

type point struct {
	x int
	y int
}

func (p point) move(mv point) point {
	if mv.x == 0 && mv.y == 0 {
		panic("wrong move")
	}
	return point{p.x + mv.x, p.y + mv.y}
}

func regions(m Map) []Map {
	xMax, yMax := len(m[0]), len(m)
	usedPoints := make(map[point]bool)
	var maps []Map
	for y := 0; y < yMax; y++ {
		for x := 0; x < xMax; x++ {
			if usedPoints[point{x, y}] {
				continue
			}
			r := m.region(point{x, y}, usedPoints)
			maps = append(maps, r)
		}
	}
	return maps
}

func fieldPrice(m Map) int {
	xMax, yMax := len(m[0]), len(m)
	a, p := 0, 0
	// Count fences (transition from occupied to vacancy)
	for y := 0; y < yMax; y++ {
		for x := 0; x < xMax; x++ {
			if m[y][x] != 0 {
				a++
				if y == 0 || m[y-1][x] == 0 {
					p++
				}
				if y == yMax-1 || m[y+1][x] == 0 {
					p++
				}
				if x == 0 || m[y][x-1] == 0 {
					p++
				}
				if x == xMax-1 || m[y][x+1] == 0 {
					p++
				}
			}
		}
	}
	return a * p
}

func fieldFencedPrice(m Map) int {
	xMax, yMax := len(m[0]), len(m)
	fences := make(Map, 2*len(m)+2)
	for y := 0; y < 2*yMax+2; y++ {
		fences[y] = make([]byte, 2*xMax+1)
	}
	a, p := 0, 0
	// place fences (transition from occupied to vacancy)
	for y := 0; y < yMax; y++ {
		for x := 0; x < xMax; x++ {
			if m[y][x] != 0 {
				a++
				fences[2*y+1][2*x+1] = m[y][x]
				if y == 0 {
					fences[0][2*x+1] = '-'
				}
				if y > 0 && m[y-1][x] == 0 {
					fences[2*y][2*x+1] = '-'
				}
				if y == yMax-1 || m[y+1][x] == 0 {
					fences[2*(y+1)][2*x+1] = '-'
				}
				if x == 0 {
					fences[2*y+1][0] = '|'
				}
				if x > 0 && m[y][x-1] == 0 {
					fences[2*y+1][2*x] = '|'
				}
				if x == xMax-1 || m[y][x+1] == 0 {
					fences[2*y+1][2*(x+1)] = '|'
				}
			}
		}
	}

	// place fence corners (turn in the fence)
	for y := 0; y < 2*yMax+1; y += 2 {
		for x := 0; x < 2*xMax+1; x += 2 {
			if fences[y][x] != 0 {
				panic("invalid fence")
			}
			corners := 0
			if x < 2*xMax && fences[y+1][x] != 0 && fences[y][x+1] != 0 {
				fences[y][x] = '+'
				corners++
			}
			if x > 0 && fences[y+1][x] != 0 && fences[y][x-1] != 0 {
				fences[y][x] = '+'
				corners++
			}
			if y > 0 && x < 2*xMax && fences[y-1][x] != 0 && fences[y][x+1] != 0 {
				fences[y][x] = '+'
				corners++
			}
			if y > 0 && x > 0 && fences[y-1][x] != 0 && fences[y][x-1] != 0 {
				fences[y][x] = '+'
				corners++
			}
			// a corner is either a turn or a cross betweeen 2 opposite part of the same field
			if corners == 1 {
				p++
			} else if corners == 4 {
				p += 2
			} else if corners != 0 {
				panic("wtf")
			}
		}
	}
	return a * p
}

func price(reg []Map) int {
	price := 0
	for _, r := range reg {
		price += fieldPrice(r)
	}
	return price
}

func fencePrice(reg []Map) int {
	price := 0
	for _, r := range reg {
		price += fieldFencedPrice(r)
	}
	return price
}

func TestDay12Part1(t *testing.T) {
	m, err := loadFile("mini.txt")
	require.NoError(t, err)
	regs := regions(m)
	assert.Equal(t, 140, price(regs))

	m, err = loadFile("miniO.txt")
	require.NoError(t, err)
	regs = regions(m)
	assert.Equal(t, 772, price(regs))

	m, err = loadFile("input_test.txt")
	require.NoError(t, err)
	regs = regions(m)
	assert.Equal(t, 1930, price(regs))

	m, err = loadFile("input.txt")
	require.NoError(t, err)
	regs = regions(m)
	assert.Equal(t, 1550156, price(regs))
	assert.True(t, 952948 > fencePrice(regs))
}

func TestDay12Part2(t *testing.T) {
	m, err := loadFile("mini.txt")
	require.NoError(t, err)
	regs := regions(m)
	assert.Equal(t, 80, fencePrice(regs))

	m, err = loadFile("miniO.txt")
	require.NoError(t, err)
	regs = regions(m)
	assert.Equal(t, 436, fencePrice(regs))

	m, err = loadFile("input_test.txt")
	require.NoError(t, err)
	regs = regions(m)
	assert.Equal(t, 1206, fencePrice(regs))

	m, err = loadFile("input_ex.txt")
	require.NoError(t, err)
	regs = regions(m)
	assert.Equal(t, 236, fencePrice(regs))

	m, err = loadFile("input_ab.txt")
	require.NoError(t, err)
	regs = regions(m)
	assert.Equal(t, 368, fencePrice(regs))

	m, err = loadFile("input.txt")
	require.NoError(t, err)
	regs = regions(m)
	assert.Equal(t, 946084, fencePrice(regs))
}
