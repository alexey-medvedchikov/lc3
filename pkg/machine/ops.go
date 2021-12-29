package machine

func (m *Machine) AddReg(dstReg Register, srcReg1 Register, srcReg2 Register) {
	result := m.Regs.ReadRI16(srcReg1) + m.Regs.ReadRI16(srcReg2)
	m.Regs.SetRI16(dstReg, result)
	m.adjustFlags(m.Regs.ReadRU16(dstReg))
}

func (m *Machine) AddImm(dstReg Register, srcReg1 Register, imm5 int16) {
	result := m.Regs.ReadRI16(srcReg1) + imm5
	m.Regs.SetRI16(dstReg, result)
	m.adjustFlags(m.Regs.ReadRU16(dstReg))
}

func (m *Machine) AndReg(dstReg Register, srcReg1 Register, srcReg2 Register) {
	result := m.Regs.ReadRU16(srcReg1) & m.Regs.ReadRU16(srcReg2)
	m.Regs.SetRU16(dstReg, result)
	m.adjustFlags(m.Regs.ReadRU16(dstReg))
}

func (m *Machine) AndImm(dstReg Register, srcReg1 Register, imm5 int16) {
	result := m.Regs.ReadRI16(srcReg1) & imm5
	m.Regs.SetRI16(dstReg, result)
	m.adjustFlags(m.Regs.ReadRU16(dstReg))
}

func (m *Machine) BRx(nzp byte, offset9 int16) {
	n := uint16(nzp>>2) & m.Regs.GetPSRFlagN()
	z := uint16(nzp>>1) & m.Regs.GetPSRFlagZ()
	p := uint16(nzp) & m.Regs.GetPSRFlagP()

	if n|z|p == 1 {
		m.Regs.PC = addOffsetU16(m.Regs.PC, offset9)
	}
}

func (m *Machine) JMP(baseReg Register) {
	/*
		https://github.com/lassandroan/golc3/blob/main/pkg/machine/machine.go
		TODO:
				if instruction&0x1 == 1 {
				if mc.getPrivilege() {
					mc.setPrivilege(false)
				} else {
					// 0x00 Privilege Violation Vector -> 0x0100 Interrupt Addr
					mc.raiseException(0x00, mc.getPriority())
				}
			}
	*/
	m.Regs.PC = m.Regs.ReadRU16(baseReg)
}

func (m *Machine) JSR(offset11 int16) {
	m.Regs.SetRU16(R7, m.Regs.PC)
	m.Regs.PC = addOffsetU16(m.Regs.PC, offset11)
}

func (m *Machine) JSRR(baseReg Register) {
	m.Regs.SetRU16(R7, m.Regs.PC)
	m.Regs.PC = m.Regs.ReadRU16(baseReg)
}

func (m *Machine) LD(dstReg Register, offset9 int16) {
	addr := addOffsetU16(m.Regs.PC, offset9)
	val := m.Memory.ReadWord(addr)
	m.Regs.SetRU16(dstReg, val)
	m.adjustFlags(val)
}

func (m *Machine) LDI(dstReg Register, offset9 int16) {
	addr := addOffsetU16(m.Regs.PC, offset9)
	nextAddr := m.Memory.ReadWord(addr)
	val := m.Memory.ReadWord(nextAddr)
	m.Regs.SetRU16(dstReg, val)
	m.adjustFlags(val)
}

func (m *Machine) LDR(dstReg Register, baseReg Register, offset6 int16) {
	baseAddr := m.Regs.ReadRU16(baseReg)
	addr := addOffsetU16(baseAddr, offset6)
	val := m.Memory.ReadWord(addr)
	m.Regs.SetRU16(dstReg, val)
	m.adjustFlags(val)
}

func (m *Machine) LEA(dstReg Register, offset9 int16) {
	val := addOffsetU16(m.Regs.PC, offset9)
	m.Regs.SetRU16(dstReg, val)
	m.adjustFlags(val)
}

func (m *Machine) Not(dstReg Register, srcReg Register) {
	result := ^m.Regs.ReadRU16(srcReg)
	m.Regs.SetRU16(dstReg, result)
	m.adjustFlags(result)
}

func (m *Machine) RTI() {
	if m.Regs.GetPrivilegeMode() != SupervisorMode {
		// TODO: Initiate privilege mode exception
		return
	}

	ssp := m.Regs.GetSSP()
	m.Regs.PC = m.Memory.ReadWord(ssp)
	m.Regs.SetRU16(R6, m.Regs.ReadRU16(R6)+1)
	psrAddr := m.Regs.ReadRU16(R6)
	psr := m.Memory.ReadWord(psrAddr)
	m.Regs.SetRU16(R6, m.Regs.ReadRU16(R6)+1)
	m.Regs.PSR = psr
}

func (m *Machine) ST(srcReg Register, offset9 int16) {
	addr := addOffsetU16(m.Regs.PC, offset9)
	m.Memory.WriteWord(addr, m.Regs.ReadRU16(srcReg))
}

func (m *Machine) STI(srcReg Register, offset9 int16) {
	addrOfAddr := addOffsetU16(m.Regs.PC, offset9)
	addr := m.Memory.ReadWord(addrOfAddr)
	val := m.Regs.ReadRU16(srcReg)
	m.Memory.WriteWord(addr, val)
}

func (m *Machine) STR(srcReg Register, baseReg Register, offset6 int16) {
	addr := addOffsetU16(m.Regs.ReadRU16(baseReg), offset6)
	m.Memory.WriteWord(addr, m.Regs.ReadRU16(srcReg))
}

func (m *Machine) Trap(vec8 uint8) {
	m.Regs.SetRU16(R7, m.Regs.PC)
	m.Regs.PC = m.Memory.ReadWord(uint16(vec8))
}

func (m *Machine) adjustFlags(r uint16) {
	res := int16(r)
	switch {
	case res < 0:
		m.Regs.SetPSRFlagsNZP(0b100)
	case res == 0:
		m.Regs.SetPSRFlagsNZP(0b010)
	case res > 0:
		m.Regs.SetPSRFlagsNZP(0b001)
	}
}
