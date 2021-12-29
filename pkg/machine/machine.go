package machine

import "github.com/alexey-medvedchikov/lc3/pkg/bytecode"

type Executor interface {
	AddReg(dstReg Register, srcReg1 Register, srcReg2 Register)
	AddImm(dstReg Register, srcReg1 Register, imm5 int16)
	AndReg(dstReg Register, srcReg1 Register, srcReg2 Register)
	AndImm(dstReg Register, srcReg1 Register, imm5 int16)
	BRx(nzp byte, offset9 int16)
	JMP(baseReg Register)
	JSR(offset11 int16)
	JSRR(baseReg Register)
	LD(dstReg Register, offset9 int16)
	LDI(dstReg Register, offset9 int16)
	LDR(dstReg Register, baseReg Register, offset6 int16)
	LEA(dstReg Register, offset9 int16)
	Not(dstReg Register, srcReg Register)
	RTI()
	ST(srcReg Register, offset9 int16)
	STI(srcReg Register, offset9 int16)
	STR(srcReg Register, baseReg Register, offset6 int16)
	Trap(vec8 uint8)
}

type Machine struct {
	Regs   Regs
	Memory Memory

	controlReg uint16
}

func (m *Machine) Init() {
	trapPos := [...]uint16{
		bytecode.TrapGETCAddr, bytecode.TrapOUTAddr, bytecode.TrapPUTSAddr,
		bytecode.TrapINAddr, bytecode.TrapPUTSPAddr, bytecode.TrapHALTAddr,
	}
	trapCode := [...][]uint16{
		bytecode.TrapGETC, bytecode.TrapOUT, bytecode.TrapPUTS,
		bytecode.TrapIN, bytecode.TrapPUTSP, bytecode.TrapHALT,
	}

	pos := PrivilegedStart
	for i := range trapPos {
		m.Memory.WriteSegment(pos, trapCode[i])
		m.Memory.WriteWord(trapPos[i], pos)
		pos += uint16(len(trapCode[i]))
	}

	m.Memory.DeviceWriteFunc = m.DeviceWriteFunc
	m.Memory.DeviceReadFunc = m.DeviceReadFunc

	m.Regs.Reset()
	m.Regs.PC = UserStart
	m.Regs.SetRU16(R6, UserEnd)
}

func (m *Machine) Start() {
	m.EnableClock()

	for m.IsClockEnabled() {
		m.Step()
	}
}

func (m *Machine) Step() {
	op := m.Memory.ReadWord(m.Regs.PC)
	m.Regs.PC++
	opTableKey := op >> 12
	opTable[opTableKey](m, op)
}

var opTable = [...]func(ex Executor, op uint16){
	decodeBR,      // 0b0000 = 0
	decodeAdd,     // 0b0001 = 1
	decodeLD,      // 0b0010 = 2
	decodeST,      // 0b0011 = 3
	decodeJSR,     // 0b0100 = 4
	decodeAnd,     // 0b0101 = 5
	decodeLDR,     // 0b0110 = 6
	decodeSTR,     // 0b0111 = 7
	decodeRTI,     // 0b1000 = 8
	decodeNot,     // 0b1001 = 9
	decodeLDI,     // 0b1010 = 10
	decodeSTI,     // 0b1011 = 11
	decodeJMP,     // 0b1100 = 12
	decodeInvalid, // 0b1101 = 13
	decodeLEA,     // 0b1110 = 14
	decodeTrap,    // 0b1111 = 15
}

func decodeInvalid(_ Executor, _ uint16) {
	panic(interface{}("Invalid operation"))
}

func decodeBR(m Executor, op uint16) {
	m.BRx(decodeB3Offset9(op))
}

func decodeAdd(m Executor, op uint16) {
	if op&0b10_0000 == 0 {
		if op&0b0000_0000_0001_1000 == 0 {
			m.AddReg(decodeR3R3R3(op))
		} else {
			decodeInvalid(m, op)
		}
	} else {
		m.AddImm(decodeR3R3Imm5(op))
	}
}

func decodeAnd(m Executor, op uint16) {
	if op&0b10_0000 == 0 {
		if op&0b0000_0000_0001_1000 == 0 {
			m.AndReg(decodeR3R3R3(op))
		} else {
			decodeInvalid(m, op)
		}
	} else {
		m.AndImm(decodeR3R3Imm5(op))
	}
}

func decodeJSR(m Executor, op uint16) {
	if op&0b1000_0000_0000 == 0 {
		br := Register(op >> 6 & 0b111)
		m.JSRR(br)
	} else {
		offset := signExtend11(op & 0b111_1111_1111)
		m.JSR(offset)
	}
}

func decodeRTI(m Executor, op uint16) {
	if op&0b0000_1111_1111_1111 == 0 {
		m.RTI()
	} else {
		decodeInvalid(m, op)
	}
}

func decodeJMP(m Executor, op uint16) {
	if op&0b0000_1110_0011_1111 == 0 {
		br := Register(op >> 6 & 0b111)
		m.JMP(br)
	} else {
		decodeInvalid(m, op)
	}
}

func decodeNot(m Executor, op uint16) {
	if op&0b0000_0000_0011_1111 == 0b11_1111 {
		dr := Register(op >> 9 & 0b111)
		sr := Register(op >> 6 & 0b111)
		m.Not(dr, sr)
	} else {
		decodeInvalid(m, op)
	}
}

func decodeTrap(m Executor, op uint16) {
	if op&0b0000_1111_0000_0000 == 0 {
		vec := byte(op & 0b1111_1111)
		m.Trap(vec)
	}
}

func decodeLD(m Executor, op uint16)  { m.LD(decodeR3Offset9(op)) }
func decodeLDI(m Executor, op uint16) { m.LDI(decodeR3Offset9(op)) }
func decodeLDR(m Executor, op uint16) { m.LDR(decodeR3R3Imm6(op)) }
func decodeLEA(m Executor, op uint16) { m.LEA(decodeR3Offset9(op)) }
func decodeST(m Executor, op uint16)  { m.ST(decodeR3Offset9(op)) }
func decodeSTI(m Executor, op uint16) { m.STI(decodeR3Offset9(op)) }
func decodeSTR(m Executor, op uint16) { m.STR(decodeR3R3Imm6(op)) }

func decodeR3R3R3(op uint16) (Register, Register, Register) {
	dr := Register(op >> 9 & 0b111)
	sr1 := Register(op >> 6 & 0b111)
	sr2 := Register(op & 0b111)

	return dr, sr1, sr2
}

func decodeR3R3Imm5(op uint16) (Register, Register, int16) {
	dr := Register(op >> 9 & 0b111)
	sr := Register(op >> 6 & 0b111)
	imm := signExtend5(op & 0b1_1111)

	return dr, sr, imm
}

func decodeR3R3Imm6(op uint16) (Register, Register, int16) {
	dr := Register(op >> 9 & 0b111)
	sr := Register(op >> 6 & 0b111)
	imm := signExtend6(op & 0b11_1111)

	return dr, sr, imm
}

func decodeR3Offset9(op uint16) (Register, int16) {
	dr := Register(op >> 9 & 0b111)
	offset := signExtend9(op & 0b1_1111_1111)

	return dr, offset
}

func decodeB3Offset9(op uint16) (byte, int16) {
	dr := byte(op >> 9 & 0b111)
	offset := signExtend9(op & 0b1_1111_1111)

	return dr, offset
}
