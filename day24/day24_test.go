package day24

import (
	"fmt"
	"math/rand/v2"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type gates map[string]int

type oper struct {
	left, right, op string
	res             string
}

func loadInput(file string) (gates, map[string]*oper, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}
	lines := strings.Split(string(content), "\n")
	var i int
	g := make(gates, 0)
	for i = 0; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])
		if l == "" {
			break
		}
		parts := strings.Split(l, ":")
		g[strings.TrimSpace(parts[0])], err = strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, nil, err
		}
	}
	ops := make(map[string]*oper, 0)
	for i++; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])
		var op oper
		split := strings.Split(l, "->")
		op.res = strings.TrimSpace(split[1])
		parts := strings.Split(strings.TrimSpace(split[0]), " ")
		op.left = parts[0]
		op.op = parts[1]
		op.right = parts[2]
		ops[op.res] = &op
	}
	return g, ops, err
}

type namedBits struct {
	name  string
	value int
}

func value(g gates, prefix string) int {
	bits := make([]namedBits, 0)
	for k, v := range g {
		if strings.HasPrefix(k, prefix) {
			bits = append(bits, namedBits{k, v})
		}
	}
	sort.SliceStable(bits, func(i, j int) bool {
		return strings.Compare(bits[i].name, bits[j].name) < 0
	})
	z := 0
	d := 0
	for _, b := range bits {
		z = z | (b.value << d)
		d++
	}
	return z
}

func computeZ(g gates, ops map[string]*oper) int {
	someOp := true
	for someOp {
		someOp = false
		for _, op := range ops {
			if _, processed := g[op.res]; processed {
				continue
			}
			l, existsL := g[op.left]
			r, existsR := g[op.right]
			if !existsL || !existsR {
				continue
			}
			switch op.op {
			case "AND":
				g[op.res] = l & r
			case "OR":
				g[op.res] = l | r
			case "XOR":
				g[op.res] = l ^ r
			default:
				panic("invlaid op")
			}
			someOp = true
		}
	}
	return value(g, "z")
}

func TestDay24Part1(t *testing.T) {
	g, ops, err := loadInput("input_min.txt")
	require.NoError(t, err)
	assert.Equal(t, 4, computeZ(g, ops))

	g, ops, err = loadInput("input_test.txt")
	require.NoError(t, err)
	assert.Equal(t, 2024, computeZ(g, ops))

	g, ops, err = loadInput("input.txt")
	require.NoError(t, err)
	assert.Equal(t, 42883464055378, computeZ(g, ops))
}

func findBugs(ops map[string]*oper) [][2]string {
	// z0 = x0 XOR y0
	// z1 = (x1 XOR y1) XOR (x0 AND y0)
	// z2 = (y2 XOR x2) XOR (y1 AND x1 OR x0 AND y0 AND x1 XOR y1)
	// z3 = (y3 XOR x3) XOR (y2 AND x2 OR x0 AND y0 AND x1 AND y1 AND x2 XOR y2)
	// zN = ...
	var changes [][2]string
	zXor := make(map[string]*oper)
	for _, op := range ops {
		if op.res[0] == 'z' && op.op != "XOR" {
			zXor[op.res] = op
		}
	}
	for _, op := range ops {
		// xNN XOR yNN -> aaa
		if op.left[1:] == op.right[1:] && op.op == "XOR" && op.res[0] != 'z' {
			for _, opXor := range ops {
				// aaa XOR bbb -> zNN (zNN not found)
				if (opXor.left == op.res || opXor.right == op.res) && opXor.op == "XOR" && zXor["z"+op.right[1:]] != nil {
					changes = append(changes, [2]string{"z" + op.right[1:], opXor.res})
					continue
				}

				if opXor.res == "z"+op.right[1:] && opXor.op == "XOR" {
					// ccc XOR ddd -> zNN : ddd or ccc must be aaa
					if opXor.left != op.res && opXor.right != op.res {
						if ops[opXor.left].op == "XOR" || ops[opXor.left].op == "AND" {
							changes = append(changes, [2]string{op.res, opXor.left})
						}
						if ops[opXor.right].op == "XOR" || ops[opXor.right].op == "AND" {
							changes = append(changes, [2]string{op.res, opXor.right})
						}
					}
				}
			}
		}
	}
	return changes
}

func TestDay24Part2(t *testing.T) {
	g, ops, err := loadInput("input.txt")
	require.NoError(t, err)

	changes := findBugs(ops)
	assert.Equal(t, 4, len(changes))
	lst := make([]string, 0, 8)
	for _, c := range changes {
		lst = append(lst, c[0], c[1])
		op0 := ops[c[0]]
		op1 := ops[c[1]]
		op0.res = c[1]
		ops[c[1]] = op0
		op1.res = c[0]
		ops[c[0]] = op1
	}
	sort.Strings(lst)
	assert.Equal(t, "dqr,dtk,pfw,shh,vgs,z21,z33,z39", strings.Join(lst, ","))

	x := value(g, "x")
	y := value(g, "y")
	assert.Equal(t, 22044444944409, x)
	assert.Equal(t, 20280608350777, y)
	z := x + y
	assert.Equal(t, 42325053295186, z)
	zC := computeZ(g, ops)
	assert.Equal(t, z, zC)
	assert.Equal(t, fmt.Sprintf("%b", z), fmt.Sprintf("%b", zC))

	// check on more values (first findBugs did only found 3 switches required)
	digits := 45
	for i := 0; i < 10; i++ {
		x := rand.IntN(1 << digits)
		y := rand.IntN(1 << digits)
		z := x + y
		g := make(map[string]int)
		for d := 0; d < digits; d++ {
			g[fmt.Sprintf("x%02d", d)] = (x & (1 << d)) >> d
			g[fmt.Sprintf("y%02d", d)] = (y & (1 << d)) >> d
		}
		zC := computeZ(g, ops)
		assert.Equal(t, z, zC)
		assert.Equal(t, fmt.Sprintf("%b", z), fmt.Sprintf("%b", zC))
		if t.Failed() {
			break
		}
	}
}
