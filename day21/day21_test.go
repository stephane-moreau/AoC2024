package day21

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	input_test = []string{
		"029A",
		"980A",
		"179A",
		"456A",
		"379A",
	}
	input = []string{
		"286A",
		"480A",
		"140A",
		"413A",
		"964A",
	}

	numPad = []string{
		"789",
		"456",
		"123",
		" 0A",
	}

	keyPad = []string{
		" ^A",
		"<v>",
	}
)

type point struct {
	x, y int
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

func isValid(code string, pad []string, p point) bool {
	for _, c := range code {
		if pad[p.y][p.x] == ' ' {
			return false
		}
		p = p.move(directions[byte(c)])
	}
	return true
}

type path map[int][]string

func mvId(s, e byte) int {
	// return fmt.Sprintf("%c%c", s, e)
	return (int(s) << 8) | int(e)
}

func paths(pad []string) path {
	res := make(path)
	for r1 := 0; r1 < len(pad); r1++ {
		for c1 := 0; c1 < len(pad[r1]); c1++ {
			if pad[r1][c1] == ' ' {
				continue
			}
			for r2 := 0; r2 < len(pad); r2++ {
				for c2 := 0; c2 < len(pad[r2]); c2++ {
					if pad[r2][c2] == ' ' {
						continue
					}
					mv := mvId(pad[r1][c1], pad[r2][c2])
					if res[mv] != nil {
						continue
					}
					var sH, sV string
					if c1 < c2 {
						for i := c1; i < c2; i++ {
							sH += ">"
						}
					} else {
						for i := c1; i > c2; i-- {
							sH += "<"
						}
					}
					if r1 < r2 {
						for i := r1; i < r2; i++ {
							sV += "v"
						}
					} else {
						for i := r1; i > r2; i-- {
							sV += "^"
						}
					}
					if sV == "" {
						res[mv] = append(res[mv], sH+"A")
					} else if sH == "" {
						res[mv] = append(res[mv], sV+"A")
					} else {
						if isValid(sV+sH, pad, point{c1, r1}) {
							res[mv] = append(res[mv], sV+sH+"A")
						}
						if isValid(sH+sV, pad, point{c1, r1}) {
							res[mv] = append(res[mv], sH+sV+"A")
						}
					}
				}
			}
		}
	}
	return res
}

type encoding []string

func encode(code string, paths path) []encoding {
	cur := byte('A')
	travels := make([]encoding, 0, 256)
	for i := range code {
		c := code[i]
		travels = append(travels, paths[mvId(cur, c)])
		cur = c
	}
	return travels
}

type entry struct {
	lvl  int
	path string
}

func shortestPath(part string, keyPaths path, limit int, cache map[entry]int) int {
	cur := byte('A')
	sp := 0
	ntr := entry{limit, part}
	cachedSP := cache[ntr]
	if cachedSP != 0 {
		return cachedSP
	}
	for i := range part {
		c := part[i]
		mvPaths := keyPaths[mvId(cur, c)]
		shortestLen := 0
		if limit > 1 {
			for _, p := range mvPaths {
				newPart := shortestPath(p, keyPaths, limit-1, cache)
				if shortestLen == 0 || shortestLen > newPart {
					shortestLen = newPart
				}
			}
		} else {
			shortestLen = len(mvPaths[0])
		}
		sp += shortestLen
		cur = c
	}
	cache[ntr] = sp
	return sp
}

func scoreWithCache(codes []string, numPath, keyPath path, limit int) int {
	//
	sc := 0
	cache := make(map[entry]int)

	for _, code := range codes {
		val, _ := strconv.ParseInt(strings.TrimRight(strings.TrimLeft(code, "0"), "A"), 10, 64)
		fmt.Printf("%v\n", val)
		paths := encode(code, numPath)
		pathLen := 0
		for _, parts := range paths {
			partLen := 0
			for _, part := range parts {
				e := entry{limit, part}
				if cache[e] == 0 {
					cache[e] = shortestPath(part, keyPath, limit, cache)
				}
				if partLen == 0 || partLen > cache[e] {
					partLen = cache[e]
				}
			}
			pathLen += partLen
		}
		//fmt.Printf("%v\n", path)
		sc += pathLen * int(val)
	}
	return sc
}

func TestDay21Encode(t *testing.T) {
	numPath := paths(numPad)
	keyPath := paths(keyPad)
	assert.Equal(t, 19448, scoreWithCache([]string{"286A"}, numPath, keyPath, 2))
	assert.Equal(t, 299728, scoreWithCache([]string{"286A"}, numPath, keyPath, 5))
	assert.Equal(t, 28603432, scoreWithCache([]string{"286A"}, numPath, keyPath, 10))
	assert.Equal(t, 2725295144, scoreWithCache([]string{"286A"}, numPath, keyPath, 15))
	assert.Equal(t, 259618624408, scoreWithCache([]string{"286A"}, numPath, keyPath, 20))
	assert.Equal(t, 24732073940288, scoreWithCache([]string{"286A"}, numPath, keyPath, 25))
}

func TestDay21(t *testing.T) {
	numPath := paths(numPad)
	keyPath := paths(keyPad)

	assert.Equal(t, 126384, scoreWithCache(input_test, numPath, keyPath, 2))

	assert.Equal(t, 163086, scoreWithCache(input, numPath, keyPath, 2))
	assert.Equal(t, 198466286401228, scoreWithCache(input, numPath, keyPath, 25))
}
