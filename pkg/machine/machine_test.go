package machine

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexey-medvedchikov/lc3/pkg/bytecode"
)

func TestMachine_Start(t *testing.T) {
	var m Machine

	program := []uint16{
		bytecode.LEA(bytecode.R0, 2),
		bytecode.Trap(0x22),
		bytecode.Trap(0x25),
		binary.LittleEndian.Uint16([]byte{'H', 0}),
		binary.LittleEndian.Uint16([]byte{'i', 0}),
		binary.LittleEndian.Uint16([]byte{'!', 0}),
		0x0000,
	}

	m.Memory.WriteSegment(UserStart, program)

	m.Init()
	m.Start()
}

func Test_decodeBRx_PosOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("BRx", byte(0b111), int16(1)).Once()

	decodeBR(m, bytecode.BRx(0b111, 1))
}

func Test_decodeBRx_NegOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("BRx", byte(0b111), int16(-1)).Once()

	decodeBR(m, bytecode.BRx(0b111, -1))
}

func Test_decodeAdd_Reg(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("AddReg", R5, R6, R7).Once()

	decodeAdd(m, bytecode.AddReg(bytecode.R5, bytecode.R6, bytecode.R7))
}

func Test_decodeAdd_PosImm(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("AddImm", R5, R6, int16(2)).Once()

	decodeAdd(m, bytecode.AddImm(bytecode.R5, bytecode.R6, 2))
}

func Test_decodeAdd_NegImm(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("AddImm", R5, R6, int16(-1)).Once()

	decodeAdd(m, bytecode.AddImm(bytecode.R5, bytecode.R6, -1))
}

func Test_decodeAdd_Invalid(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)

	assert.Panics(t, func() { decodeAdd(m, 0b0001_000_000_010_000) })
	assert.Panics(t, func() { decodeAdd(m, 0b0001_000_000_001_000) })
}

func Test_decodeAnd_Reg(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("AndReg", R5, R6, R7).Once()

	decodeAnd(m, bytecode.AndReg(bytecode.R5, bytecode.R6, bytecode.R7))
}

func Test_decodeAnd_PosImm(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("AndImm", R5, R6, int16(2)).Once()

	decodeAnd(m, bytecode.AndImm(bytecode.R5, bytecode.R6, 2))
}

func Test_decodeAnd_NegImm(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("AndImm", R5, R6, int16(-1)).Once()

	decodeAnd(m, bytecode.AndImm(bytecode.R5, bytecode.R6, -1))
}

func Test_decodeAnd_Invalid(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)

	assert.Panics(t, func() { decodeAnd(m, 0b0101_000_000_010_000) })
	assert.Panics(t, func() { decodeAnd(m, 0b0101_000_000_001_000) })
}

func Test_decodeLD_PosOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("LD", R5, int16(1)).Once()

	decodeLD(m, bytecode.LD(bytecode.R5, 1))
}

func Test_decodeLD_NegOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("LD", R5, int16(-1)).Once()

	decodeLD(m, bytecode.LD(bytecode.R5, -1))
}

func Test_decodeST_PosOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("ST", R5, int16(1)).Once()

	decodeST(m, bytecode.ST(bytecode.R5, 1))
}

func Test_decodeST_NegOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("ST", R5, int16(-1)).Once()

	decodeST(m, bytecode.ST(bytecode.R5, -1))
}

func Test_decodeJSR_PosOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("JSR", int16(1)).Once()

	decodeJSR(m, bytecode.JSR(1))
}

func Test_decodeJSR_NegOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("JSR", int16(-1)).Once()

	decodeJSR(m, bytecode.JSR(-1))
}

func Test_decodeJSRR(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("JSRR", R6).Once()

	decodeJSR(m, bytecode.JSRR(bytecode.R6))
}

func Test_decodeLDR_PosOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("LDR", R5, R6, int16(2)).Once()

	decodeLDR(m, bytecode.LDR(bytecode.R5, bytecode.R6, 2))
}

func Test_decodeLDR_NegOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("LDR", R6, R7, int16(-1)).Once()

	decodeLDR(m, bytecode.LDR(bytecode.R6, bytecode.R7, -1))
}

func Test_decodeSTR_PosOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("STR", R5, R6, int16(2)).Once()

	decodeSTR(m, bytecode.STR(bytecode.R5, bytecode.R6, 2))
}

func Test_decodeSTR_NegOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("STR", R5, R6, int16(-1)).Once()

	decodeSTR(m, bytecode.STR(bytecode.R5, bytecode.R6, -1))
}

func Test_decodeRTI(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("RTI").Once()

	decodeRTI(m, bytecode.RTI())
}

func Test_decodeRTI_Invalid(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)

	assert.Panics(t, func() { decodeRTI(m, 0b1000_1010_1010_1010) })
}

func Test_decodeNot(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("Not", R5, R6).Once()

	decodeNot(m, bytecode.Not(bytecode.R5, bytecode.R6))
}

func Test_decodeNot_Invalid(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)

	assert.Panics(t, func() { decodeNot(m, 0b1001_010_010_010101) })
}

func Test_decodeLDI_PosOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("LDI", R5, int16(1)).Once()

	decodeLDI(m, bytecode.LDI(bytecode.R5, 1))
}

func Test_decodeLDI_NegOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("LDI", R5, int16(-1)).Once()

	decodeLDI(m, bytecode.LDI(bytecode.R5, -1))
}

func Test_decodeSTI_PosOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("STI", R5, int16(1)).Once()

	decodeSTI(m, bytecode.STI(bytecode.R5, 1))
}

func Test_decodeSTI_NegOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("STI", R5, int16(-1)).Once()

	decodeSTI(m, bytecode.STI(bytecode.R5, -1))
}

func Test_decodeJMP(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("JMP", R6).Once()

	decodeJMP(m, bytecode.JMP(bytecode.R6))
}

func Test_decodeJMP_Invalid(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)

	assert.Panics(t, func() { decodeJMP(m, 0b1100_010_010_010101) })
}

func Test_decodeLEA_PosOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("LEA", R5, int16(1)).Once()

	decodeLEA(m, bytecode.LEA(bytecode.R5, 1))
}

func Test_decodeLEA_NegOffset(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("LEA", R5, int16(-1)).Once()

	decodeLEA(m, bytecode.LEA(bytecode.R5, -1))
}

func Test_decodeTrap(t *testing.T) {
	m := &mockExecutor{}
	defer m.AssertExpectations(t)
	m.On("Trap", uint8(1)).Once()

	decodeTrap(m, bytecode.Trap(1))
}
