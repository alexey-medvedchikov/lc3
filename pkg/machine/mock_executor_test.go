package machine

import "github.com/stretchr/testify/mock"

type mockExecutor struct {
	mock.Mock
}

var _ Executor = &mockExecutor{}

func (m *mockExecutor) AddReg(dstReg Register, srcReg1 Register, srcReg2 Register) {
	_ = m.Called(dstReg, srcReg1, srcReg2)
}

func (m *mockExecutor) AddImm(dstReg Register, srcReg1 Register, imm5 int16) {
	_ = m.Called(dstReg, srcReg1, imm5)
}

func (m *mockExecutor) AndReg(dstReg Register, srcReg1 Register, srcReg2 Register) {
	_ = m.Called(dstReg, srcReg1, srcReg2)
}

func (m *mockExecutor) AndImm(dstReg Register, srcReg1 Register, imm5 int16) {
	_ = m.Called(dstReg, srcReg1, imm5)
}

func (m *mockExecutor) BRx(nzp byte, offset9 int16) {
	_ = m.Called(nzp, offset9)
}

func (m *mockExecutor) JMP(baseReg Register) {
	_ = m.Called(baseReg)
}

func (m *mockExecutor) JSR(offset11 int16) {
	_ = m.Called(offset11)
}

func (m *mockExecutor) JSRR(baseReg Register) {
	_ = m.Called(baseReg)
}

func (m *mockExecutor) LD(dstReg Register, offset9 int16) {
	_ = m.Called(dstReg, offset9)
}

func (m *mockExecutor) LDI(dstReg Register, offset9 int16) {
	_ = m.Called(dstReg, offset9)
}

func (m *mockExecutor) LDR(dstReg Register, baseReg Register, offset6 int16) {
	_ = m.Called(dstReg, baseReg, offset6)
}

func (m *mockExecutor) LEA(dstReg Register, offset9 int16) {
	_ = m.Called(dstReg, offset9)
}

func (m *mockExecutor) Not(dstReg Register, srcReg Register) {
	_ = m.Called(dstReg, srcReg)
}

func (m *mockExecutor) RTI() {
	_ = m.Called()
}

func (m *mockExecutor) ST(srcReg Register, offset9 int16) {
	_ = m.Called(srcReg, offset9)
}

func (m *mockExecutor) STI(srcReg Register, offset9 int16) {
	_ = m.Called(srcReg, offset9)
}

func (m *mockExecutor) STR(srcReg Register, baseReg Register, offset6 int16) {
	_ = m.Called(srcReg, baseReg, offset6)
}

func (m *mockExecutor) Trap(vec8 uint8) {
	_ = m.Called(vec8)
}
