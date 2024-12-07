package day7

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type calibration struct {
	res     int
	numbers []int
}

func (c calibration) isValid(withConcat bool) bool {
	results := make([]int, 0)
	for i, n := range c.numbers {
		if i == 0 {
			results = append(results, n)
			continue
		}
		newRes := make([]int, 2*len(results), 3*len(results))
		for i, r := range results {
			newRes[i] = r * n
			newRes[len(results)+i] = r + n
			if withConcat {
				v, err := strconv.Atoi(fmt.Sprintf("%d%d", r, n))
				if err != nil {
					panic(err)
				}
				newRes = append(newRes, v)
			}
		}
		results = newRes
	}
	for _, r := range results {
		if r == c.res {
			return true
		}
	}
	return false
}

func loadFile(file string) ([]calibration, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	calibrations := make([]calibration, len(lines))
	for i, l := range lines {
		parts := strings.Split(strings.TrimSpace(l), ":")
		calibrations[i].res, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		nums := strings.Split(parts[1], " ")
		calibrations[i].numbers = make([]int, 0, len(nums))
		for _, n := range nums {
			if n == "" {
				continue
			}
			c, err := strconv.Atoi(n)
			if err != nil {
				return nil, err
			}
			calibrations[i].numbers = append(calibrations[i].numbers, c)
		}
	}
	return calibrations, nil
}

func result(cals []calibration, withConcat bool) int {
	s := 0
	for _, c := range cals {
		if c.isValid(withConcat) {
			s += c.res
		}
	}
	return s
}

func TestDay7Part1(t *testing.T) {
	cals, err := loadFile("day7_test.txt")
	require.NoError(t, err)
	assert.Equal(t, 3749, result(cals, false))

	cals, err = loadFile("day7.txt")
	require.NoError(t, err)
	r := result(cals, false)
	assert.NotEqual(t, 1620717387505, r) // Too high
	assert.Equal(t, 1620690235709, r)
}

func TestDay7Part3(t *testing.T) {
	cals, err := loadFile("day7_test.txt")
	require.NoError(t, err)
	assert.Equal(t, 11387, result(cals, true))

	cals, err = loadFile("day7.txt")
	require.NoError(t, err)
	r := result(cals, true)
	assert.NotEqual(t, 1620717387505, r) // Too high
	assert.Equal(t, 145397611075341, r)
}
