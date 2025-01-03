package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "container/heap"
    "sort"
)

type HuffmanCode struct {
    Symbol  int
    Bits    int
    Depth   int
}

type Node struct {
    IsBranch    bool
    Weight      int
    Symbol      int
    BranchLeft  *Node
    BranchRight *Node
}

type NodeHeap []*Node
func (h NodeHeap) Len() int             { return len(h) }
func (h NodeHeap) Less(i, j int) bool   { return h[i].Weight < h[j].Weight }
func (h NodeHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *NodeHeap) Push(x interface{})  { *h = append(*h, x.(*Node)) }
func (h *NodeHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

func buildHuffmanTree(histo []int, maxDepth int) *Node {
    sum := 0
    for _, x := range histo {
        sum += x
    }

    minWeight := sum >> (maxDepth - 2)

    nHeap := &NodeHeap{}
    heap.Init(nHeap)

    for s, w := range histo {
        if w > 0 {
            if w < minWeight {
                w = minWeight
            }

            heap.Push(nHeap, &Node{
                Weight: w, 
                Symbol: s,
            })
        }
    }
    
    for nHeap.Len() < 1 {
        heap.Push(nHeap, &Node{
            Weight: minWeight, 
            Symbol: 0,
        })
    }
    
    for nHeap.Len() > 1 {
        n1 := heap.Pop(nHeap).(*Node)
        n2 := heap.Pop(nHeap).(*Node)
        heap.Push(nHeap, &Node{
            IsBranch: true, 
            Weight: n1.Weight + n2.Weight, 
            BranchLeft: n1, 
            BranchRight: n2,
        })
    }

    return heap.Pop(nHeap).(*Node)
}

func buildHuffmanCodes(histo []int, maxDepth int) []HuffmanCode {
    codes := make([]HuffmanCode, len(histo))

    tree := buildHuffmanTree(histo, maxDepth)
    if !tree.IsBranch {
        codes[tree.Symbol] = HuffmanCode{tree.Symbol, 0, -1}
        return codes
    }
    
    var symbols []HuffmanCode
    setBitDepths(tree, &symbols, 0)

    sort.Slice(symbols, func(i, j int) bool {
        if symbols[i].Depth == symbols[j].Depth {
            return symbols[i].Symbol < symbols[j].Symbol
        }

        return symbols[i].Depth < symbols[j].Depth
    })

    bits := 0
    prevDepth := 0
    for _, sym := range symbols {
        bits <<= (sym.Depth - prevDepth)
        codes[sym.Symbol].Symbol = sym.Symbol
        codes[sym.Symbol].Bits = bits
        codes[sym.Symbol].Depth = sym.Depth
        bits++

        prevDepth = sym.Depth
    }

    return codes
}

func setBitDepths(node *Node, codes *[]HuffmanCode, level int) {
    if node == nil {
        return
    }

    if !node.IsBranch {
        *codes = append(*codes, HuffmanCode{
            Symbol: node.Symbol,
            Depth: level,
        })

        return
    }

    setBitDepths(node.BranchLeft, codes, level + 1)
    setBitDepths(node.BranchRight, codes, level + 1)
}

func writeHuffmanCodes(w *BitWriter, codes []HuffmanCode) {
    var symbols [2]int
    
    cnt := 0
    for _, code := range codes {
        if code.Depth != 0 {
            if cnt < 2 {
                symbols[cnt] = code.Symbol
            }

            cnt++
        }

        if cnt > 2 {
            break
        }
    }
    
    if cnt == 0 {
        w.writeBits(1, 1)
        w.writeBits(0, 3)
    } else if cnt <= 2 && symbols[0] < 1 << 8 && symbols[1] < 1 << 8 {
        w.writeBits(1, 1)
        w.writeBits(uint64(cnt - 1), 1)
        if symbols[0] <= 1 {
            w.writeBits(0, 1)
            w.writeBits(uint64(symbols[0]), 1)
        } else {
            w.writeBits(1, 1)
            w.writeBits(uint64(symbols[0]), 8)
        }

        if cnt > 1 {
            w.writeBits(uint64(symbols[1]), 8)
        }
    } else {
        writeFullHuffmanCode(w, codes)
    }
}

func writeFullHuffmanCode(w *BitWriter, codes []HuffmanCode) {
    histo := make([]int, 19)
    for _, c := range codes {
        histo[c.Depth]++
    }

    // lengthCodeOrder comes directly from the WebP specs!
    var lengthCodeOrder = []int{
        17, 18, 0, 1, 2, 3, 4, 5, 16, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
    }

    cnt := 0
    for i, c := range lengthCodeOrder {
        if histo[c] > 0 {
            cnt = max(i + 1, 4)
        }
    }

    w.writeBits(0, 1)
    w.writeBits(uint64(cnt - 4), 4)

    lenghts := buildHuffmanCodes(histo, 7)
    for i := 0; i < cnt; i++ {
        w.writeBits(uint64(lenghts[lengthCodeOrder[i]].Depth), 3)
    }

    w.writeBits(0, 1)

    for _, c := range codes {
        w.writeCode(lenghts[c.Depth])
    }
}