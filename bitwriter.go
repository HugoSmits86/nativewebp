package nativewebp

import (
	//------------------------------
	//general
	//------------------------------
	"bytes"
	"sync"
)

type bitWriter struct {
	Buffer        *bytes.Buffer
	BitBuffer     uint64
	BitBufferSize int
}

const maxBitWriterSize = 1 << 20 // 1MB

var bitWriterPool = &sync.Pool{
	New: func() any {
		return &bitWriter{Buffer: new(bytes.Buffer)}
	},
}

func getBitWriter() *bitWriter {
	return bitWriterPool.Get().(*bitWriter)
}

func putBitWriter(w *bitWriter) {
	if w.Buffer.Cap() > maxBitWriterSize {
		return
	}

	w.Buffer.Reset()
	w.BitBuffer = 0
	w.BitBufferSize = 0

	bitWriterPool.Put(w)
}

func (w *bitWriter) writeBits(value uint64, n int) {
	if n <= 0 || n > 64 {
		panic("Invalid bit count: must be between 1 and 64")
	}

	if value >= (1 << n) {
		panic("too many bits for the given value")
	}

	w.BitBuffer |= (value << w.BitBufferSize)
	w.BitBufferSize += n
	w.writeThrough()
}

func (w *bitWriter) writeCode(code huffmanCode) {
	if code.Depth <= 0 {
		return
	}

	value := uint64(code.Bits)
	reversed := uint64(0)
	for i := 0; i < code.Depth; i++ {
		reversed = (reversed << 1) | (value & 1)
		value >>= 1
	}

	w.writeBits(reversed, code.Depth)
}

func (w *bitWriter) AlignByte() {
	w.BitBufferSize = (w.BitBufferSize + 7) &^ 7
	w.writeThrough()
}

func (w *bitWriter) writeThrough() {
	for w.BitBufferSize >= 8 {
		w.Buffer.WriteByte(byte(w.BitBuffer & 0xFF))
		w.BitBuffer >>= 8
		w.BitBufferSize -= 8
	}
}
