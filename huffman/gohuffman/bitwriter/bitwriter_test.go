package bitwriter

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBitWriter(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)
	
	assert.Equal(t, byte(0), bw.buffer)
	assert.Equal(t, 0, bw.count)
	assert.Equal(t, &buf, bw.writer)
}

func TestWriteBit(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)
	
	// Write a single bit (1)
	err := bw.WriteBit(true)
	assert.NoError(t, err)
	assert.Equal(t, byte(128), bw.buffer) // 10000000
	assert.Equal(t, 1, bw.count)
	
	// Write another bit (0)
	err = bw.WriteBit(false)
	assert.NoError(t, err)
	assert.Equal(t, byte(128), bw.buffer) // 10000000
	assert.Equal(t, 2, bw.count)
	
	// Write another bit (1)
	err = bw.WriteBit(true)
	assert.NoError(t, err)
	assert.Equal(t, byte(160), bw.buffer) // 10100000
	assert.Equal(t, 3, bw.count)
}

func TestWriteBitFlushOnFullByte(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)
	
	// Write 8 bits: 10101010
	bits := []bool{true, false, true, false, true, false, true, false}
	for _, bit := range bits {
		err := bw.WriteBit(bit)
		assert.NoError(t, err)
	}
	
	// Buffer should be flushed and reset
	assert.Equal(t, byte(0), bw.buffer)
	assert.Equal(t, 0, bw.count)
	
	// Check that byte was written to buffer
	assert.Equal(t, 1, buf.Len())
	assert.Equal(t, byte(170), buf.Bytes()[0]) // 10101010 = 170
}

func TestWriteBitsFromString(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)
	
	// Write "10101010"
	err := bw.WriteBitsFromString("10101010")
	assert.NoError(t, err)
	
	// Should have flushed automatically
	assert.Equal(t, byte(0), bw.buffer)
	assert.Equal(t, 0, bw.count)
	assert.Equal(t, 1, buf.Len())
	assert.Equal(t, byte(170), buf.Bytes()[0]) // 10101010 = 170
}

func TestWriteBitsFromStringPartial(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)
	
	// Write "101" (3 bits)
	err := bw.WriteBitsFromString("101")
	assert.NoError(t, err)
	
	// Should not have flushed yet
	assert.Equal(t, byte(160), bw.buffer) // 10100000
	assert.Equal(t, 3, bw.count)
	assert.Equal(t, 0, buf.Len())
}

func TestFlush(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)
	
	// Write 3 bits: "101"
	err := bw.WriteBitsFromString("101")
	assert.NoError(t, err)
	
	// Manually flush
	err = bw.Flush()
	assert.NoError(t, err)
	
	// Buffer should be reset
	assert.Equal(t, byte(0), bw.buffer)
	assert.Equal(t, 0, bw.count)
	
	// Byte should be written with padding zeros
	assert.Equal(t, 1, buf.Len())
	assert.Equal(t, byte(160), buf.Bytes()[0]) // 10100000 = 160
}

func TestFlushEmptyBuffer(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)
	
	// Flush empty buffer
	err := bw.Flush()
	assert.NoError(t, err)
	
	// Nothing should be written
	assert.Equal(t, 0, buf.Len())
	assert.Equal(t, byte(0), bw.buffer)
	assert.Equal(t, 0, bw.count)
}

func TestMultipleBytes(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)
	
	// Write 16 bits: "1010101011001100"
	err := bw.WriteBitsFromString("1010101011001100")
	assert.NoError(t, err)
	
	// Should have written 2 bytes
	assert.Equal(t, 2, buf.Len())
	assert.Equal(t, byte(170), buf.Bytes()[0]) // 10101010 = 170
	assert.Equal(t, byte(204), buf.Bytes()[1]) // 11001100 = 204
	
	// Buffer should be empty
	assert.Equal(t, byte(0), bw.buffer)
	assert.Equal(t, 0, bw.count)
}

func TestMixedWriteOperations(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)
	
	// Write individual bits
	err := bw.WriteBit(true)
	assert.NoError(t, err)
	err = bw.WriteBit(false)
	assert.NoError(t, err)
	
	// Write string of bits
	err = bw.WriteBitsFromString("101010")
	assert.NoError(t, err)
	
	// Should have written 1 byte (8 bits total)
	assert.Equal(t, 1, buf.Len())
	assert.Equal(t, byte(170), buf.Bytes()[0]) // 10101010 = 170
}

func TestBitOrdering(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)
	
	// Write bits one by one: 1,0,0,0,0,0,0,1
	bits := []bool{true, false, false, false, false, false, false, true}
	for _, bit := range bits {
		err := bw.WriteBit(bit)
		assert.NoError(t, err)
	}
	
	// Should result in 10000001 = 129
	assert.Equal(t, 1, buf.Len())
	assert.Equal(t, byte(129), buf.Bytes()[0])
}