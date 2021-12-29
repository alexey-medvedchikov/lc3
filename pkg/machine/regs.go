package machine

import (
	"fmt"
	"math/bits"
)

type Regs struct {
	// R is General Purpose Registers
	R [8]uint16

	// PC is a Program Counter
	PC uint16

	// PSR is a Processor Status Register
	// 15     privilege mode
	// 14..11 reserved
	// 10..8  priority level
	// 7..3   reserved
	// 2      Condition code N (Negative)
	// 1      Condition code Z (Zero)
	// 0      Condition code P (Positive)
	PSR uint16
}

type Register byte

func (r Register) String() string {
	return fmt.Sprintf("R%d", r)
}

const (
	SupervisorMode = 0b0
	UserMode       = 0b1

	R0 Register = 0
	R1 Register = 1
	R2 Register = 2
	R3 Register = 3
	R4 Register = 4
	R5 Register = 5
	R6 Register = 6
	R7 Register = 7

	// PL0 is for Priority Level 0, lowest to highest
	PL0 = 0x0
	PL1 = 0x1
	PL2 = 0x2
	PL3 = 0x3
	PL4 = 0x4
	PL5 = 0x5
	PL6 = 0x6
	PL7 = 0x7

	privModeMask  uint16 = 0b1 << 15
	prioLevelMask uint16 = 0b111 << 8

	// regSize is register size in bytes
	regSize = 2
)

// ReadRU16 reads the contents of a register reg and returns it as uint16
func (r *Regs) ReadRU16(reg Register) uint16 {
	return r.R[reg]
}

// ReadRI16 reads the contents of a register reg and returns it as int16
func (r *Regs) ReadRI16(reg Register) int16 {
	return int16(r.R[reg])
}

// SetRU16 sets contents of a register reg to uint16 value of v
func (r *Regs) SetRU16(reg Register, v uint16) {
	r.R[reg] = v
}

// SetRI16 sets contents of a register reg to int16 value of v
func (r *Regs) SetRI16(reg Register, v int16) {
	r.R[reg] = uint16(v)
}

// GetPrivilegeMode returns value of privilege mode from PSR register as a bit in position 0
func (r *Regs) GetPrivilegeMode() uint16 {
	return (r.PSR & privModeMask) >> 15
}

// SetPrivilegeMode sets privilege mode in PSR register
func (r *Regs) SetPrivilegeMode(v uint16) {
	r.PSR = (r.PSR & ^privModeMask) | bits.RotateLeft16(v, 15)
}

// GetPriorityLevel returns value of priority level from PSR register as a bit in position 0
func (r *Regs) GetPriorityLevel() uint16 {
	return (r.PSR & prioLevelMask) >> 8
}

// SetPriorityLevel sets priority level in PSR register
func (r *Regs) SetPriorityLevel(v uint16) {
	r.PSR = (r.PSR & ^prioLevelMask) | bits.RotateLeft16(v, 8)
}

func (r *Regs) SetPSRFlagsNZP(b uint16) {
	r.PSR = (r.PSR & ^uint16(0b111)) | (b & uint16(0b111))
}

func (r *Regs) GetPSRFlagN() uint16 { return r.PSR >> 2 & 0b1 }
func (r *Regs) GetPSRFlagZ() uint16 { return r.PSR >> 1 & 0b1 }
func (r *Regs) GetPSRFlagP() uint16 { return r.PSR & 0b1 }

func (r *Regs) GetSSP() uint16 { return r.ReadRU16(R6) }

func (r *Regs) Reset() {
	for i := range r.R {
		r.SetRU16(Register(i), 0)
	}
	r.PC = 0
	r.PSR = 0
}
