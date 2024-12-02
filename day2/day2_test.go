package day2

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type path []int

func invalidPos(p path) int {
	var prev int
	sign := p[1] - p[0]

	invalidPos := -1
	log := strconv.Itoa(p[0])
	for i, cur := range p {
		if i == 0 {
			prev = cur
			continue
		}
		d := cur - prev
		prev = cur
		if d == 0 ||
			(sign > 0 && (d > 3 || d < 0)) ||
			(sign < 0 && (d < -3 || d > 0)) {
			if invalidPos == -1 {
				invalidPos = i
			}
			log += " * "
		} else {
			log += "   "
		}
		log += strconv.Itoa(cur)
	}
	if invalidPos != -1 {
		fmt.Printf("invalid %s\n", log)
	}
	return invalidPos
}

func readFile(f string) ([]path, error) {
	content, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	paths := make([]path, len(lines))
	for i, line := range lines {
		values := strings.Split(line, " ")
		path := make([]int, len(values))
		for j, s := range values {
			path[j], err = strconv.Atoi(strings.TrimSpace(s))
			if err != nil {
				return nil, err
			}
		}
		paths[i] = path
	}
	return paths, err
}

func countValidPaths(paths []path) int {
	s := 0
	for _, path := range paths {
		if invalidPos(path) == -1 {
			s++
		}
	}
	return s
}

func countValidTolarablePaths(paths []path) int {
	s := 0
	for _, p := range paths {
		pos := invalidPos(p)
		if pos == -1 {
			s++
		} else if invalidPos(append(append(path{}, p[:pos]...), p[pos+1:]...)) == -1 {
			s++
		} else if invalidPos(append(append(path{}, p[:pos-1]...), p[pos:]...)) == -1 {
			s++
		} else if invalidPos(append(path{}, p[1:]...)) == -1 {
			s++
		} else {
			fmt.Printf("invalid Path at %d: %v\n", pos, p)
		}
	}
	return s
}

func TestValid(t *testing.T) {
	paths, err := readFile("day2_test.txt")
	if err != nil {
		t.Fatal(err)
	}

	c := countValidPaths(paths)
	assert.Equal(t, 2, c)
	fmt.Printf("number of valid path is %d\n", c)

	paths, err = readFile("day2.txt")
	if err != nil {
		t.Fatal(err)
	}

	c = countValidPaths(paths)
	assert.Equal(t, 502, c)
	fmt.Printf("number of valid is %d\n", c)
}

func TestValidWIthErrors(t *testing.T) {
	paths, err := readFile("day2_test.txt")
	if err != nil {
		t.Fatal(err)
	}

	c := countValidTolarablePaths(paths)
	assert.Equal(t, 4, c)
	fmt.Printf("number of valid path is %d\n", c)

	paths, err = readFile("day2.txt")
	if err != nil {
		t.Fatal(err)
	}

	c = countValidTolarablePaths(paths)
	assert.Equal(t, 544, c)
	fmt.Printf("number of valid is %d\n", c)
}
