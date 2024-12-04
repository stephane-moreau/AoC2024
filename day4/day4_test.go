package day4

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type grid [][]byte

func load(file string) (grid, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	g := make([][]byte, len(lines))
	for i, l := range lines {
		g[i] = []byte(strings.TrimSpace(l))
	}
	return g, nil
}

type point struct {
	x, y int
}

func (p point) move(d point) point {
	return point{p.x + d.x, p.y + d.y}
}

func (g grid) At(p point) byte {
	if p.x < 0 || p.x >= len(g[0]) {
		return 0
	}
	if p.y < 0 || p.y >= len(g) {
		return 0
	}
	return g[p.y][p.x]
}

func isXmas(g grid, p, dir point) bool {
	for _, c := range "XMAS" {
		if g.At(p) != byte(c) {
			return false
		}
		p = p.move(dir)
	}
	return true
}

var directions = []point{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}

func isXmasAllDir(g grid, p point) int {
	if g.At(p) != byte('X') {
		return 0
	}
	c := 0
	for _, d := range directions {
		if isXmas(g, p, d) {
			c++
		}
	}
	return c
}

func masInX(g grid, p point) int {
	if g.At(p) != byte('A') {
		return 0
	}
	tl := g.At(p.move(point{-1, -1}))
	tr := g.At(p.move(point{1, -1}))
	bl := g.At(p.move(point{-1, 1}))
	br := g.At(p.move(point{1, 1}))

	if tl == byte('M') && br == byte('S') &&
		((bl == byte('M') && tr == byte('S')) ||
			(tr == byte('M') && bl == byte('S'))) {
		return 1
	}
	if br == byte('M') && tl == byte('S') &&
		((bl == byte('M') && tr == byte('S')) ||
			(tr == byte('M') && bl == byte('S'))) {
		return 1
	}
	return 0
}

func count(g grid, isXmas func(g grid, p point) int) int {
	c := 0
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[0]); x++ {
			start := point{x, y}
			c += isXmas(g, start)
		}
	}
	return c
}

func TestCountXmas(t *testing.T) {
	g, err := load("day4_test.txt")
	require.NoError(t, err)
	assert.Equal(t, 18, count(g, isXmasAllDir))

	g, err = load("day4.txt")
	require.NoError(t, err)
	assert.Equal(t, 2524, count(g, isXmasAllDir))
}

func TestCountMasInX(t *testing.T) {
	g, err := load("day4_test.txt")
	require.NoError(t, err)
	assert.Equal(t, 9, count(g, masInX))

	g, err = load("day4.txt")
	require.NoError(t, err)
	assert.Equal(t, 1873, count(g, masInX))
}
