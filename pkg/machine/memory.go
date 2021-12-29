package machine

import (
	"io"
	"unsafe"
)

const (
	MemoryEnd  uint16 = 0xFFFF
	MemorySize int    = int(MemoryEnd) + 1

	TrapVecTblStart uint16 = 0x0000
	TrapVecTblEnd   uint16 = 0x00FF
	IntVecTblStart  uint16 = 0x0100
	IntVecTblEnd    uint16 = 0x01FF
	PrivilegedStart uint16 = 0x0200
	PrivilegedEnd   uint16 = 0x2FFF
	UserStart       uint16 = 0x3000
	UserEnd         uint16 = 0xFDFF
	DeviceRegStart  uint16 = 0xFE00
	DeviceRegEnd    uint16 = 0xFFFF
)

type Memory struct {
	mem             [DeviceRegStart]uint16
	DeviceReadFunc  func(addr uint16) uint16
	DeviceWriteFunc func(addr uint16, data uint16)
}

func (m *Memory) ReadWord(addr uint16) uint16 {
	if addr >= DeviceRegStart {
		return m.DeviceReadFunc(addr)
	}
	return m.mem[addr]
}

func (m *Memory) WriteWord(addr uint16, data uint16) {
	if addr >= DeviceRegStart {
		m.DeviceWriteFunc(addr, data)
		return
	}
	m.mem[addr] = data
}

func (m *Memory) WriteSegment(addr uint16, data []uint16) {
	for i, v := range data {
		m.mem[addr+uint16(i)] = v
	}
}

type MemoryWriter struct {
	mem *Memory
	// pos is in bytes, not words
	pos int
}

func NewMemoryWriter(m *Memory) *MemoryWriter {
	return &MemoryWriter{mem: m, pos: 0}
}

func (w *MemoryWriter) Write(p []byte) (int, error) {
	mem := memoryAsByteArray(w.mem)

	var err error
	bytesToWrite := len(p)
	remainingSpace := int(DeviceRegStart)*2 - w.pos
	if bytesToWrite > remainingSpace {
		bytesToWrite = remainingSpace
		err = io.ErrShortWrite
	}

	for i := 0; i < bytesToWrite; i++ {
		mem[w.pos+i] = p[i]
	}
	w.pos += bytesToWrite

	return bytesToWrite, err
}

func memoryAsByteArray(m *Memory) *[int(DeviceRegStart) * 2]byte {
	return (*[int(DeviceRegStart) * 2]byte)(unsafe.Pointer(&m.mem))
}
