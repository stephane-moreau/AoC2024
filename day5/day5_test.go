package day4

import (
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadFile(file string) (map[int]map[int]bool, [][]int, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}
	lines := strings.Split(string(content), "\n")
	rules := make(map[int]map[int]bool)
	orders := make([][]int, 0, 10)
	var i int
	for i = 0; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])
		if l == "" {
			break
		}
		vals := strings.Split(l, "|")
		first, err := strconv.Atoi(vals[0])
		if err != nil {
			return nil, nil, err
		}
		next, err := strconv.Atoi(vals[1])
		if err != nil {
			return nil, nil, err
		}
		if _, exist := rules[first]; !exist {
			rules[first] = make(map[int]bool)
		}
		rules[first][next] = true
	}
	for ; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])
		if l == "" {
			continue
		}
		vals := strings.Split(l, ",")
		order := make([]int, len(vals))
		for j, v := range vals {
			order[j], err = strconv.Atoi(v)
			if err != nil {
				return nil, nil, err
			}
		}
		orders = append(orders, order)
	}
	return rules, orders, nil
}

func errorNdx(order []int, rules map[int]map[int]bool) int {
	for i := 0; i < len(order)-1; i++ {
		if !rules[order[i]][order[i+1]] {
			return i
		}
	}
	return -1
}

func filterOrders(orders [][]int, rules map[int]map[int]bool) ([][]int, int) {
	res := make([][]int, 0)
	val := 0
	for _, order := range orders {
		if errorNdx(order, rules) == -1 {
			res = append(res, order)
			val += order[len(order)/2]
		}
	}
	return res, val
}

func TestDay5Part1(t *testing.T) {
	rules, orders, err := loadFile("day5_test.txt")
	require.NoError(t, err)

	valids, count := filterOrders(orders, rules)
	assert.Equal(t, 3, len(valids))
	assert.Equal(t, 143, count)

	rules, orders, err = loadFile("day5.txt")
	require.NoError(t, err)
	valids, count = filterOrders(orders, rules)
	assert.Equal(t, 106, len(valids))
	assert.Equal(t, 6034, count)
}

func repairOrders(orders [][]int, rules map[int]map[int]bool) ([][]int, int) {
	res := make([][]int, 0)
	val := 0
	for _, order := range orders {
		errNdx := errorNdx(order, rules)
		if errNdx == -1 {
			continue
		}

		slices.SortFunc(order, func(a, b int) int {
			if rules[a][b] {
				return -1
			}
			if rules[b][a] {
				return 1
			}
			return 0
		})
		res = append(res, order)
		val += order[len(order)/2]
	}
	return res, val
}

func TestDay5Part2(t *testing.T) {
	rules, orders, err := loadFile("day5_test.txt")
	require.NoError(t, err)

	valids, count := repairOrders(orders, rules)
	assert.Equal(t, 3, len(valids))
	assert.Equal(t, 123, count)

	rules, orders, err = loadFile("day5.txt")
	require.NoError(t, err)
	valids, count = repairOrders(orders, rules)
	assert.Equal(t, 111, len(valids))
	assert.Equal(t, 6305, count)
}
