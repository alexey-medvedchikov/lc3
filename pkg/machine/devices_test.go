package machine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachine_DisableClock(t *testing.T) {
	m := Machine{
		controlReg: 0xFFFF,
	}
	m.DisableClock()
	assert.Equal(t, uint16(0b0111_1111_1111_1111), m.controlReg)
}

func TestMachine_EnableClock(t *testing.T) {
	var m Machine
	m.EnableClock()
	assert.Equal(t, uint16(0b1000_0000_0000_0000), m.controlReg)
}

func TestMachine_IsClockEnabled_False(t *testing.T) {
	var m Machine
	assert.False(t, m.IsClockEnabled())
}

func TestMachine_IsClockEnabled_True(t *testing.T) {
	m := Machine{
		controlReg: 0xFFFF,
	}
	assert.True(t, m.IsClockEnabled())
}
