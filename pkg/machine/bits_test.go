package machine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_signExtend5(t *testing.T) {
	assert.Equal(t, int16(-16), signExtend5(16))
	assert.Equal(t, int16(-1), signExtend5(31))
	assert.Equal(t, int16(0), signExtend5(0))
	assert.Equal(t, int16(1), signExtend5(1))
	assert.Equal(t, int16(15), signExtend5(15))
}

func Test_signExtend6(t *testing.T) {
	assert.Equal(t, int16(-32), signExtend6(32))
	assert.Equal(t, int16(-1), signExtend6(63))
	assert.Equal(t, int16(0), signExtend6(0))
	assert.Equal(t, int16(1), signExtend6(1))
	assert.Equal(t, int16(31), signExtend6(31))
}

func Test_signExtend9(t *testing.T) {
	assert.Equal(t, int16(-256), signExtend9(256))
	assert.Equal(t, int16(-1), signExtend9(511))
	assert.Equal(t, int16(0), signExtend9(0))
	assert.Equal(t, int16(1), signExtend9(1))
	assert.Equal(t, int16(255), signExtend9(255))
}

func Test_signExtend11(t *testing.T) {
	assert.Equal(t, int16(-1024), signExtend11(1024))
	assert.Equal(t, int16(-1), signExtend11(2047))
	assert.Equal(t, int16(0), signExtend11(0))
	assert.Equal(t, int16(1), signExtend11(1))
	assert.Equal(t, int16(1023), signExtend11(1023))
}
