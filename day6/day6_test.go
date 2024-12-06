package day6

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Map [][]byte

var debug = false

func (m Map) dump(pos point) {
	if !debug {
		return
	}
	fmt.Println("")
	var sb strings.Builder
	for y, l := range m {
		for x, c := range l {
			if pos.x == x && pos.y == y {
				sb.WriteRune('*')
			} else if c == North || c == South {
				sb.WriteRune('|')
			} else if c == East || c == West {
				sb.WriteRune('-')
			} else if c == 0 {
				sb.WriteRune(' ')
			} else if c < 32 {
				sb.WriteRune('+')
			} else {
				sb.WriteRune(rune(c))
			}
		}
		sb.WriteRune('\n')
	}
	fmt.Println(sb.String())
}

func (m Map) copy() Map {
	res := make(Map, len(m))
	for i, l := range m {
		res[i] = []byte(string(l))
	}
	return res
}

func loadFile(file string) (Map, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	museumMap := make(Map, len(lines))
	for i := 0; i < len(lines); i++ {
		line := strings.ReplaceAll(strings.TrimSpace(lines[i]), ".", "\x00")
		museumMap[i] = []byte(line)
	}
	return museumMap, nil
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

func findGuard(museumMap Map) point {
	yMax := len(museumMap)
	xMax := len(museumMap[0])
	for y := 0; y < yMax; y++ {
		for x := 0; x < xMax; x++ {
			if museumMap[y][x] == '^' {
				museumMap[y][x] = North
				return point{x, y}
			}
		}
	}
	return point{-1, -1}
}

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

func nextPos(pos point, dir byte, museumMap Map) (point, byte) {
	xMax := len(museumMap[0])
	yMax := len(museumMap)

	nextPos := pos.move(directions[dir])
	if nextPos.x < 0 || nextPos.x >= xMax || nextPos.y < 0 || nextPos.y >= yMax {
		return pos, 0
	}
	for museumMap[nextPos.y][nextPos.x] == '#' || museumMap[nextPos.y][nextPos.x] == 'O' {
		// turn right
		dir = (dir * 2) % 15
		nextPos = pos.move(directions[dir])
	}
	return pos.move(directions[dir]), dir
}

func patrol(pos point, museumMap Map) int {
	var dir byte = North
	visits := 1
	steps := 1
	for {
		nextPos, newDir := nextPos(pos, dir, museumMap)
		if newDir == 0 {
			break
		}
		steps++
		if (museumMap[nextPos.y][nextPos.x] & newDir) != 0 {
			return -1
		}
		if dir != newDir {
			museumMap.dump(pos)
		}
		dir = newDir
		pos = pos.move(directions[dir])
		if museumMap[pos.y][pos.x] == 0 {
			visits++
		}
		museumMap[pos.y][pos.x] = museumMap[pos.y][pos.x] + dir
	}
	return visits
}

func patrolLoops(pos point, museumMap Map) int {
	timeLoops := make(map[point]bool)
	guardStart := pos

	blockPosition := museumMap.copy()
	var dir byte = North
	for {
		nextPos, newDir := nextPos(pos, dir, blockPosition)
		if newDir == 0 {
			break
		}

		if blockPosition[nextPos.y][nextPos.x] != North {
			newMap := museumMap.copy()
			newMap[nextPos.y][nextPos.x] = 'O'
			path := patrol(guardStart, newMap)
			if path == -1 {
				newMap.dump(point{})
				timeLoops[nextPos] = true
			}
		}
		dir = newDir
		pos = pos.move(directions[dir])
		blockPosition[pos.y][pos.x] = dir
	}
	return len(timeLoops)
}

func TestDay6Part1(t *testing.T) {
	museumMap, err := loadFile("day6_test.txt")
	require.NoError(t, err)

	guardPos := findGuard(museumMap)
	visits := patrol(guardPos, museumMap)
	assert.Equal(t, 41, visits)

	museumMap, err = loadFile("day6.txt")
	require.NoError(t, err)

	guardPos = findGuard(museumMap)
	visits = patrol(guardPos, museumMap)
	assert.Equal(t, 5153, visits)
}

func TestDay6Part2(t *testing.T) {
	museumMap, err := loadFile("day6_test.txt")
	require.NoError(t, err)

	guardPos := findGuard(museumMap)
	loops := patrolLoops(guardPos, museumMap)
	assert.Equal(t, 6, loops)

	museumMap, err = loadFile("day6.txt")
	require.NoError(t, err)

	guardPos = findGuard(museumMap)
	loops = patrolLoops(guardPos, museumMap)
	assert.Equal(t, 1711, loops)
}
