package day10

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

type testCase struct {
	fileName         string
	trailHeadScores  []int
	trailHeadRatings []int
	score            int
	rating           int
}

var cases = []testCase{
	{fileName: "mini_input.txt", trailHeadScores: []int{1}, score: 1},
	{fileName: "input1_test.txt", trailHeadScores: []int{2}, score: 2},
	{fileName: "input2_test.txt", trailHeadScores: []int{4}, score: 4},
	{fileName: "input3_test.txt", trailHeadScores: []int{1, 2}, score: 3},
	{fileName: "input4_test.txt", trailHeadScores: []int{5, 6, 5, 3, 1, 3, 5, 3, 5},
		trailHeadRatings: []int{20, 24, 10, 4, 1, 4, 5, 8, 5}, score: 36, rating: 81},
	{"input.txt", nil, nil, 688, 1459},
}

func (tc testCase) test(t *testing.T) {
	m, err := loadFile(tc.fileName)
	require.NoError(t, err)
	zeros := make([]point, 0)
	xMax, yMax := len(m[0]), len(m)
	for y := 0; y < yMax; y++ {
		for x := 0; x < xMax; x++ {
			if m[y][x] == '0' {
				zeros = append(zeros, point{x, y})
			}
		}
	}
	if tc.trailHeadScores != nil {
		assert.Equal(t, len(tc.trailHeadScores), len(zeros), "error in '%s'", tc.fileName)
	}
	s := 0
	r := 0
	for i, z := range zeros {
		pos := []point{z}
		var c byte
		for c = '0'; c <= '9'; c++ {
			newPos := make([]point, 0, len(pos))
			for _, p := range pos {
				for _, d := range directions {
					n := p.move(d)
					if n.x < 0 || n.x >= xMax {
						continue
					}
					if n.y < 0 || n.y >= yMax {
						continue
					}
					if m[n.y][n.x] == c+1 {
						newPos = append(newPos, n)
					}
				}
			}
			if len(newPos) == 0 {
				break
			}
			pos = newPos
		}
		if c == '9' {
			dedup := make(map[point]bool)
			for _, p := range pos {
				dedup[p] = true
			}
			if tc.trailHeadScores != nil {
				assert.Equal(t, tc.trailHeadScores[i], len(dedup), "error in '%s'", tc.fileName)
			}
			if tc.trailHeadRatings != nil {
				assert.Equal(t, tc.trailHeadRatings[i], len(pos), "error in '%s'", tc.fileName)
			}
			s += len(dedup)
			r += len(pos)
		}
	}
	assert.Equal(t, tc.score, s, "error in '%s'", tc.fileName)
	if tc.rating > 0 {
		assert.Equal(t, tc.rating, r, "error in '%s'", tc.fileName)
	}
}

func TestDay10(t *testing.T) {
	for _, tc := range cases {
		tc.test(t)
	}
}
