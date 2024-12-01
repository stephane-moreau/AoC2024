package day1

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readFile(f string) ([]int, []int, error) {
	content, err := os.ReadFile(f)
	if err != nil {
		return nil, nil, err
	}

	lines := strings.Split(string(content), "\n")
	left := make([]int, len(lines))
	right := make([]int, len(lines))
	for i, line := range lines {
		values := strings.Split(line, " ")
		left[i], err = strconv.Atoi(values[0])
		if err == nil {
			right[i], err = strconv.Atoi(strings.TrimSpace(strings.Join(values[1:], "")))
		}
		if err != nil {
			return nil, nil, err
		}
	}
	return left, right, err
}

func sumDistance(left, right []int) int {
	sort.Ints(left)
	sort.Ints(right)

	s := 0
	for i := range left {
		if left[i] > right[i] {
			s += left[i] - right[i]
		} else {
			s += right[i] - left[i]
		}
	}
	return s
}

func calcSimilarity(left, right []int) int {
	counts := make(map[int]int)
	for _, r := range right {
		counts[r] = counts[r] + 1
	}
	s := 0
	for _, l := range left {
		s += l * counts[l]
	}
	return s
}

func TestDistance(t *testing.T) {
	left, right, err := readFile("day1_test.txt")
	if err != nil {
		t.Fatal(err)
	}

	s := sumDistance(left, right)
	assert.Equal(t, 11, s)
	fmt.Printf("sum is %d\n", s)

	left, right, err = readFile("day1.txt")
	if err != nil {
		t.Fatal(err)
	}

	s = sumDistance(left, right)
	assert.Equal(t, 1666427, s)
	fmt.Printf("final sum is %d\n", s)
}

func TestSimilarity(t *testing.T) {
	left, right, err := readFile("day1_test.txt")
	if err != nil {
		t.Fatal(err)
	}

	s := calcSimilarity(left, right)
	assert.Equal(t, 31, s)
	fmt.Printf("similarity score is %d\n", s)

	left, right, err = readFile("day1.txt")
	if err != nil {
		t.Fatal(err)
	}

	s = calcSimilarity(left, right)
	assert.Equal(t, 24316233, s)
	fmt.Printf("final similarity is %d\n", s)
}
