package day19

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dict struct {
	entries map[string]bool
	cache   map[string]int
}

func loadInput(file string) (*regexp.Regexp, *dict, []string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, nil, nil, err
	}
	var r *regexp.Regexp
	words := dict{
		entries: map[string]bool{},
		cache:   map[string]int{},
	}
	lines := make([]string, 0, 10)
	var re string
	for i, l := range strings.Split(string(content), "\n") {
		l := strings.TrimSpace(l)
		if i == 0 {
			re = "^("
			wrds := strings.Split(l, ", ")
			sort.SliceStable(wrds, func(i int, j int) bool {
				return len(wrds[i]) < len(wrds[j])
			})
			for _, w := range wrds {
				wc := countPaths(w, &words)
				words.entries[w] = true
				words.cache[w] = 1 + wc
				re += w
			}
			re += ")*$"
			continue
		}
		if l == "" {
			continue
		}
		lines = append(lines, l)
	}
	r, err = regexp.Compile(re)
	return r, &words, lines, err
}

func TestDay19Part1(t *testing.T) {
	r, _, lines, err := loadInput("input_test.txt")
	require.NoError(t, err)
	res := 0
	for _, w := range lines {
		if r.MatchString(w) {
			res++
		}
	}
	assert.Equal(t, 6, res)

	r, _, lines, err = loadInput("input.txt")
	require.NoError(t, err)
	res = 0
	for _, w := range lines {
		if r.MatchString(w) {
			res++
		}
	}
	assert.Equal(t, 371, res)
}

func countPaths(l string, d *dict) int {
	if d.cache[l] != 0 {
		return d.cache[l]
	}
	solCount := 0
	for w := range d.entries {
		if strings.HasPrefix(l, w) {
			solCount += countPaths(l[len(w):], d)
		}
	}
	d.cache[l] = solCount
	return solCount
}

func TestDay19(t *testing.T) {
	_, words, lines, err := loadInput("input_test.txt")
	require.NoError(t, err)

	res := 0
	count := 0
	for _, l := range lines {
		solCount := countPaths(l, words)
		if solCount > 0 {
			res += solCount
			count++
		}
	}
	assert.Equal(t, 16, res)
	assert.Equal(t, 6, count)

	_, words, lines, err = loadInput("input.txt")
	require.NoError(t, err)

	res = 0
	count = 0
	for _, l := range lines {
		solCount := countPaths(l, words)
		if solCount > 0 {
			res += solCount
			count++
		}
	}
	assert.Equal(t, 650354687260341, res)
	assert.Equal(t, 371, count)
}
