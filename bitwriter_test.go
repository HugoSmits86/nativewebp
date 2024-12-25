package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "bytes"
    //------------------------------
    //testing
    //------------------------------
    "testing"
)

func TestWriteBits(t *testing.T) {
    for id, tt := range []struct {
        initialBuffer   []byte
        initialBitBuf   uint64
        initialBufSize  int
        value           uint64
        bitCount        int
        expectedBuffer  []byte
        expectedBitBuf  uint64
        expectedBufSize int
        expectPanic     bool
    }{
        // Valid cases
        {nil, 0, 0, 0b1, 1, nil, 0b1, 1, false},                                                // Write 1 bit
        {nil, 0, 0, 0b11010101, 8, []byte{0b11010101}, 0, 0, false},                            // Write 8 bits, flush to buffer
        {nil, 0, 0, 0xFFFF, 16, []byte{0xFF, 0xFF}, 0, 0, false},                               // Write 16 bits, flush to buffer
        {nil, 0, 0, 0b101, 3, nil, 0b101, 3, false},                                            // Write 3 bits
        {nil, 0b1, 1, 0b10, 2, nil, 0b101, 3, false},                                           // Append 2 bits
        {nil, 0b101, 3, 0b1111, 4, nil, 0b1111101, 7, false},                                   // Append 4 bits
        {[]byte{0xFF}, 0, 0, 0b101, 3, []byte{0xFF}, 0b101, 3, false},                          // Preserve buffer
        // Multiple writes, testing flush
        {nil, 0, 0, 0b1101, 4, nil, 0b1101, 4, false},                                          // First write
        {[]byte{}, 0b1101, 4, 0b1111, 4, []byte{0xFD}, 0, 0, false},                            // Flush to buffer (8 bits)
        {[]byte{0xAB}, 0, 0, 0b1010101010101010, 16, []byte{0xAB, 0xAA, 0xAA}, 0, 0, false},    // Write 16 bits after flush
        // Invalid cases (expect panic)
        {nil, 0, 0, 0b101, 0, nil, 0, 0, true},                                                 // Bit count is 0
        {nil, 0, 0, 0b101, 65, nil, 0, 0, true},                                                // Bit count exceeds 64
        {nil, 0, 0, 0b101, 2, nil, 0, 0, true},                                                 // Value too large for bit count
    } {
        // Use defer to catch panics
        func() {
            defer func() {
                if r := recover(); r != nil {
                    if !tt.expectPanic {
                        t.Errorf("test %v: unexpected panic: %v", id, r)
                    }
                } else if tt.expectPanic {
                    t.Errorf("test %v: expected panic but did not occur", id)
                }
            }()

            buffer := &bytes.Buffer{}
            buffer.Write(tt.initialBuffer)
            writer := BitWriter{
                Buffer:        buffer,
                BitBuffer:     tt.initialBitBuf,
                BitBufferSize: tt.initialBufSize,
            }

            writer.WriteBits(tt.value, tt.bitCount)

            // Validate state
            if !tt.expectPanic {
                if !bytes.Equal(writer.Buffer.Bytes(), tt.expectedBuffer) {
                    t.Errorf("test %v: buffer mismatch: expected %v, got %v", id, tt.expectedBuffer, writer.Buffer.Bytes())
                }
                if writer.BitBuffer != tt.expectedBitBuf {
                    t.Errorf("test %v: bit buffer mismatch: expected %v, got %v", id, tt.expectedBitBuf, writer.BitBuffer)
                }
                if writer.BitBufferSize != tt.expectedBufSize {
                    t.Errorf("test %v: bit buffer size mismatch: expected %v, got %v", id, tt.expectedBufSize, writer.BitBufferSize)
                }
            }
        }()
    }
}

func TestWriteCode(t *testing.T) {
    for id, tt := range []struct {
        initialBuffer   []byte
        initialBitBuf   uint64
        initialBufSize  int
        code            HuffmanCode
        expectedBuffer  []byte
        expectedBitBuf  uint64
        expectedBufSize int
    }{
        {nil, 0, 0, HuffmanCode{Bits: 0b101, Depth: 3}, nil, 0b101, 3},                             // Basic 3-bit code
        {nil, 0, 0, HuffmanCode{Bits: 0b10, Depth: 2}, nil, 0b01, 2},                               // 2-bit code, reversed
        {nil, 0, 0, HuffmanCode{Bits: 0b1011, Depth: 4}, nil, 0b1101, 4},                           // 4-bit code, reversed
        {nil, 0b1, 1, HuffmanCode{Bits: 0b10, Depth: 2}, nil, 0b011, 3},                            // Append 2 bits to existing buffer
        {nil, 0, 0, HuffmanCode{Bits: 0, Depth: 0}, nil, 0, 0},                                     // Zero-Depth: code, no operation
        {nil, 0b10101010, 8, HuffmanCode{Bits: 0b1111, Depth: 4}, []byte{0b10101010}, 0b1111, 4},   // Flush full byte, 4 bits remaining
        {nil, 0, 0, HuffmanCode{Bits: 0b10011, Depth: 5}, nil, 0b11001, 5},                         // 5-bit code, reversed
        {nil, 0, 0, HuffmanCode{Bits: 0b1, Depth: -1}, nil, 0, 0},                                  // Negative Depth:, no operation
    } {
        buffer := &bytes.Buffer{}
        buffer.Write(tt.initialBuffer)
        writer := BitWriter{
            Buffer:        buffer,
            BitBuffer:     tt.initialBitBuf,
            BitBufferSize: tt.initialBufSize,
        }

        func() {
            defer func() {
                if r := recover(); r != nil {
                    t.Errorf("test %v: unexpected panic: %v", id, r)
                }
            }()
            writer.WriteCode(tt.code)
        }()

        if !bytes.Equal(writer.Buffer.Bytes(), tt.expectedBuffer) {
            t.Errorf("test %v: buffer mismatch: expected %v, got %v", id, tt.expectedBuffer, writer.Buffer.Bytes())
        }

        if writer.BitBuffer != tt.expectedBitBuf {
            t.Errorf("test %v: bit buffer mismatch: expected %064b, got %064b", id, tt.expectedBitBuf, writer.BitBuffer)
        }

        if writer.BitBufferSize != tt.expectedBufSize {
            t.Errorf("test %v: bit buffer size mismatch: expected %v, got %v", id, tt.expectedBufSize, writer.BitBufferSize)
        }
    }
}


func TestWriteThrough(t *testing.T) {
    for id, tt := range []struct {
        initialBuffer   []byte
        initialBitBuf   uint64
        initialBufSize  int
        expectedBuffer  []byte
        expectedBitBuf  uint64
        expectedBufSize int
    }{
        {nil, 0b11010101, 8, []byte{0b11010101}, 0, 0},                             // Exactly 8 bits
        {nil, 0b1111111111111111, 16, []byte{0xFF, 0xFF}, 0, 0},                    // Multiple of 8 bits
        {nil, 0b1010101010101010, 12, []byte{0b10101010}, 0b10101010, 4},           // More than 8 bits, remainder in buffer
        {nil, 0b11110000, 4, nil, 0b11110000, 4},                                   // Less than 8 bits, nothing flushed
        {[]byte{0xAB}, 0b11010101, 8, []byte{0xAB, 0xD5}, 0, 0},                    // Preserves existing buffer contents
        {[]byte{0xAB}, 0b1010101010101010, 12, []byte{0xAB, 0xAA}, 0b10101010, 4},  // Mixed existing buffer and partial flush
    } {
        buffer := &bytes.Buffer{}
        buffer.Write(tt.initialBuffer)
        writer := BitWriter{
            Buffer:        buffer,
            BitBuffer:     tt.initialBitBuf,
            BitBufferSize: tt.initialBufSize,
        }

        writer.writeThrough()

        if !bytes.Equal(writer.Buffer.Bytes(), tt.expectedBuffer) {
            t.Errorf("test %v: buffer mismatch: expected %v, got %v", id, tt.expectedBuffer, writer.Buffer.Bytes())
        }

        if writer.BitBuffer != tt.expectedBitBuf {
            t.Errorf("test %v: bit buffer mismatch: expected %064b, got %064b", id, tt.expectedBitBuf, writer.BitBuffer)
        }

        if writer.BitBufferSize != tt.expectedBufSize {
            t.Errorf("test %v: bit buffer size mismatch: expected %v, got %v", id, tt.expectedBufSize, writer.BitBufferSize)
        }
    }
}

func TestAlignByte(t *testing.T) {
    for id, tt := range []struct {
        initialBuffer   []byte
        initialBitBuf   uint64
        initialBufSize  int
        expectedBuffer  []byte
        expectedBitBuf  uint64
        expectedBufSize int
    }{
        {nil, 0b1101, 4, []byte{0x0D}, 0, 0},                                   // Align 4 bits, no padding
        {nil, 0b10101010, 8, []byte{0b10101010}, 0, 0},                         // Already aligned
        {nil, 0b1010101010101010, 12, []byte{0xAA, 0xAA}, 0, 0},                // Align 12 bits
        {[]byte{0xAB}, 0b1111, 4, []byte{0xAB, 0x0F}, 0, 0},                    // Existing buffer, no padding
        {[]byte{0xAB}, 0b1010101010101010, 10, []byte{0xAB, 0xAA, 0xAA}, 0, 0}, // Align 10 bits
        {nil, 0, 0, nil, 0, 0},                                                 // Empty buffer
    } {
        buffer := &bytes.Buffer{}
        buffer.Write(tt.initialBuffer)
        writer := BitWriter{
            Buffer:        buffer,
            BitBuffer:     tt.initialBitBuf,
            BitBufferSize: tt.initialBufSize,
        }

        writer.AlignByte()

        if !bytes.Equal(writer.Buffer.Bytes(), tt.expectedBuffer) {
            t.Errorf("test %v: buffer mismatch: expected %v, got %v", id, tt.expectedBuffer, writer.Buffer.Bytes())
        }

        if writer.BitBuffer != tt.expectedBitBuf {
            t.Errorf("test %v: bit buffer mismatch: expected %064b, got %064b", id, tt.expectedBitBuf, writer.BitBuffer)
        }

        if writer.BitBufferSize != tt.expectedBufSize {
            t.Errorf("test %v: bit buffer size mismatch: expected %v, got %v", id, tt.expectedBufSize, writer.BitBufferSize)
        }
    }
}
