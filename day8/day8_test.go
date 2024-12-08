package day8

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type point struct {
	x int
	y int
}

func (p point) move(vx, vy int) point {
	return point{p.x + vx, p.y + vy}
}

type cityMap map[byte][]point
type CityMap struct {
	xMax, yMax int
	cm         cityMap
}

func loadCityMap(file string) (*CityMap, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	cm := make(cityMap)
	yMax := len(lines)
	xMax := len(strings.TrimSpace(lines[0]))
	for y := 0; y < yMax; y++ {
		l := strings.TrimSpace(lines[y])
		for x := 0; x < xMax; x++ {
			if l[x] >= 'a' && l[x] <= 'z' ||
				l[x] >= 'A' && l[x] <= 'Z' ||
				l[x] >= '0' && l[x] <= '9' {
				freq := l[x]
				cm[freq] = append(cm[freq], point{x, y})
			}
		}
	}
	return &CityMap{xMax, yMax, cm}, nil
}

func (cm CityMap) CreateResonance() map[point]bool {
	res := make(map[point]bool)
	for _, nw := range cm.cm {
		for _, s := range nw {
			for _, e := range nw {
				if s == e {
					continue
				}
				resX := e.x + (e.x - s.x)
				resY := e.y + (e.y - s.y)
				if resX >= 0 && resX < cm.xMax &&
					resY >= 0 && resY < cm.yMax {
					res[point{resX, resY}] = true
				}
			}
		}
	}
	return res
}

func (cm CityMap) CreateResonanceNet() map[point]bool {
	res := make(map[point]bool)
	for _, nw := range cm.cm {
		if len(nw) <= 1 {
			continue
		}
		for _, s := range nw {
			for _, e := range nw {
				if s == e {
					continue
				}
				for p := e; p.x >= 0 && p.x < cm.xMax &&
					p.y >= 0 && p.y < cm.yMax; p = p.move(e.x-s.x, e.y-s.y) {
					res[p] = true
				}
			}
		}
	}
	return res
}

func TestDay8Part1(t *testing.T) {
	cm, err := loadCityMap("day8_test.txt")
	require.NoError(t, err)
	points := cm.CreateResonance()
	assert.Equal(t, 14, len(points))

	cm, err = loadCityMap("day8.txt")
	require.NoError(t, err)
	points = cm.CreateResonance()
	assert.Equal(t, 313, len(points))
}

func TestDay8Part2(t *testing.T) {
	cm, err := loadCityMap("day8_test.txt")
	require.NoError(t, err)
	points := cm.CreateResonanceNet()
	assert.Equal(t, 34, len(points))

	cm, err = loadCityMap("day8.txt")
	require.NoError(t, err)
	points = cm.CreateResonanceNet()
	assert.Equal(t, 1064, len(points))
}
