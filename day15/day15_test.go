package day15

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Map [][]byte
type Directions []byte

var debug = false

func (m Map) dump(force bool) {
	if !debug && !force {
		return
	}
	fmt.Println("")
	var sb strings.Builder
	for _, l := range m {
		for _, c := range l {
			sb.WriteRune(rune(c))
		}
		sb.WriteRune('\n')
	}
	fmt.Println(sb.String())
}

func findRobot(garden Map) point {
	yMax := len(garden)
	xMax := len(garden[0])
	for y := 0; y < yMax; y++ {
		for x := 0; x < xMax; x++ {
			if garden[y][x] == '@' {
				return point{x, y}
			}
		}
	}
	return point{-1, -1}
}

func loadFile(file string, scale bool) (Map, Directions, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}
	lines := strings.Split(string(content), "\n")
	garden := make(Map, 0, len(lines))
	var moves Directions
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			moves = make(Directions, 0, len(garden[0]))
		}
		if moves == nil {
			garden = append(garden, []byte(line))
		} else {
			moves = append(moves, []byte(line)...)
		}
	}
	if scale {
		for y, l := range garden {
			garden[y] = make([]byte, 2*len(l))
			for x, c := range l {
				if c == '#' || c == '.' {
					garden[y][2*x] = c
					garden[y][2*x+1] = c
				}
				if c == 'O' {
					garden[y][2*x] = '['
					garden[y][2*x+1] = ']'
				}
				if c == '@' {
					garden[y][2*x] = c
					garden[y][2*x+1] = '.'
				}
			}
		}
	}
	return garden, moves, nil
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

var dirs = map[byte]point{
	'^': {0, -1},
	'>': {1, 0},
	'v': {0, 1},
	'<': {-1, 0},
}

func last[T any](lst []T) T {
	return lst[len(lst)-1]
}

func applyMoves(garden Map, moves Directions) {
	s := findRobot(garden)
	newPos := make([]point, 0, len(garden))
	for _, m := range moves {
		newPos = append(newPos[:0], s)
		for garden[last(newPos).y][last(newPos).x] == 'O' ||
			garden[last(newPos).y][last(newPos).x] == '@' {
			newPos = append(newPos, last(newPos).move(dirs[m]))
		}
		if garden[last(newPos).y][last(newPos).x] == '#' {
			continue
		}
		for i := len(newPos) - 1; i >= 1; i-- {
			n := newPos[i]
			p := newPos[i-1]
			c := garden[n.y][n.x]
			garden[n.y][n.x] = garden[p.y][p.x]
			garden[p.y][p.x] = c
			if i == 1 {
				s = n
			}
		}
		garden.dump(false)
	}
}

func score(garden Map, m byte) int {
	s := 0
	for y, l := range garden {
		for x, c := range l {
			if c == m {
				s += 100*y + x
			}
		}
	}
	return s
}

func TestDay15Part1(t *testing.T) {
	garden, moves, err := loadFile("input_min.txt", false)
	require.NoError(t, err)
	garden.dump(false)
	applyMoves(garden, moves)
	assert.Equal(t, 2028, score(garden, 'O'))

	garden, moves, err = loadFile("input_test.txt", false)
	require.NoError(t, err)
	garden.dump(false)
	applyMoves(garden, moves)
	assert.Equal(t, 10092, score(garden, 'O'))

	garden, moves, err = loadFile("input.txt", false)
	require.NoError(t, err)
	garden.dump(false)
	applyMoves(garden, moves)
	assert.Equal(t, 1349898, score(garden, 'O'))
}

type set []point

func (s set) move(garden Map, mv point) set {
	if mv.x == 0 && mv.y == 0 {
		panic("wrong move")
	}
	newSet := make(set, 0, len(s))
	for _, p := range s {
		m := p.move(mv)
		if garden[m.y][m.x] == '.' {
			continue
		}
		newSet = append(newSet, p.move(mv))
	}
	return newSet
}

func canPush(garden Map, s set) bool {
	hasBox := false
	for _, p := range s {
		if garden[p.y][p.x] != '[' && garden[p.y][p.x] != ']' && garden[p.y][p.x] != '.' {
			return false
		}
		hasBox = hasBox || garden[p.y][p.x] == '[' || garden[p.y][p.x] == ']'
	}
	return hasBox
}

func hasHash(garden Map, s set) bool {
	for _, p := range s {
		if garden[p.y][p.x] == '#' {
			return true
		}
	}
	return false
}

func expand(garden Map, s set, m byte) set {
	if m == '<' || m == '>' || len(s) == 0 {
		return s
	}
	newSet := make(set, 0, len(s)+2)
	if garden[s[0].y][s[0].x] == ']' {
		newSet = append(newSet, point{s[0].x - 1, s[0].y})
	}
	newSet = append(newSet, s...)
	if garden[last(s).y][last(s).x] == '[' {
		newSet = append(newSet, point{last(s).x + 1, last(s).y})
	}
	return newSet
}

func applyLargeMoves(garden Map, moves Directions) {
	s := findRobot(garden)
	newPos := make([]set, 0, len(garden))
	for _, m := range moves {
		newPos = append(newPos[:0], set{s})
		for {
			newPos = append(newPos, expand(garden, last(newPos).move(garden, dirs[m]), m))
			if !canPush(garden, last(newPos)) {
				break
			}

		}
		if hasHash(garden, last(newPos)) {
			continue
		}
		for i := len(newPos) - 1; i >= 0; i-- {
			n := newPos[i]
			for _, np := range n {
				p := np.move(dirs[m])
				pC, nC := garden[p.y][p.x], garden[np.y][np.x]
				garden[p.y][p.x] = nC
				garden[np.y][np.x] = pC
			}
		}
		garden.dump(false)
		s = s.move(dirs[m])
	}
}

func TestDay15Part2(t *testing.T) {
	garden, moves, err := loadFile("input_min2.txt", true)
	require.NoError(t, err)
	garden.dump(true)
	applyLargeMoves(garden, moves)
	garden.dump(true)
	assert.Equal(t, 618, score(garden, '['))

	garden, moves, err = loadFile("input_min3.txt", true)
	require.NoError(t, err)
	garden.dump(true)
	applyLargeMoves(garden, moves)
	garden.dump(true)
	assert.Equal(t, 413, score(garden, '['))

	garden, moves, err = loadFile("input_test.txt", true)
	require.NoError(t, err)
	garden.dump(true)
	applyLargeMoves(garden, moves)
	garden.dump(true)
	assert.Equal(t, 9021, score(garden, '['))

	garden, moves, err = loadFile("input.txt", true)
	require.NoError(t, err)
	garden.dump(true)
	applyLargeMoves(garden, moves)
	assert.Equal(t, 1376686, score(garden, '['))
}
