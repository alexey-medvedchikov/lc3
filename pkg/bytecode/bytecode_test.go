package bytecode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddReg(t *testing.T) {
	op := AddReg(1, 1, 1)
	assert.Equal(t, uint16(0b0001_001_001_000_001), op)
}

func Test_AddImm(t *testing.T) {
	op := AddImm(1, 1, 1)
	assert.Equal(t, uint16(0b0001_001_001_1_00001), op)

	op = AddImm(1, 1, -1)
	assert.Equal(t, uint16(0b0001_001_001_1_11111), op)
}

func Test_AndReg(t *testing.T) {
	op := AndReg(1, 1, 1)
	assert.Equal(t, uint16(0b0101_001_001_000_001), op)
}

func Test_AndImm(t *testing.T) {
	op := AndImm(1, 1, 1)
	assert.Equal(t, uint16(0b0101_001_001_1_00001), op)

	op = AndImm(1, 1, -1)
	assert.Equal(t, uint16(0b0101_001_001_1_11111), op)
}

func Test_BRx(t *testing.T) {
	op := BRx(0b010, 1)
	assert.Equal(t, uint16(0b0000_0_1_0_000000001), op)

	op = BRx(0b010, -1)
	assert.Equal(t, uint16(0b0000_0_1_0_111111111), op)
}

func Test_NOP(t *testing.T) {
	op := NOP()
	assert.Equal(t, uint16(0b0000_11100_0000_0000), op)
}

func Test_JMP(t *testing.T) {
	op := JMP(0b010)
	assert.Equal(t, uint16(0b1100_000_010_000000), op)
}

func Test_JMPT(t *testing.T) {
	op := JMPT(0b010)
	assert.Equal(t, uint16(0b1100_000_010_000001), op)
}

func Test_RET(t *testing.T) {
	op := RET()
	assert.Equal(t, uint16(0b1100_000_111_000000), op)
}

func Test_JSR(t *testing.T) {
	op := JSR(1)
	assert.Equal(t, uint16(0b0100_1_00000000001), op)

	op = JSR(-1)
	assert.Equal(t, uint16(0b0100_1_11111111111), op)
}

func Test_JSRR(t *testing.T) {
	op := JSRR(0b010)
	assert.Equal(t, uint16(0b0100_000_010_000000), op)
}

func Test_LD(t *testing.T) {
	op := LD(0b010, 1)
	assert.Equal(t, uint16(0b0010_010_000000001), op)

	op = LD(0b010, -1)
	assert.Equal(t, uint16(0b0010_010_111111111), op)
}

func Test_LDI(t *testing.T) {
	op := LDI(0b010, 1)
	assert.Equal(t, uint16(0b1010_010_000000001), op)

	op = LDI(0b010, -1)
	assert.Equal(t, uint16(0b1010_010_111111111), op)
}

func Test_LDR(t *testing.T) {
	op := LDR(0b010, 0b100, 1)
	assert.Equal(t, uint16(0b0110_010_100_000001), op)

	op = LDR(0b010, 0b100, -1)
	assert.Equal(t, uint16(0b0110_010_100_111111), op)
}

func Test_LEA(t *testing.T) {
	op := LEA(0b010, 1)
	assert.Equal(t, uint16(0b1110_010_000000001), op)

	op = LEA(0b010, -1)
	assert.Equal(t, uint16(0b1110_010_111111111), op)
}

func Test_Not(t *testing.T) {
	op := Not(0b010, 0b001)
	assert.Equal(t, uint16(0b1001_010_001_111111), op)
}

func Test_RTI(t *testing.T) {
	op := RTI()
	assert.Equal(t, uint16(0b1000_0000_0000_0000), op)
}

func Test_ST(t *testing.T) {
	op := ST(0b010, 1)
	assert.Equal(t, uint16(0b0011_010_000000001), op)

	op = ST(0b010, -1)
	assert.Equal(t, uint16(0b0011_010_111111111), op)
}

func Test_STI(t *testing.T) {
	op := STI(0b010, 1)
	assert.Equal(t, uint16(0b1011_010_000000001), op)

	op = STI(0b010, -1)
	assert.Equal(t, uint16(0b1011_010_111111111), op)
}

func Test_STR(t *testing.T) {
	op := STR(0b010, 0b001, 1)
	assert.Equal(t, uint16(0b0111_010_001_000001), op)

	op = STR(0b010, 0b001, -1)
	assert.Equal(t, uint16(0b0111_010_001_111111), op)
}

func Test_Trap(t *testing.T) {
	op := Trap(0b1)
	assert.Equal(t, uint16(0b1111_0000_00000001), op)

	op = Trap(0b1111_1111)
	assert.Equal(t, uint16(0b1111_0000_11111111), op)
}
