package day14

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type point struct {
	x int
	y int
}

type robot struct {
	p point
	v point
}

func loadFile(file string) ([]robot, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var robots []robot
	lines := strings.Split(string(content), "\n")
	for _, l := range lines {
		pv := strings.Split(strings.TrimSpace(l), " ")
		p := strings.Split(pv[0], ",")
		v := strings.Split(pv[1], ",")
		px, err := strconv.Atoi(p[0][2:])
		if err != nil {
			return nil, err
		}
		py, err := strconv.Atoi(p[1])
		if err != nil {
			return nil, err
		}
		vx, err := strconv.Atoi(v[0][2:])
		if err != nil {
			return nil, err
		}
		vy, err := strconv.Atoi(v[1])
		if err != nil {
			return nil, err
		}
		robots = append(robots, robot{
			p: point{px, py},
			v: point{vx, vy},
		})
	}
	return robots, nil
}

func (r robot) move(iter, xMax, yMax int) point {
	n := point{
		x: (r.p.x + iter*r.v.x) % xMax,
		y: (r.p.y + iter*r.v.y) % yMax,
	}
	if n.x < 0 {
		n.x += xMax
	}
	if n.y < 0 {
		n.y += yMax
	}
	return n
}

func TestMoves(t *testing.T) {
	r := robot{
		p: point{2, 4},
		v: point{2, -3},
	}
	assert.Equal(t, point{4, 1}, r.move(1, 11, 7))
	assert.Equal(t, point{6, 5}, r.move(2, 11, 7))
	assert.Equal(t, point{8, 2}, r.move(3, 11, 7))
	assert.Equal(t, point{10, 6}, r.move(4, 11, 7))
	assert.Equal(t, point{1, 3}, r.move(5, 11, 7))
}

func score(robots []robot, iter, xMax, yMax int) int {
	var q1, q2, q3, q4 int
	for _, r := range robots {
		p := r.move(iter, xMax, yMax)
		if p.x < xMax/2 {
			if p.y < yMax/2 {
				q1++
			} else if p.y > yMax/2 {
				q2++
			}
		} else if p.x > xMax/2 {
			if p.y < yMax/2 {
				q3++
			} else if p.y > yMax/2 {
				q4++
			}
		}
	}
	return q1 * q2 * q3 * q4
}

func TestDay14Part1(t *testing.T) {
	robots, err := loadFile("input_test.txt")
	require.NoError(t, err)
	assert.Equal(t, 12, score(robots, 100, 11, 7))

	robots, err = loadFile("input.txt")
	require.NoError(t, err)
	assert.Equal(t, 216772608, score(robots, 100, 101, 103))
}

var (
	g = color.RGBA{9, 0x52, 0x28, 0xff}
)

func TestDay14LineSearch(t *testing.T) {
	robots, err := loadFile("input.txt")
	require.NoError(t, err)
	var xMax, yMax = 101, 103

	points := make([]point, len(robots))
	maxMvLength := 0
	maxMv := 0
	for mv := 0; mv <= 10000; mv++ {
		for i, r := range robots {
			points[i] = r.move(mv, xMax, yMax)
		}
		sort.SliceStable(points, func(i, j int) bool {
			return points[i].x < points[j].x || (points[i].x == points[j].x && points[i].y < points[j].y)
		})
		lineSatrt := -1
		var maxLength int
		for ndx := 0; ndx < len(points)-1; ndx++ {
			if lineSatrt == -1 {
				lineSatrt = ndx
			}
			p, n := points[ndx], points[ndx+1]
			if p == n {
				continue
			}
			if p.x != n.x || p.y+1 != n.y {
				if maxLength < (p.y - points[lineSatrt].y) {
					maxLength = p.y - points[lineSatrt].y
				}
				lineSatrt = -1
			}
		}
		if maxMvLength < maxLength {
			maxMv = mv
			maxMvLength = maxLength
		}
	}
	for i, r := range robots {
		points[i] = r.move(maxMv, xMax, yMax)
	}
	sort.SliceStable(points, func(i, j int) bool {
		return points[i].y < points[j].y || (points[i].y == points[j].y && points[i].x < points[j].x)
	})
	s := ""
	curPoint := 0
	for y := 0; y < yMax; y++ {
		for x := 0; x < xMax; x++ {
			if curPoint < len(points) && points[curPoint].x == x && points[curPoint].y == y {
				s += "X"
			} else {
				s += " "
			}
			for curPoint < len(points) && points[curPoint].x == x && points[curPoint].y == y {
				curPoint++
			}
		}
		s += "\n"
	}
	print(s)
	assert.Equal(t, 6888, maxMv)
}

func TestDay14Img(t *testing.T) {
	robots, err := loadFile("input.txt")
	require.NoError(t, err)
	var xMax, yMax = 101, 103
	framesCount := 88
	a := gif.GIF{
		Image:     make([]*image.Paletted, framesCount+1),
		Delay:     make([]int, framesCount+1),
		LoopCount: 1,
	}
	for i := 0; i <= framesCount; i++ {
		img := image.NewPaletted(image.Rect(0, 0, xMax, yMax), color.Palette{color.Black, g})
		mv := 6800 + i
		for _, r := range robots {
			p := r.move(mv, xMax, yMax)
			img.Set(p.x, p.y, g)
		}
		a.Image[i] = img
		if mv == 6888 {
			f, err := os.Create("tree.gif")
			require.NoError(t, err)
			gif.Encode(f, img, nil)
			f.Close()
		}
		a.Delay[i] = 2
	}
	f, err := os.Create("robots.gif")
	require.NoError(t, err)
	gif.EncodeAll(f, &a)
	f.Close()
}
