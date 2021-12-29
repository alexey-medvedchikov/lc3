package machine

import (
	"math"
	"math/bits"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisters_GetPrivilegeMode(t *testing.T) {
	r := Regs{PSR: bits.RotateLeft16(0b1, 15)}
	assert.Equal(t, uint16(0b1), r.GetPrivilegeMode())

	r = Regs{PSR: ^bits.RotateLeft16(0b1, 15)}
	assert.Equal(t, uint16(0b0), r.GetPrivilegeMode())
}

func TestRegisters_SetPrivilegeMode(t *testing.T) {
	r := Regs{}
	r.SetPrivilegeMode(UserMode)
	assert.Equal(t, bits.RotateLeft16(0b1, 15), r.PSR)

	r = Regs{PSR: math.MaxUint16}
	r.SetPrivilegeMode(SupervisorMode)
	assert.Equal(t, ^bits.RotateLeft16(0b1, 15), r.PSR)
}

func TestRegisters_SetPriorityLevel(t *testing.T) {
	r := Regs{PSR: bits.RotateLeft16(0b111, 8)}
	assert.Equal(t, uint16(0b111), r.GetPriorityLevel())

	r = Regs{PSR: ^bits.RotateLeft16(0b111, 8)}
	assert.Equal(t, uint16(0b0), r.GetPriorityLevel())
}

func TestRegisters_GetPriorityLevel(t *testing.T) {
	r := Regs{}
	r.SetPriorityLevel(0b111)
	assert.Equal(t, bits.RotateLeft16(0b111, 8), r.PSR)

	r = Regs{PSR: math.MaxUint16}
	r.SetPriorityLevel(0b000)
	assert.Equal(t, ^bits.RotateLeft16(0b111, 8), r.PSR)
}

func TestRegs_SetPSRFlagsNZP(t *testing.T) {
	r := Regs{PSR: math.MaxUint16}
	r.SetPSRFlagsNZP(0b000)
	assert.Equal(t, ^uint16(0b111), r.PSR)
}

func TestRegisters_GetPSRFlagN(t *testing.T) {
	r := Regs{PSR: bits.RotateLeft16(0b1, 2)}
	assert.Equal(t, uint16(0b1), r.GetPSRFlagN())
}

func TestRegisters_GetPSRFlagZ(t *testing.T) {
	r := Regs{PSR: bits.RotateLeft16(0b1, 1)}
	assert.Equal(t, uint16(0b1), r.GetPSRFlagZ())
}

func TestRegisters_GetPSRFlagP(t *testing.T) {
	r := Regs{PSR: 0b1}
	assert.Equal(t, uint16(0b1), r.GetPSRFlagP())
}

func TestRegisters_GetSSP(t *testing.T) {
	var r Regs
	v := uint16(0x0BAD)
	r.SetRU16(R6, v)
	assert.Equal(t, v, r.GetSSP())
}

func TestRegs_Reset(t *testing.T) {
	r := Regs{
		R:   [...]uint16{1, 2, 3, 4, 5, 6, 7, 8},
		PC:  0xFEFE,
		PSR: 0xABAB,
	}
	r.Reset()

	assert.Zero(t, r)
}
