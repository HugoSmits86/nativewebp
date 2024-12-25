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

func TestBuildHuffmanTree(t *testing.T) {
    for id, tt := range []struct {
        histo        []int
        maxDepth    int
        expectedTree *Node // Expected structure of the Huffman tree
    }{
        // Simple case with 2 symbols
        {
            histo:     []int{5, 10},
            maxDepth: 4,
            expectedTree: &Node{
                IsBranch: true,
                Weight:   15,
                BranchLeft: &Node{
                    IsBranch: false,
                    Weight:   5,
                    Symbol:   0,
                },
                BranchRight: &Node{
                    IsBranch: false,
                    Weight:   10,
                    Symbol:   1,
                },
            },
        },
        // Histogram with more symbols
        {
            histo:     []int{5, 9, 12, 13},
            maxDepth: 5,
            expectedTree: &Node{
                IsBranch: true,
                Weight:   39,
                BranchLeft: &Node{
                    IsBranch: true,
                    Weight:   14,
                    BranchLeft: &Node{
                        IsBranch: false,
                        Weight:   5,
                        Symbol:   0,
                    },
                    BranchRight: &Node{
                        IsBranch: false,
                        Weight:   9,
                        Symbol:   1,
                    },
                },
                BranchRight: &Node{
                    IsBranch: true,
                    Weight:   25,
                    BranchLeft: &Node{
                        IsBranch: false,
                        Weight:   12,
                        Symbol:   2,
                    },
                    BranchRight: &Node{
                        IsBranch: false,
                        Weight:   13,
                        Symbol:   3,
                    },
                },
            },
        },
        // Test case that triggers the for nHeap.Len() < 1 loop
        {
            histo:     []int{}, // Empty histogram
            maxDepth: 4,
            expectedTree: &Node{
                IsBranch: false,
                Weight:   0,
                Symbol:   0,
            },
        },
        // Test case with all zero weights
        {
            histo:     []int{0, 0, 0},
            maxDepth: 4,
            expectedTree: &Node{
                IsBranch: false,
                Weight:   0,
                Symbol:   0,
            },
        },
    } {
        resultTree := buildHuffmanTree(tt.histo, tt.maxDepth)

        var compareTrees func(a, b *Node) bool
        compareTrees = func(a, b *Node) bool {
            if a == nil && b == nil {
                return true
            }
            if a == nil || b == nil {
                return false
            }
            if a.IsBranch != b.IsBranch || a.Weight != b.Weight || a.Symbol != b.Symbol {
                return false
            }
            return compareTrees(a.BranchLeft, b.BranchLeft) && compareTrees(a.BranchRight, b.BranchRight)
        }

        if !compareTrees(resultTree, tt.expectedTree) {
            t.Errorf("test %v: Huffman tree mismatch: got %+v, expected %+v", id, resultTree, tt.expectedTree)
        }
    }
}

func TestBuildHuffmanCodes(t *testing.T) {
    for id, tt := range []struct {
        histo        []int
        maxDepth    int
        expectedBits map[int]HuffmanCode // Expected results as a map for clarity
    }{
        // Test case with a single symbol
        {
            histo:     []int{10},
            maxDepth: 4,
            expectedBits: map[int]HuffmanCode{
                0: {Symbol: 0, Bits: 0, Depth: -1}, // Single symbol, no actual code assigned
            },
        },
        // Test case with two symbols
        {
            histo:     []int{5, 15},
            maxDepth: 4,
            expectedBits: map[int]HuffmanCode{
                0: {Symbol: 0, Bits: 0b0, Depth: 1}, // Symbol 0 gets code '0'
                1: {Symbol: 1, Bits: 0b1, Depth: 1}, // Symbol 1 gets code '1'
            },
        },
        // Test case with symbols requiring different depthss
        {
            histo:     []int{5, 9, 12, 13, 1}, // Fifth symbol has lower weight, longer code
            maxDepth: 4,
            expectedBits: map[int]HuffmanCode{
                0: {Symbol: 0, Bits: 0b110, Depth: 3}, // Symbol 0 gets code '110'
                1: {Symbol: 1, Bits: 0b0, Depth: 2},   // Symbol 1 gets code '0'
                2: {Symbol: 2, Bits: 0b1, Depth: 2},   // Symbol 2 gets code '1'
                3: {Symbol: 3, Bits: 0b10, Depth: 2},  // Symbol 3 gets code '10'
                4: {Symbol: 4, Bits: 0b111, Depth: 3}, // Symbol 4 gets code '111'
            },
        },
    } {
        resultCodes := buildHuffmanCodes(tt.histo, tt.maxDepth)

        for sym, expectedCode := range tt.expectedBits {
            if sym >= len(resultCodes) {
                t.Errorf("test %v: missing code for symbol %v", id, expectedCode.Symbol)
                continue
            }

            resultCode := resultCodes[sym]
            if resultCode.Bits != expectedCode.Bits || resultCode.Depth != expectedCode.Depth {
                t.Errorf("test %v: code mismatch for symbol %v: got {Bits: %b, Depth: %d}, expected {Bits: %b, Depth: %d}",
                    id, expectedCode.Symbol, resultCode.Bits, resultCode.Depth, expectedCode.Bits, expectedCode.Depth)
            }
        }
    }
}

func TestSetBitDepths(t *testing.T) {
    for id, tt := range []struct {
        tree           *Node
        expectedCodes  []HuffmanCode
    }{
        // Test case with a nil node
        {
            tree:          nil, // Nil node
            expectedCodes: []HuffmanCode{}, // No codes generated
        },
        // Test case with a single node (no branches)
        {
            tree: &Node{
                IsBranch: false,
                Weight:   5,
                Symbol:   0,
            },
            expectedCodes: []HuffmanCode{
                {Symbol: 0, Depth: 0}, // Root node has depth 0
            },
        },
        // Test case with a simple binary tree
        {
            tree: &Node{
                IsBranch: true,
                Weight:   15,
                BranchLeft: &Node{
                    IsBranch: false,
                    Weight:   5,
                    Symbol:   0,
                },
                BranchRight: &Node{
                    IsBranch: false,
                    Weight:   10,
                    Symbol:   1,
                },
            },
            expectedCodes: []HuffmanCode{
                {Symbol: 0, Depth: 1}, // Left branch depth = 1
                {Symbol: 1, Depth: 1}, // Right branch depth = 1
            },
        },
        // Test case with a more complex tree
        {
            tree: &Node{
                IsBranch: true,
                Weight:   30,
                BranchLeft: &Node{
                    IsBranch: true,
                    Weight:   15,
                    BranchLeft: &Node{
                        IsBranch: false,
                        Weight:   5,
                        Symbol:   0,
                    },
                    BranchRight: &Node{
                        IsBranch: false,
                        Weight:   10,
                        Symbol:   1,
                    },
                },
                BranchRight: &Node{
                    IsBranch: false,
                    Weight:   15,
                    Symbol:   2,
                },
            },
            expectedCodes: []HuffmanCode{
                {Symbol: 0, Depth: 2},
                {Symbol: 1, Depth: 2}, 
                {Symbol: 2, Depth: 1},
            },
        },
    } {
        var codes []HuffmanCode
        setBitDepths(tt.tree, &codes, 0)

        if len(codes) != len(tt.expectedCodes) {
            t.Errorf("test %v: depths mismatch: got %v, expected %v", id, len(codes), len(tt.expectedCodes))
            continue
        }

        for i, expectedCode := range tt.expectedCodes {
            if codes[i] != expectedCode {
                t.Errorf("test %v: mismatch at index %v: got %+v, expected %+v", id, i, codes[i], expectedCode)
            }
        }
    }
}

func TestWriteHuffmanCodes(t *testing.T) {
    for id, tt := range []struct {
        codes          []HuffmanCode
        expectedBits   []byte
        expectedBitBuf uint64
        expectedBufSize int
    }{
        // No codes present
        {
            codes: []HuffmanCode{},
            expectedBits: []byte{},
            expectedBitBuf: 0b0001,       
            expectedBufSize: 4,
        },
        // Single symbol, symbol[0] <= 1
        {
            codes: []HuffmanCode{
                {Symbol: 0, Bits: 0, Depth: 1},
            },
            expectedBits: []byte{},       
            expectedBitBuf: 0b0001,       
            expectedBufSize: 4,           
        },
        // Single symbol, symbol[0] > 1
        {
            codes: []HuffmanCode{
                {Symbol: 3, Bits: 0b11, Depth: 1},
            },
            expectedBits: []byte{0b00011101},       
            expectedBitBuf: 0b0000,       
            expectedBufSize: 3,           
        },
        // Two symbols, symbol[0] > 1
        {
            codes: []HuffmanCode{
                {Symbol: 2, Bits: 0b10, Depth: 1},
                {Symbol: 3, Bits: 0b11, Depth: 1},
            },
            expectedBits: []byte{0b00010111, 0b00011000},
            expectedBitBuf: 0b00,
            expectedBufSize: 3,    
        },
        // Write full Huffman code (trigger writeFullHuffmanCode)
        {
            codes: []HuffmanCode{
                {Symbol: 0, Bits: 0, Depth: 3},
                {Symbol: 1, Bits: 1, Depth: 3},
                {Symbol: 2, Bits: 2, Depth: 2},
            },
            expectedBits: []byte{0b00000100, 0b00000000, 0b00010010},
            expectedBitBuf: 0b0011,
            expectedBufSize: 3,
        },
    } {
        buffer := &bytes.Buffer{}
        writer := &BitWriter{
            Buffer:        buffer,
            BitBuffer:     0,
            BitBufferSize: 0,
        }

        writeHuffmanCodes(writer, tt.codes)

        if !bytes.Equal(buffer.Bytes(), tt.expectedBits) {
            t.Errorf("test %d: buffer mismatch\nexpected: %064b\n     got: %064b\n", id, tt.expectedBits, buffer.Bytes())
        }

        if writer.BitBuffer != tt.expectedBitBuf {
            t.Errorf("test %d: bit buffer mismatch\nexpected: %064b\n     got: %064b\n", id, tt.expectedBitBuf, writer.BitBuffer)
        }

        if writer.BitBufferSize != tt.expectedBufSize {
            t.Errorf("test %d: bit buffer size mismatch\nexpected: %d\n     got: %d\n", id, tt.expectedBufSize, writer.BitBufferSize)
        }
    }
}
