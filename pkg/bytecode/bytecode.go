package bytecode

import (
	"math/bits"
)

func AddReg(dstReg Register, srcReg1 Register, srcReg2 Register) uint16 {
	// ADD       DR      SR1             SR2
	// 0 0 0 1 | x x x | x x x | 0 0 0 | x x x
	return encodeOp(0b0001) | encodeR3R3R3(dstReg, srcReg1, srcReg2)
}

func AddImm(dstReg Register, srcReg Register, imm5 int16) uint16 {
	// ADD       DR      SR1         imm5
	// 0 0 0 1 | x x x | x x x | 1 | x x x x x
	return encodeOp(0b0001) | encodeR3R3Imm5(dstReg, srcReg, imm5)
}

func AndReg(dstReg Register, srcReg1 Register, srcReg2 Register) uint16 {
	// AND       DR      SR             SR2
	// 0 1 0 1 | x x x | x x x | 0 0 0 | x x x
	return encodeOp(0b0101) | encodeR3R3R3(dstReg, srcReg1, srcReg2)
}

func AndImm(dstReg Register, srcReg Register, imm5 int16) uint16 {
	// AND       DR      SR         imm5
	// 0 1 0 1 | x x x | x x x | 1 | x x x x x
	return encodeOp(0b0101) | encodeR3R3Imm5(dstReg, srcReg, imm5)
}

func BRx(nzp byte, offset9 int16) uint16 {
	// BRx        N   Z   P   PCOffset9
	// 0 0 0 0 | x | x | x | x x x x x x x x x
	return /* encodeOp(0b0000) | */ encodeR3Imm9(Register(nzp), offset9)
}

func NOP() uint16 {
	// BRx        N   Z   P   PCOffset9
	// 0 0 0 0 | 1 | 1 | 1 | 0 0 0 0 0 0 0 0 0
	return 0b0000_11100_0000_0000
}

func JMP(baseReg Register) uint16 {
	// JMP               BaseR
	// 1 1 0 0 | 0 0 0 | x x x | 0 0 0 0 0 0
	return encodeOp(0b1100) | encodeR3R3Imm6(0, baseReg, 0)
}

func JMPT(baseReg Register) uint16 {
	// JMPT              BaseR
	// 1 1 0 0 | 0 0 0 | x x x | 0 0 0 0 0 1
	return encodeOp(0b1100) | encodeR3R3Imm6(0, baseReg, 1)
}

func RET() uint16 {
	return JMP(7)
}

func JSR(offset11 int16) uint16 {
	// JSR           PCOffset11
	// 0 1 0 0 | 1 | x x x x x x x x x x x
	return encodeOp(0b0100) |
		bits.RotateLeft16(0b1, 11) |
		uint16(offset11&0b111_1111_1111)
}

func JSRR(baseReg Register) uint16 {
	// JSRR              BaseR
	// 0 1 0 0 | 0 0 0 | x x x | 0 0 0 0 0 0
	return encodeOp(0b0100) | encodeR3R3Imm6(0, baseReg, 0)
}

func LD(dstReg Register, offset9 int16) uint16 {
	// LD        DR      PCOffset9
	// 0 0 1 0 | x x x | x x x x x x x x x
	return encodeOp(0b0010) | encodeR3Imm9(dstReg, offset9)
}

func LDI(dstReg Register, offset9 int16) uint16 {
	// LDI       DR      PCOffset9
	// 1 0 1 0 | x x x | x x x x x x x x x
	return encodeOp(0b1010) | encodeR3Imm9(dstReg, offset9)
}

func LDR(dstReg Register, baseReg Register, offset6 int16) uint16 {
	// LDR       DR      BaseR   Offset6
	// 0 1 1 0 | x x x | x x x | x x x x x x
	return encodeOp(0b0110) | encodeR3R3Imm6(dstReg, baseReg, offset6)
}

func LEA(dstReg Register, offset9 int16) uint16 {
	// LEA       DR      PCOffset9
	// 1 1 1 0 | x x x | x x x x x x x x x
	return encodeOp(0b1110) | encodeR3Imm9(dstReg, offset9)
}

func Not(dstReg Register, srcReg Register) uint16 {
	// NOT       DR      SR
	// 1 0 0 1 | x x x | x x x | 1 1 1 1 1 1
	return encodeOp(0b1001) | encodeR3R3Imm6(dstReg, srcReg, 0b11_1111)
}

func RTI() uint16 {
	return 0b1000_0000_0000_0000
}

func ST(srcReg Register, offset9 int16) uint16 {
	// ST        DR      PCOffset9
	// 0 0 1 1 | x x x | x x x x x x x x x
	return encodeOp(0b0011) | encodeR3Imm9(srcReg, offset9)
}

func STI(srcReg Register, offset9 int16) uint16 {
	// STI       DR      PCOffset9
	// 1 0 1 1 | x x x | x x x x x x x x x
	return encodeOp(0b1011) | encodeR3Imm9(srcReg, offset9)
}

func STR(srcReg Register, baseReg Register, offset6 int16) uint16 {
	// STR       SR      BaseR   PCOffset6
	// 0 1 1 1 | x x x | x x x | x x x x x x
	return encodeOp(0b0111) | encodeR3R3Imm6(srcReg, baseReg, offset6)
}

func Trap(vec uint8) uint16 {
	// TRAP                Vec8
	// 1 1 1 1 | 0 0 0 0 | x x x x x x x x
	return encodeOp(0b1111) | uint16(vec)
}

func encodeOp(op byte) uint16 {
	return bits.RotateLeft16(uint16(op), 12)
}

func encodeR3R3R3(a Register, b Register, c Register) uint16 {
	return bits.RotateLeft16(uint16(a&0b111), 9) |
		bits.RotateLeft16(uint16(b&0b111), 6) |
		uint16(c&0b111)
}

func encodeR3R3Imm6(a Register, b Register, c int16) uint16 {
	return bits.RotateLeft16(uint16(a&0b111), 9) |
		bits.RotateLeft16(uint16(b&0b111), 6) |
		uint16(c&0b11_1111)
}

func encodeR3Imm9(a Register, imm9 int16) uint16 {
	return bits.RotateLeft16(uint16(a&0b111), 9) |
		uint16(imm9&0b1_1111_1111)
}

func encodeR3R3Imm5(a Register, b Register, imm5 int16) uint16 {
	return bits.RotateLeft16(uint16(a&0b111), 9) |
		bits.RotateLeft16(0b1, 5) |
		bits.RotateLeft16(uint16(b&0b111), 6) |
		uint16(imm5&0b1_1111)
}
