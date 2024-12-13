package day13

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type clawValues struct {
	lines [][]int
}

var ()

func loadFile(file string, offset int) ([]clawValues, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	claws := make([]clawValues, 0, len(lines)/3)
	for i := 0; i < len(lines); i += 4 {
		l1 := strings.Split(lines[i], ",")
		l2 := strings.Split(lines[i+1], ",")
		l3 := strings.Split(lines[i+2], ",")
		if strings.TrimSpace(lines[i+3]) != "" {
			return nil, fmt.Errorf("invalid line %s - should be empty", lines[i+3])
		}
		cw := clawValues{lines: [][]int{make([]int, 3), make([]int, 3)}}

		cw.lines[0][0], err = strconv.Atoi(strings.TrimSpace(strings.Split(l1[0], "+")[1]))
		if err != nil {
			return nil, err
		}
		cw.lines[1][0], err = strconv.Atoi(strings.TrimSpace(strings.Split(l1[1], "+")[1]))
		if err != nil {
			return nil, err
		}

		cw.lines[0][1], err = strconv.Atoi(strings.TrimSpace(strings.Split(l2[0], "+")[1]))
		if err != nil {
			return nil, err
		}
		cw.lines[1][1], err = strconv.Atoi(strings.TrimSpace(strings.Split(l2[1], "+")[1]))
		if err != nil {
			return nil, err
		}

		cw.lines[0][2], err = strconv.Atoi(strings.TrimSpace(strings.Split(l3[0], "=")[1]))
		if err != nil {
			return nil, err
		}
		cw.lines[0][2] += offset
		cw.lines[1][2], err = strconv.Atoi(strings.TrimSpace(strings.Split(l3[1], "=")[1]))
		if err != nil {
			return nil, err
		}
		cw.lines[1][2] += offset
		claws = append(claws, cw)
	}
	return claws, nil
}

func (cw clawValues) det() int {
	return cw.lines[0][0]*cw.lines[1][1] - cw.lines[0][1]*cw.lines[1][0]
}

// System:
// ax + by = e
// cx + dy = f
// Solutions:
// x = (de-bf)/(ad-bc)
// y = (af-ce)/(ad-bc)
func score(claws []clawValues, limitHits bool) int {
	res := 0
	for _, cw := range claws {
		if d := cw.det(); d != 0 {
			rA := float64(cw.lines[0][2]*cw.lines[1][1]-cw.lines[1][2]*cw.lines[0][1]) / float64(d)
			rB := float64((cw.lines[1][2]*cw.lines[0][0] - cw.lines[0][2]*cw.lines[1][0])) / float64(d)
			if rA < 0 || rB < 0 {
				continue
			}
			if limitHits && (rA > 100 || rB > 100) {
				continue
			}
			a := int(rA)
			b := int(rB)
			if rA-float64(a) == 0 && rB-float64(b) == 0 {
				res += 3*a + b
			}
		}
	}
	return res
}

func TestDay13Part1(t *testing.T) {
	claws, err := loadFile("input_test.txt", 0)
	require.NoError(t, err)
	assert.Equal(t, 480, score(claws, true))

	claws, err = loadFile("input.txt", 0)
	require.NoError(t, err)
	assert.Equal(t, 26599, score(claws, true))
}

func TestDay13Part2(t *testing.T) {
	claws, err := loadFile("input_test.txt", 10_000_000_000_000)
	require.NoError(t, err)
	assert.Equal(t, 875318608908, score(claws, false))

	claws, err = loadFile("input.txt", 10_000_000_000_000)
	require.NoError(t, err)
	assert.Equal(t, 106228669504887, score(claws, false))
}
