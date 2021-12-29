package machine

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexey-medvedchikov/lc3/pkg/bytecode"
)

func TestTracedMachine_Start(t *testing.T) {
	var buf bytes.Buffer
	m := NewTracedMachine(&buf, &Machine{})

	program := []uint16{
		bytecode.AndImm(bytecode.R0, bytecode.R0, 0),
		bytecode.STI(bytecode.R0, 0),
		ControlReg,
	}

	m.Memory.WriteSegment(UserStart, program)

	m.Init()

	steps := 0
	m.Start(func(m *Machine) {
		if steps > 50 {
			m.DisableClock()
		}
		steps++
	})

	assert.Equal(t, "0: AND R0 R0 #0\n1: STI R0 #0\n", buf.String())
}
