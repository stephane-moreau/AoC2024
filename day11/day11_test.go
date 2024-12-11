package day11

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	input_test = []int{125, 17}
	input      = []int{0, 5601550, 3914, 852, 50706, 68, 6, 645371}
)

func blink(stone int) []int {
	if stone == 0 {
		return []int{1}
	}

	c := int(math.Log10(float64(stone))) + 1
	if c%2 == 0 {
		n := int(math.Pow10(c / 2))
		i1 := stone / n
		i2 := stone % n
		return []int{i1, i2}
	}
	return []int{stone * 2024}
}

type splitValues struct {
	stone      int
	subStones  []int
	stepCounts map[int]int
}

func numBlinks(stone int, blinks map[int]*splitValues, maxSteps int) int {
	for l := 0; l < maxSteps; l++ {
	stones:
		for _, sv := range blinks {
			level := len(sv.stepCounts)
			if level == 0 {
				level = 1
			}
			n := 0
			for _, ss := range sv.subStones {
				bss := blinks[ss]
				if bss == nil {
					continue stones
				}
				n += blinks[ss].stepCounts[level]
			}
			sv.stepCounts[level+1] = n
		}
	}
	res := 0
	for _, s := range blinks[stone].subStones {
		res += blinks[s].stepCounts[maxSteps]
	}
	return res
}

func processBlinks(stones []int, maxSteps int) map[int]*splitValues {
	blinks := map[int]*splitValues{
		-1: {
			stone:      -1,
			subStones:  stones,
			stepCounts: map[int]int{},
		},
	}
	for i := 0; i < maxSteps; i++ {
		cur := make(map[int]bool)
		for _, sv := range blinks {
			for _, s := range sv.subStones {
				if blinks[s] != nil || cur[s] {
					continue
				}
				cur[s] = true
			}
		}
		for s := range cur {
			ss := blink(s)
			blinks[s] = &splitValues{
				stone:      s,
				subStones:  ss,
				stepCounts: map[int]int{1: len(ss)},
			}
		}
	}
	return blinks
}

func TestDay11Part1(t *testing.T) {
	blinksRes := processBlinks(input_test, 6)
	assert.Equal(t, 3, numBlinks(-1, blinksRes, 1))
	assert.Equal(t, 4, numBlinks(-1, blinksRes, 2))
	assert.Equal(t, 5, numBlinks(-1, blinksRes, 3))
	assert.Equal(t, 9, numBlinks(-1, blinksRes, 4))
	assert.Equal(t, 13, numBlinks(-1, blinksRes, 5))
	assert.Equal(t, 22, numBlinks(-1, blinksRes, 6))
	blinksRes = processBlinks(input_test, 25)
	assert.Equal(t, 55312, numBlinks(-1, blinksRes, 25))
	blinksRes = processBlinks(input, 25)
	assert.Equal(t, 189092, numBlinks(-1, blinksRes, 25))
	blinksRes = processBlinks(input, 75)
	assert.Equal(t, 189092, numBlinks(-1, blinksRes, 25))
	assert.Equal(t, 224869647102559, numBlinks(-1, blinksRes, 75))
}
