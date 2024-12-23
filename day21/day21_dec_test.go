package day21

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func decode(code string, pad []string) string {
	s := ""
	var p point
	if len(pad) == 2 {
		p = point{2, 0}
	} else {
		p = point{2, 3}
	}
	for _, c := range code {
		if c != 'A' {
			p = p.move(directions[byte(c)])
		} else {
			s += fmt.Sprintf("%c", pad[p.y][p.x])
		}
	}
	return s
}

func TestDecode(t *testing.T) {
	assert.Equal(t, "v<<A>>^A<A>A<AAv>A^Av<AAA>^A", decode("v<A<AA>>^AvAA<^A>Av<<A>>^AvA^Av<<A>>^AAv<A>A^A<A>Av<A<A>>^AAAvA<^A>A", keyPad))
	assert.Equal(t, "<A^A^^>AvvvA", decode("v<<A>>^A<A>A<AAv>A^Av<AAA>^A", keyPad))
	assert.Equal(t, "029A", decode("<A^A^^>AvvvA", numPad))

	assert.Equal(t, "<A>Av<<AA>^AA>AvAA^A<vAAA>^A", decode("<v<A>>^AvA^A<vA<AA>>^AAvA<^A>AAvA^A<vA>^AA<A>A<v<A>A>^AAAvA<^A>A", keyPad), "ref")
	assert.Equal(t, "<A>A<AAv<AA>>^AvAA^Av<AAA>^A", decode("v<<A>>^AvA^Av<<A>>^AAv<A<A>>^AAvAA<^A>Av<A>^AA<A>Av<A<A>>^AAAvA<^A>A", keyPad), "home")
	assert.Equal(t, "^A<<^^A>>AvvvA", decode("<A>Av<<AA>^AA>AvAA^A<vAAA>^A", keyPad), "ref")
	assert.Equal(t, "^A^^<<A>>AvvvA", decode("<A>A<AAv<AA>>^AvAA^Av<AAA>^A", keyPad), "home")
	assert.Equal(t, "379A", decode("^A<<^^A>>AvvvA", numPad), "ref")
	assert.Equal(t, "379A", decode("^A^^<<A>>AvvvA", numPad), "home")

	assert.Equal(t, "v<<A>^A>A<AA>Av<A>A^Av<AA^>A", decode("v<A<AA>>^AvA^<A>AvA^Av<<A>>^AAvA^Av<A<A>>^AvA^A<A>Av<A<A>>^AA<Av>A^A", keyPad))
	assert.Equal(t, "<^A^^Av>AvvA", decode("v<<A>^A>A<AA>Av<A>A^Av<AA^>A", keyPad))
	assert.Equal(t, "286A", decode("<^A^^Av>AvvA", numPad))
}
