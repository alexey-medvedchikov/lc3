package machine

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemory_ReadWord(t *testing.T) {
	var m Memory
	v := uint16(0x0BAD)
	addr := uint16(256)
	m.mem[addr] = v
	assert.Equal(t, v, m.ReadWord(addr))
}

func TestMemory_WriteWord(t *testing.T) {
	var m Memory
	v := uint16(0x0BAD)
	addr := uint16(256)
	m.WriteWord(addr, v)
	assert.Equal(t, v, m.mem[addr])
}

func TestMemory_WriteSegment(t *testing.T) {
	var m Memory
	seg := [...]uint16{0x0102, 0x0304}
	m.WriteSegment(0, seg[:])
	assert.Equal(t, uint16(0x0102), m.mem[0])
	assert.Equal(t, uint16(0x0304), m.mem[1])
}

func TestMemoryWriter_Write_None(t *testing.T) {
	var m Memory

	w := NewMemoryWriter(&m)
	n, err := w.Write([]byte{})

	assert.Equal(t, 0, n)
	assert.NoError(t, err)
	for i := 0; i < int(DeviceRegStart); i++ {
		assert.Equal(t, uint16(0), m.ReadWord(uint16(i)))
	}
}

func TestMemoryWriter_Write(t *testing.T) {
	var m Memory

	w := NewMemoryWriter(&m)
	n, err := w.Write([]byte{0xFF, 0xFE})

	assert.Equal(t, 2, n)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0xFEFF), m.mem[0])
	for i := 1; i < int(DeviceRegStart); i++ {
		assert.Equal(t, uint16(0), m.ReadWord(uint16(i)))
	}
}

func TestMemoryWriter_Write_Twice(t *testing.T) {
	var m Memory

	w := NewMemoryWriter(&m)
	_, _ = w.Write([]byte{0xFF, 0xFE})
	n, err := w.Write([]byte{0xFD, 0xFC})

	assert.Equal(t, 2, n)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0xFEFF), m.ReadWord(0))
	assert.Equal(t, uint16(0xFCFD), m.ReadWord(1))
	for i := 2; i < int(DeviceRegStart); i++ {
		assert.Equal(t, uint16(0), m.ReadWord(uint16(i)))
	}
}

func TestMemoryWriter_Write_Overflow(t *testing.T) {
	var m Memory

	buf := make([]byte, MemorySize*2*2)
	for i := 0; i < len(buf); i++ {
		buf[i] = 0xFE
	}

	w := NewMemoryWriter(&m)
	n, err := w.Write(buf)

	assert.Equal(t, int(DeviceRegStart)*2, n)
	assert.ErrorIs(t, err, io.ErrShortWrite)
	for i := 0; i < int(DeviceRegStart); i++ {
		assert.Equal(t, uint16(0xFEFE), m.ReadWord(uint16(i)))
	}
}
