package machine

import (
	"fmt"
	"io"
)

type TracedMachine struct {
	*Machine

	w     io.Writer
	cycle int
}

var _ Executor = &TracedMachine{}

func NewTracedMachine(w io.Writer, m *Machine) *TracedMachine {
	return &TracedMachine{
		Machine: m,
		w:       w,
		cycle:   0,
	}
}

func (t *TracedMachine) Start(trace func(*Machine)) {
	t.Machine.EnableClock()

	for {
		trace(t.Machine)
		if !t.Machine.IsClockEnabled() {
			break
		}
		t.Step()
	}
}

func (t *TracedMachine) Step() {
	op := t.Machine.Memory.ReadWord(t.Machine.Regs.PC)
	t.Machine.Regs.PC++
	opTableKey := op >> 12
	opTable[opTableKey](t, op)
	t.cycle++
}

func (t *TracedMachine) log(format string, args ...interface{}) {
	args = append([]interface{}{t.cycle}, args...)
	s := fmt.Sprintf("%d: "+format, args...)
	_, _ = t.w.Write([]byte(s))
}

func (t *TracedMachine) AddReg(dstReg Register, srcReg1 Register, srcReg2 Register) {
	t.log("ADD %s %s %s\n", dstReg, srcReg1, srcReg2)
	t.Machine.AddReg(dstReg, srcReg1, srcReg2)
}

func (t *TracedMachine) AddImm(dstReg Register, srcReg1 Register, imm5 int16) {
	t.log("ADD %s %s #%d\n", dstReg, srcReg1, imm5)
	t.Machine.AddImm(dstReg, srcReg1, imm5)
}

func (t *TracedMachine) AndReg(dstReg Register, srcReg1 Register, srcReg2 Register) {
	t.log("AND %s %s %s\n", dstReg, srcReg1, srcReg2)
	t.Machine.AndReg(dstReg, srcReg1, srcReg2)
}

func (t *TracedMachine) AndImm(dstReg Register, srcReg1 Register, imm5 int16) {
	t.log("AND %s %s #%d\n", dstReg, srcReg1, imm5)
	t.Machine.AndImm(dstReg, srcReg1, imm5)
}

func (t *TracedMachine) BRx(nzp byte, offset9 int16) {
	t.log("BRx b%0.3b #%d\n", nzp, offset9)
	t.Machine.BRx(nzp, offset9)
}

func (t *TracedMachine) JMP(baseReg Register) {
	t.log("JMP %s\n", baseReg)
	t.Machine.JMP(baseReg)
}

func (t *TracedMachine) JSR(offset11 int16) {
	t.log("JSR #%d\n", offset11)
	t.Machine.JSR(offset11)
}

func (t *TracedMachine) JSRR(baseReg Register) {
	t.log("JSRR %s\n", baseReg)
	t.Machine.JSRR(baseReg)
}

func (t *TracedMachine) LD(dstReg Register, offset9 int16) {
	t.log("LD %s #%d\n", dstReg, offset9)
	t.Machine.LD(dstReg, offset9)
}

func (t *TracedMachine) LDI(dstReg Register, offset9 int16) {
	t.log("LDI %s #%d\n", dstReg, offset9)
	t.Machine.LDI(dstReg, offset9)
}

func (t *TracedMachine) LDR(dstReg Register, baseReg Register, offset6 int16) {
	t.log("LDR %s %s #%d\n", dstReg, baseReg, offset6)
	t.Machine.LDR(dstReg, baseReg, offset6)
}

func (t *TracedMachine) LEA(dstReg Register, offset9 int16) {
	t.log("LEA %s #%d\n", dstReg, offset9)
	t.Machine.LEA(dstReg, offset9)
}

func (t *TracedMachine) Not(dstReg Register, srcReg Register) {
	t.log("Not %s %s\n", dstReg, srcReg)
	t.Machine.Not(dstReg, srcReg)
}

func (t *TracedMachine) RTI() {
	t.log("RTI\n")
	t.Machine.RTI()
}

func (t *TracedMachine) ST(srcReg Register, offset9 int16) {
	t.log("ST %s #%d\n", srcReg, offset9)
	t.Machine.ST(srcReg, offset9)
}

func (t *TracedMachine) STI(srcReg Register, offset9 int16) {
	t.log("STI %s #%d\n", srcReg, offset9)
	t.Machine.STI(srcReg, offset9)
}

func (t *TracedMachine) STR(srcReg Register, baseReg Register, offset6 int16) {
	t.log("STR %s %s #%d\n", srcReg, baseReg, offset6)
	t.Machine.STR(srcReg, baseReg, offset6)
}

func (t *TracedMachine) Trap(vec8 uint8) {
	t.log("Trap #%d\n", vec8)
	t.Machine.Trap(vec8)
}
