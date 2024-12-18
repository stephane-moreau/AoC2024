package day17

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type registers struct {
	a, b, c int
	out     []int
}

func adv(reg *registers, op int) int {
	reg.a = reg.a >> op
	return -1
}
func bxl(reg *registers, op int) int {
	reg.b = reg.b ^ op
	return -1
}

func bst(reg *registers, op int) int {
	reg.b = op % 8
	return -1
}

func jnz(reg *registers, op int) int {
	if reg.a != 0 {
		return op
	}
	return -1
}

func bxc(reg *registers, op int) int {
	reg.b = (reg.b ^ reg.c)
	return -1
}

func out(reg *registers, op int) int {
	reg.out = append(reg.out, op%8)
	return -1
}

func bdv(reg *registers, op int) int {
	reg.b = reg.a >> op
	return -1
}

func cdv(reg *registers, op int) int {
	reg.c = reg.a >> op
	return -1
}

type opCode func(*registers, int) int

var functions = map[int]opCode{
	0: adv,
	1: bxl,
	2: bst,
	3: jnz,
	4: bxc,
	5: out,
	6: bdv,
	7: cdv,
}

func operand(reg *registers, op, oper int) *int {
	if op == 1 {
		return &oper
	}
	if oper <= 3 {
		return &oper
	}
	if oper == 4 {
		return &reg.a
	}
	if oper == 5 {
		return &reg.b
	}
	if oper == 6 {
		return &reg.c
	}
	panic("unsupported opcode")
}

func exec(reg registers, ops []int) registers {
	for i := 0; i < len(ops); i += 2 {
		ptr := functions[ops[i]](&reg, *operand(&reg, ops[i], ops[i+1]))
		if ptr != -1 {
			i = ptr - 2
		}
	}
	return reg
}

func TestFucntions(t *testing.T) {
	r := exec(registers{c: 9}, []int{2, 6})
	assert.Equal(t, 1, r.b)

	r = exec(registers{a: 10}, []int{5, 0, 5, 1, 5, 4})
	assert.Equal(t, []int{0, 1, 2}, r.out)

	r = exec(registers{a: 2024}, []int{0, 1, 5, 4, 3, 0})
	assert.Equal(t, []int{4, 2, 5, 6, 7, 7, 7, 7, 3, 1, 0}, r.out)
	assert.Equal(t, 0, r.a)

	r = exec(registers{b: 29}, []int{1, 7})
	assert.Equal(t, 26, r.b)

	r = exec(registers{b: 2024, c: 43690}, []int{4, 0})
	assert.Equal(t, 44354, r.b)
}

func TestDay17Part1(t *testing.T) {
	r := exec(registers{a: 729}, []int{0, 1, 5, 4, 3, 0})
	assert.Equal(t, []int{4, 6, 3, 5, 6, 3, 5, 2, 1, 0}, r.out)

	r = exec(registers{a: 41_644_071}, []int{2, 4, 1, 2, 7, 5, 1, 7, 4, 4, 0, 3, 5, 5, 3, 0})
	assert.Equal(t, []int{3, 1, 5, 3, 7, 4, 2, 7, 5}, r.out)
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func reverseProg(prog []int) int {
	res := 0
	for i := len(prog) - 1; i >= 0; i-- {
		found := false
		for !found {
			resBase := res
			for c := 0; c < 2048; c++ {
				t := resBase + c
				r := registers{a: t}
				u := exec(r, prog)
				if equal(u.out, prog[i:]) {
					if i != 0 {
						res = t << 3
					} else {
						res = t
					}
					found = true
					break
				}
			}
			if !found {
				panic("WTF")
			}
		}
	}
	return res
}

func TestDay17Part2(t *testing.T) {
	// 117440
	prog := []int{0, 3, 5, 4, 3, 0}
	assert.Equal(t, 117440, reverseProg(prog))

	prog = []int{2, 4, 1, 2, 7, 5, 1, 7, 4, 4, 0, 3, 5, 5, 3, 0}
	assert.Equal(t, 190593310997519, reverseProg(prog))
}
