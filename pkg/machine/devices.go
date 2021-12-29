package machine

import (
	"fmt"
	"math/bits"
)

const (
	KeyboardStatusReg = 0xFE00
	KeyboardDataReg   = 0xFE02

	DisplayStatusReg = 0xFE04
	DisplayDataReg   = 0xFE06

	ControlReg = 0xFFFE
)

func (m *Machine) DeviceReadFunc(addr uint16) uint16 {
	switch addr {
	case ControlReg:
		return m.controlReg
	case DisplayStatusReg:
		return 0b1000_0000_0000_0000
	}
	return 0
}

func (m *Machine) DeviceWriteFunc(addr uint16, data uint16) {
	switch addr {
	case ControlReg:
		m.controlReg = data
	case DisplayDataReg:
		fmt.Printf("Display output: %c (0x%0.4x)\n", byte(data), data)
	}
}

func (m *Machine) IsClockEnabled() bool {
	return m.controlReg&0b1000_0000_0000_0000 == 0b1000_0000_0000_0000
}

func (m *Machine) EnableClock() {
	v := m.controlReg | bits.RotateLeft16(1, 15)
	m.controlReg = v
}

func (m *Machine) DisableClock() {
	v := m.controlReg & ^bits.RotateLeft16(1, 15)
	m.controlReg = v
}
