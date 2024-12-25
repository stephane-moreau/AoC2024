package day25

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type keyOrLock []int

func loadInput(file string) ([]keyOrLock, []keyOrLock, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}
	locks := make([]keyOrLock, 0)
	keys := make([]keyOrLock, 0)
	blocks := strings.Split(strings.TrimSpace(string(content)), "\r\n\r\n")
	for _, b := range blocks {
		blockLines := strings.Split(b, "\r\n")
		lockOrKey := make([]int, len(blockLines[0]))
		for i, l := range blockLines {
			if i == 0 && strings.Trim(l, "#") == "" {
				locks = append(locks, lockOrKey)
				continue
			}
			if i == len(blockLines)-1 && strings.Trim(l, "#") == "" {
				keys = append(keys, lockOrKey)
				continue
			}
			for j, c := range l {
				if c == '#' {
					lockOrKey[j]++
				}
			}
		}
	}
	return locks, keys, err
}

func isCompatible(lock, key keyOrLock) bool {
	for i := range lock {
		if lock[i]+key[i] > 5 {
			return false
		}
	}
	return true
}
func countCompatible(locks, keys []keyOrLock) int {
	compat := 0
	for _, l := range locks {
		for _, k := range keys {
			if isCompatible(l, k) {
				compat++
			}
		}
	}
	return compat
}

func TestDay25Part1(t *testing.T) {
	locks, keys, err := loadInput("input_test.txt")
	require.NoError(t, err)
	assert.Equal(t, 2, len(locks))
	assert.Equal(t, 3, len(keys))
	assert.Equal(t, 3, countCompatible(locks, keys))

	locks, keys, err = loadInput("input.txt")
	require.NoError(t, err)
	assert.Equal(t, 3196, countCompatible(locks, keys))
}
