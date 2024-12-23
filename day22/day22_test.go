package day22

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const LOOP = 2000

func loadInput(file string) ([]int, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	input := make([]int, 0, len(lines))
	for _, l := range lines {
		i, err := strconv.Atoi(strings.TrimSpace(l))
		if err != nil {
			return nil, err
		}
		input = append(input, i)
	}
	return input, err
}

func encode(code int) int {
	next := code << 6
	code = next ^ code
	code = code % 16_777_216
	next = code >> 5
	code = next ^ code
	code = code % 16_777_216
	next = code << 11
	code = next ^ code
	code = code % 16_777_216
	return code
}

func TestEncode(t *testing.T) {
	assert.Equal(t, 15887950, encode(123))
	assert.Equal(t, 16495136, encode(15887950))
	assert.Equal(t, 527345, encode(16495136))
	assert.Equal(t, 704524, encode(527345))
	assert.Equal(t, 1553684, encode(704524))
	assert.Equal(t, 12683156, encode(1553684))
	assert.Equal(t, 11100544, encode(12683156))
	assert.Equal(t, 12249484, encode(11100544))
	assert.Equal(t, 7753432, encode(12249484))
	assert.Equal(t, 5908254, encode(7753432))
}

func secretCode(input []int) int {
	res := 0
	for _, code := range input {
		for i := 0; i < LOOP; i++ {
			code = encode(code)
		}
		res += code
	}
	return res
}
func TestDay22Part1(t *testing.T) {
	input, err := loadInput("input_test.txt")
	require.NoError(t, err)
	assert.Equal(t, 37327623, secretCode(input))

	input, err = loadInput("input.txt")
	require.NoError(t, err)
	assert.Equal(t, 19927218456, secretCode(input))
}

func genPriceMap(codes []int, loop int) map[string][]int {
	var prev [4]int
	res := make(map[string][]int)
	for ndx, code := range codes {
		for i := 0; i < loop; i++ {
			newCode := encode(code)
			delta := newCode%10 - code%10
			prev[3] = prev[2]
			prev[2] = prev[1]
			prev[1] = prev[0]
			prev[0] = delta
			if i >= 3 {
				price := newCode % 10
				seq := fmt.Sprintf("%v", prev)
				_, exists := res[seq]
				if !exists {
					prices := make([]int, len(codes))
					for i := range prices {
						prices[i] = -1
					}
					res[seq] = prices
				}
				if res[seq][ndx] == -1 {
					res[seq][ndx] = price
				}
			}
			code = newCode
		}
	}
	return res
}

func sum(prices []int) int {
	s := 0
	for _, p := range prices {
		if p == -1 {
			continue
		}
		s += p
	}
	return s
}

func collectedBananas(input []int) int {
	priceMap := genPriceMap(input, 2000)
	collect := 0
	for seq, prices := range priceMap {
		s := sum(prices)
		if s > collect {
			fmt.Printf("%d: %s\n", s, seq)
			collect = s
		}
	}
	return collect
}

func TestCollect(t *testing.T) {
	deltas := genPriceMap([]int{123}, 10)
	fmt.Printf("%v\n", deltas)
}

func TestDay22Part2(t *testing.T) {
	input, err := loadInput("input_test2.txt")
	require.NoError(t, err)
	assert.Equal(t, 23, collectedBananas(input))

	input, err = loadInput("input.txt")
	require.NoError(t, err)
	// 2222 too high
	assert.Equal(t, 2189, collectedBananas(input))
}
