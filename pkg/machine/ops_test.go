package machine

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachine_AddReg(t *testing.T) {
	for v1 := 0; v1 < math.MaxUint16; v1 += 2048 {
		for v2 := 0; v2 < math.MaxUint16; v2 += 2048 {
			vs1 := uint16(v1)
			vs2 := uint16(v2)
			result := uint16(int16(vs1) + int16(vs2))
			m := machineWithRegs(vs1, vs2)

			m.AddReg(R2, R0, R1)

			assert.Equal(t, result, m.Regs.R[R2])
			assertPSRFlags(t, m.Regs, result)

			if t.Failed() {
				t.Logf("Input: %d + %d == %d", vs1, vs2, result)
			}
		}
	}
}

func TestMachine_AddImm(t *testing.T) {
	for v1 := 0; v1 < math.MaxUint16; v1 += 2048 {
		for v2 := -16; v2 < 16; v2++ {
			vs1 := uint16(v1)
			imm := int16(v2)
			result := uint16(int16(vs1) + imm)
			m := machineWithRegs(vs1)

			m.AddImm(R2, R0, imm)

			assert.Equal(t, result, m.Regs.R[R2])
			assertPSRFlags(t, m.Regs, result)

			if t.Failed() {
				t.Logf("Input: %d + %d == %d", vs1, imm, result)
			}
		}
	}
}

func TestMachine_AndReg(t *testing.T) {
	for v1 := 0; v1 < math.MaxUint16; v1 += 2048 {
		for v2 := 0; v2 < math.MaxUint16; v2 += 2048 {
			vs1 := uint16(v1)
			vs2 := uint16(v2)
			result := uint16(v1 & v2)
			m := machineWithRegs(vs1, vs2)

			m.AndReg(R2, R0, R1)

			assert.Equal(t, result, m.Regs.R[R2])
			assertPSRFlags(t, m.Regs, result)

			if t.Failed() {
				t.Logf("Input: %d & %d == %d", vs1, vs2, result)
			}
		}
	}
}

func TestMachine_AndImm(t *testing.T) {
	for v1 := 0; v1 < math.MaxUint16; v1 += 2048 {
		for v2 := -16; v2 < 16; v2++ {
			vs1 := uint16(v1)
			imm := int16(v2)
			result := vs1 & uint16(imm)
			m := machineWithRegs(vs1)

			m.AndImm(R2, R0, imm)

			assert.Equal(t, result, m.Regs.R[R2])
			assertPSRFlags(t, m.Regs, result)

			if t.Failed() {
				t.Logf("Input: %d & %d == %d", vs1, imm, result)
			}
		}
	}
}

func TestMachine_BRx(t *testing.T) {
	tests := []struct {
		OldPC  uint16
		PSR    uint16
		Offset int16
		NZP    byte
		NewPC  uint16
	}{
		{OldPC: 100, Offset: 10, PSR: 0b001, NZP: 0b000, NewPC: 100},
		{OldPC: 100, Offset: 10, PSR: 0b010, NZP: 0b001, NewPC: 100},
		{OldPC: 100, Offset: 10, PSR: 0b100, NZP: 0b010, NewPC: 100},
		{OldPC: 100, Offset: 10, PSR: 0b001, NZP: 0b011, NewPC: 110},
		{OldPC: 100, Offset: 10, PSR: 0b010, NZP: 0b100, NewPC: 100},
		{OldPC: 100, Offset: 10, PSR: 0b100, NZP: 0b101, NewPC: 110},
		{OldPC: 100, Offset: 10, PSR: 0b001, NZP: 0b110, NewPC: 100},
		{OldPC: 100, Offset: 10, PSR: 0b100, NZP: 0b111, NewPC: 110},
	}
	for _, tc := range tests {
		var m Machine
		m.Regs.PC = tc.OldPC
		m.Regs.PSR = tc.PSR

		m.BRx(tc.NZP, tc.Offset)

		assert.Equal(t, tc.NewPC, m.Regs.PC)
	}
}

func TestMachine_JMP(t *testing.T) {
	var m Machine
	addr := uint16(10)
	m.Regs.SetRU16(R0, addr)

	m.JMP(R0)

	assert.Equal(t, addr, m.Regs.PC)
	assert.Equal(t, uint16(0b0000_0000_0000_0000), m.Regs.PSR)
}

func TestMachine_JSR(t *testing.T) {
	tests := []struct {
		OldPC  uint16
		Offset int16
		NewPC  uint16
	}{
		{OldPC: 100, Offset: 100, NewPC: 200},
		{OldPC: 100, Offset: -50, NewPC: 50},
	}
	for _, tc := range tests {
		var m Machine
		m.Regs.PC = tc.OldPC

		m.JSR(tc.Offset)

		assert.Equal(t, tc.OldPC, m.Regs.ReadRU16(7))
		assert.Equal(t, tc.NewPC, m.Regs.PC)
	}
}

func TestMachine_JSRR(t *testing.T) {
	var m Machine
	oldPC := uint16(10)
	newPC := uint16(100)
	m.Regs.SetRU16(R1, newPC)
	m.Regs.PC = oldPC

	m.JSRR(R1)

	assert.Equal(t, oldPC, m.Regs.ReadRU16(R7))
	assert.Equal(t, newPC, m.Regs.PC)
}

func TestMachine_LD_PosOffset(t *testing.T) {
	var m Machine
	m.Memory.WriteWord(m.Regs.PC+10, 100)

	m.LD(R0, 10)

	assert.Equal(t, uint16(100), m.Regs.ReadRU16(R0))
	assert.Equal(t, uint16(0b0000_0000_0000_0001), m.Regs.PSR)
}

func TestMachine_LD_NegOffset(t *testing.T) {
	var m Machine
	m.Regs.PC = 100
	m.Memory.WriteWord(m.Regs.PC-10, 100)

	m.LD(R0, -10)

	assert.Equal(t, uint16(100), m.Regs.ReadRU16(R0))
	assert.Equal(t, uint16(0b0000_0000_0000_0001), m.Regs.PSR)
}

func TestMachine_LDI_PosOffset(t *testing.T) {
	var m Machine
	m.Memory.WriteWord(m.Regs.PC+10, 5)
	m.Memory.WriteWord(5, 200)

	m.LDI(R0, 10)

	assert.Equal(t, uint16(200), m.Regs.ReadRU16(R0))
	assert.Equal(t, uint16(0b0000_0000_0000_0001), m.Regs.PSR)
}

func TestMachine_LDI_NegOffset(t *testing.T) {
	var m Machine
	m.Regs.PC = 100
	m.Memory.WriteWord(m.Regs.PC-10, 5)
	m.Memory.WriteWord(5, 200)

	m.LDI(R0, -10)

	assert.Equal(t, uint16(200), m.Regs.ReadRU16(R0))
	assert.Equal(t, uint16(0b0000_0000_0000_0001), m.Regs.PSR)
}

func TestMachine_LDR_PosOffset(t *testing.T) {
	var m Machine
	m.Regs.SetRU16(R1, 10)
	m.Memory.WriteWord(10+10, 5)

	m.LDR(R0, R1, 10)

	assert.Equal(t, uint16(5), m.Regs.ReadRU16(R0))
	assert.Equal(t, uint16(0b0000_0000_0000_0001), m.Regs.PSR)
}

func TestMachine_LDR_NegOffset(t *testing.T) {
	var m Machine
	m.Regs.SetRU16(R1, 10)
	m.Memory.WriteWord(10-5, 5)

	m.LDR(R0, R1, -5)

	assert.Equal(t, uint16(5), m.Regs.ReadRU16(R0))
	assert.Equal(t, uint16(0b0000_0000_0000_0001), m.Regs.PSR)
}

func TestMachine_LEA_PosOffset(t *testing.T) {
	var m Machine
	pc := m.Regs.PC

	m.LEA(R0, 10)

	assert.Equal(t, pc+10, m.Regs.ReadRU16(R0))
	assert.Equal(t, uint16(0b0000_0000_0000_0001), m.Regs.PSR)
}

func TestMachine_LEA_NegOffset(t *testing.T) {
	var m Machine
	m.Regs.PC = 100
	pc := m.Regs.PC

	m.LEA(R0, -10)

	assert.Equal(t, pc-10, m.Regs.ReadRU16(R0))
	assert.Equal(t, uint16(0b0000_0000_0000_0001), m.Regs.PSR)
}

func TestMachine_Not(t *testing.T) {
	var m Machine
	m.Regs.SetRU16(R1, 0xFFFF)

	m.Not(R0, R1)

	assert.Equal(t, uint16(0), m.Regs.ReadRU16(R0))
	assert.Equal(t, uint16(0b0000_0000_0000_0010), m.Regs.PSR)
}

func TestMachine_RTI_Supervisor(t *testing.T) {
	var m Machine
	m.Regs.SetPrivilegeMode(SupervisorMode)
	m.Regs.SetRU16(R6, 10)
	m.Memory.WriteWord(10, 11)
	m.Memory.WriteWord(10+1, 20)

	m.RTI()

	assert.Equal(t, uint16(12), m.Regs.ReadRU16(R6))
	assert.Equal(t, uint16(11), m.Regs.PC)
	assert.Equal(t, uint16(20), m.Regs.PSR)
}

func TestMachine_RTI_Exception(t *testing.T) {
	var m Machine
	m.Regs.SetPrivilegeMode(UserMode)

	m.RTI()

	wantMachine := Machine{}
	wantMachine.Regs.SetPrivilegeMode(UserMode)
	assert.Equal(t, wantMachine, m)
}

func TestMachine_ST_PosOffset(t *testing.T) {
	var m Machine
	m.Regs.PC = 100
	m.Regs.SetRU16(R0, 100)

	m.ST(R0, 10)

	assert.Equal(t, uint16(100), m.Memory.ReadWord(m.Regs.PC+10))
	assert.Equal(t, uint16(0b0000_0000_0000_0000), m.Regs.PSR)
}

func TestMachine_ST_NegOffset(t *testing.T) {
	var m Machine
	m.Regs.PC = 100
	m.Regs.SetRU16(R0, 100)

	m.ST(R0, -10)

	assert.Equal(t, uint16(100), m.Memory.ReadWord(m.Regs.PC-10))
	assert.Equal(t, uint16(0b0000_0000_0000_0000), m.Regs.PSR)
}

func TestMachine_STI_PosOffset(t *testing.T) {
	var m Machine
	m.Regs.PC = 20
	m.Memory.WriteWord(m.Regs.PC+40, 20)
	m.Regs.SetRU16(R0, 100)

	m.STI(R0, 40)

	assert.Equal(t, uint16(100), m.Memory.ReadWord(20))
	assert.Equal(t, uint16(0b0000_0000_0000_0000), m.Regs.PSR)
}

func TestMachine_STI_NegOffset(t *testing.T) {
	var m Machine
	m.Regs.PC = 20
	m.Memory.WriteWord(m.Regs.PC-10, 20)
	m.Regs.SetRU16(R0, 100)

	m.STI(R0, -10)

	assert.Equal(t, uint16(100), m.Memory.ReadWord(20))
	assert.Equal(t, uint16(0b0000_0000_0000_0000), m.Regs.PSR)
}

func TestMachine_STR_PosOffset(t *testing.T) {
	var m Machine
	m.Regs.SetRU16(R0, 100)
	m.Regs.SetRU16(R1, 20)

	m.STR(R0, R1, 10)

	assert.Equal(t, uint16(100), m.Memory.ReadWord(30))
	assert.Equal(t, uint16(0b0000_0000_0000_0000), m.Regs.PSR)
}

func TestMachine_STR_NegOffset(t *testing.T) {
	var m Machine
	m.Regs.SetRU16(R0, 100)
	m.Regs.SetRU16(R1, 20)

	m.STR(R0, R1, -10)

	assert.Equal(t, uint16(100), m.Memory.ReadWord(10))
	assert.Equal(t, uint16(0b0000_0000_0000_0000), m.Regs.PSR)
}

func TestMachine_Trap(t *testing.T) {
	var m Machine
	oldPC := uint16(100)
	m.Regs.PC = oldPC
	vec := uint8(0x0F)
	vecAddr := uint16(400)
	m.Memory.WriteWord(uint16(vec), vecAddr)

	m.Trap(vec)

	assert.Equal(t, oldPC, m.Regs.ReadRU16(R7))
	assert.Equal(t, vecAddr, m.Regs.PC)
}

func assertPSRFlags(t *testing.T, r Regs, res uint16) {
	ri16 := int16(res)
	switch {
	case ri16 < 0:
		assert.Equal(t, uint16(1), r.GetPSRFlagN())
		assert.Equal(t, uint16(0), r.GetPSRFlagZ())
		assert.Equal(t, uint16(0), r.GetPSRFlagP())
	case ri16 == 0:
		assert.Equal(t, uint16(0), r.GetPSRFlagN())
		assert.Equal(t, uint16(1), r.GetPSRFlagZ())
		assert.Equal(t, uint16(0), r.GetPSRFlagP())
	case ri16 > 0:
		assert.Equal(t, uint16(0), r.GetPSRFlagN())
		assert.Equal(t, uint16(0), r.GetPSRFlagZ())
		assert.Equal(t, uint16(1), r.GetPSRFlagP())
	}
}

func machineWithRegs(regs ...uint16) Machine {
	var m Machine
	for i, v := range regs {
		m.Regs.R[i] = v
	}

	return m
}
