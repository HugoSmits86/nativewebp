## LLM Code Context Prompt

You are an AI assistant acting as a **Senior Software Engineer/Architect**. This bundle contains curated context from a software project, including its structure, key configurations, dependencies, and relevant source code.
**Note:** The context has been filtered to prioritize source code and essential configurations, excluding most build artifacts, caches, and generated files (like those typically found in `.next/` build subdirectories, `node_modules/`, `*.js.map`, `*.js.nft.json`, etc.) to provide a cleaner view for analysis.

Your primary goal is to assist the user with their software development tasks. Based on this context:
- Analyze the provided information to understand the project's architecture, conventions, and logic.
- Offer solutions and suggestions that adhere to software engineering best practices.
- When requested to write or refactor code, aim for solutions that are:
    - **Scalable:** Consider how the code will perform and adapt as the project grows.
    - **Maintainable:** Write clear, well-organized, and understandable code.
    - **Secure:** Be mindful of potential security vulnerabilities and suggest safeguards.
    - **Performant:** Optimize for efficiency where appropriate.
- Help identify potential issues, suggest improvements, or explain complex parts of the codebase.
- Keep the project's existing conventions, language features (e.g., Go, Next.js/TypeScript), and dependencies in mind.

Use this focused context to respond to the user's request accurately and efficiently.

## Project Overview

- **Generated:** 2025-06-25T19:18:32-04:00
- **Root Directory Scanned:** `/Users/khaliljouaneh/Desktop/OSS/nativewebp`
- **Filtering Mode:** Default (using ignores + .contextignore)
- **Included Files in Content:** 13
- **Included Paths (Tree):** 17
- **Ignored/Filtered Items:** 1
- **Processing Time:** 1ms

## Project Conventions (User Provided)

- **Entrypoint:** Code execution typically starts in the `cmd/` directory.
- **Core Logic:** Internal packages and business logic are primarily located in the `internal/` directory.

## Go Dependencies

### Go Modules (`go.mod`)

```go
module github.com/HugoSmits86/nativewebp

go 1.22.2

require golang.org/x/image v0.24.0

```

## Project Summary (from README.md)

*(Content from `README.md`)*

```markdown
[![Codecov Coverage](https://codecov.io/gh/HugoSmits86/nativewebp/branch/main/graph/badge.svg)](https://codecov.io/gh/HugoSmits86/nativewebp)
[![Go Reference](https://pkg.go.dev/badge/github.com/HugoSmits86/nativewebp.svg)](https://pkg.go.dev/github.com/HugoSmits86/nativewebp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# Native WebP for Go

This is a native WebP encoder written entirely in Go, with **no dependencies on libwebp** or other external libraries. Designed for performance and efficiency, this encoder generates smaller files than the standard Go PNG encoder and is approximately **50% faster** in execution.

Currently, the encoder supports only WebP lossless images (VP8L).

## Decoding Support

We provide WebP decoding through a wrapper around `golang.org/x/image/webp`, with an additional `DecodeIgnoreAlphaFlag` function to handle VP8X images where the alpha flag causes decoding issues.
## Benchmark

We conducted a quick benchmark to showcase file size reduction and encoding performance. Using an image from Google’s WebP Lossless and Alpha Gallery, we compared the results of our nativewebp encoder with the standard PNG encoder. <br/><br/>
For the PNG encoder, we applied the `png.BestCompression` setting to achieve the most competitive compression outcomes.
<br/><br/>

<table align="center">
  <tr>
    <th></th>
    <th></th>
    <th>PNG encoder</th>
    <th>nativeWebP encoder</th>
    <th>reduction</th>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/1.png" height="100px"></p></td>
    <td>file size</td>
    <td>120 kb</td>
    <td>96 kb</td>
    <td>20% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>42945049 ns/op</td>
    <td>27716447 ns/op</td>
    <td>35% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/2.png" height="100px"></p></td>
    <td>file size</td>
    <td>46 kb</td>
    <td>36 kb</td>
    <td>22% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>98509399 ns/op</td>
    <td>31461759 ns/op</td>
    <td>68% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/3.png" height="100px"></p></td>
    <td>file size</td>
    <td>236 kb</td>
    <td>194 kb</td>
    <td>18% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>178205535 ns/op</td>
    <td>102454192 ns/op</td>
    <td>43% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/4.png" height="60px"></p></td>
    <td>file size</td>
    <td>53 kb</td>
    <td>41 kb</td>
    <td>23% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>29088555 ns/op</td>
    <td>14959849 ns/op</td>
    <td>49% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/5.png" height="100px"></p></td>
    <td>file size</td>
    <td>139 kb</td>
    <td>123 kb</td>
    <td>12% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>63423995 ns/op</td>
    <td>21717392 ns/op</td>
    <td>66% faster</td>
  </tr>
</table>
<p align="center">
<sub>image source: https://developers.google.com/speed/webp/gallery2</sub>
</p>


## Installation

To install the nativewebp package, use the following command:
```Bash
go get github.com/HugoSmits86/nativewebp
```
## Usage

Here’s a simple example of how to encode an image:
```Go
file, err := os.Create(name)
if err != nil {
  log.Fatalf("Error creating file %s: %v", name, err)
}
defer file.Close()

err = nativewebp.Encode(file, img, nil)
if err != nil {
  log.Fatalf("Error encoding image to WebP: %v", err)
}
```

Here’s a simple example of how to encode an animation:
```Go
file, err := os.Create(name)
if err != nil {
  log.Fatalf("Error creating file %s: %v", name, err)
}
defer file.Close()

ani := nativewebp.Animation{
  Images: []image.Image{
    frame1,
    frame2,
  },
  Durations: []uint {
    100,
    100,
  },
  Disposals: []uint {
    0,
    0,
  },
  LoopCount: 0,
  BackgroundColor: 0xffffffff,
}

err = nativewebp.EncodeAll(file, &ani, nil)
if err != nil {
  log.Fatalf("Error encoding WebP animation: %v", err)
}
```

```

## Ignored Patterns & Exclusions Used

The following patterns/extensions were used to exclude files/directories:
```
# Defaults:
- *.bak
- *.css.map
- *.dll
- *.dylib
- *.exe
- *.js.map
- *.js.nft.json
- *.log
- *.prof
- *.pyc
- *.pyo
- *.rar
- *.so
- *.swo
- *.swp
- *.tar.gz
- *.test
- *.tf
- *.tfstate
- *.tfvars
- *.tmp
- *.zip
- .DS_Store
- .env
- .env.*.local
- .env.development
- .env.local
- .env.production
- .env.test
- .git
- .hg
- .husky/_
- .idea
- .next/app-build-manifest.json
- .next/build-manifest.json
- .next/cache
- .next/diagnostics
- .next/export-detail.txt
- .next/export-marker.txt
- .next/fallback-build-manifest.json
- .next/images-manifest.json
- .next/package.json
- .next/prerender-manifest.js
- .next/prerender-manifest.json
- .next/react-loadable-manifest.json
- .next/required-server-files.json
- .next/routes-manifest.json
- .next/server/app
- .next/server/app-paths-manifest.json
- .next/server/chunks
- .next/server/flight-manifest.js
- .next/server/flight-server-css-manifest.js
- .next/server/flight-server-css-manifest.json
- .next/server/font-loader-manifest.js
- .next/server/font-loader-manifest.json
- .next/server/interception-route-rewrite-manifest.js
- .next/server/middleware-build-manifest.js
- .next/server/middleware-manifest.json
- .next/server/next-font-manifest.js
- .next/server/next-font-manifest.json
- .next/server/pages
- .next/server/pages-manifest.json
- .next/server/server-reference-manifest.js
- .next/server/server-reference-manifest.json
- .next/server/webpack-runtime.js
- .next/static/chunks
- .next/static/development
- .next/static/media
- .next/static/runtime
- .next/static/webpack
- .next/trace
- .next/types
- .svn
- .terraform
- .venv
- .vscode
- __pycache__
- build
- coverage
- dist
- env
- examples
- instance
- node_modules
- out
- package-lock.json
- pnpm-lock.yaml
- storybook-static
- target
- terraform
- testdata
- vendor
- venv
- yarn.lock
```

## Directory Structure (Filtered)

*Filtered to show source code & key configs; build artifacts (e.g. from `.next/`, `node_modules/`) are generally excluded.*

```
.
├── .github
│   └── workflows
│       └── codecov.yml
├── LICENSE
├── README.md
├── bitwriter.go
├── bitwriter_test.go
├── go.mod
├── go.sum
├── huffman.go
├── huffman_test.go
├── reader.go
├── reader_test.go
├── transform.go
├── transform_test.go
├── writer.go
└── writer_test.go
```

## Included File Contents

### File: `.github/workflows/codecov.yml`

```yaml
name: Codecov Coverage

on:
  push:
    branches:
      - main
  pull_request:
    branches:
      - main

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.x'

      - name: Install dependencies
        run: |
          go mod tidy

      - name: Run tests with coverage
        run: |
          go test -v -coverprofile=coverage.txt ./...

      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v5
        with:
          token: ${{ secrets.CODECOV_TOKEN }}
          slug: HugoSmits86/nativewebp

```

### File: `LICENSE`

```
MIT License

Copyright (c) 2024 Hugo Smits

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```

### File: `bitwriter.go`

```go
package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "bytes"
)

type bitWriter struct {
    Buffer          *bytes.Buffer
    BitBuffer       uint64
    BitBufferSize   int
}

func (w *bitWriter) writeBits(value uint64, n int) {
    if n < 0 || n > 64 {
        panic("Invalid bit count: must be between 1 and 64")
    }

    if value >= (1 << n) {
        panic("too many bits for the given value")
    }
    
    w.BitBuffer |= (value << w.BitBufferSize)
    w.BitBufferSize += n
    w.writeThrough()
}

func (w *bitWriter) writeBytes(values []byte) {
    for _, v := range values {
        w.writeBits(uint64(v), 8)
    }
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

func (w *bitWriter) alignByte() {
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
```

### File: `bitwriter_test.go`

```go
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
        {nil, 0, 0, 0b101, -1, nil, 0, 0, true},                                                // Bit count exceeds 64
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
            writer := bitWriter{
                Buffer:        buffer,
                BitBuffer:     tt.initialBitBuf,
                BitBufferSize: tt.initialBufSize,
            }

            writer.writeBits(tt.value, tt.bitCount)

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

func TestWriteBytes(t *testing.T) {
    for id, tt := range []struct {
        initialBuffer   []byte
        initialBitBuf   uint64
        initialBufSize  int
        values          []byte
        expectedBuffer  []byte
        expectedBitBuf  uint64
        expectedBufSize int
    }{
        {nil, 0, 0, []byte{0xFF}, []byte{0xFF}, 0, 0},                      // Write single byte
        {nil, 0, 0, []byte{0x12, 0x34}, []byte{0x12, 0x34}, 0, 0},          // Write two bytes
        {[]byte{0xAB}, 0, 0, []byte{0xCD}, []byte{0xAB, 0xCD}, 0, 0},       // Preserve existing buffer
        {nil, 0b1, 1, []byte{0x80}, []byte{0x01}, 0b1, 1},                  // Partial bit buffer (1 bit) + new byte
        {[]byte{0x00}, 0b1111, 4, []byte{0x0F}, []byte{0x00, 0xFF}, 0, 4},  // Partial + full flush
        {nil, 0, 0, nil, nil, 0, 0},                                        // No values to write
    } {
        buffer := &bytes.Buffer{}
        buffer.Write(tt.initialBuffer)
        writer := bitWriter{
            Buffer:        buffer,
            BitBuffer:     tt.initialBitBuf,
            BitBufferSize: tt.initialBufSize,
        }

        writer.writeBytes(tt.values)

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

func TestWriteCode(t *testing.T) {
    for id, tt := range []struct {
        initialBuffer   []byte
        initialBitBuf   uint64
        initialBufSize  int
        code            huffmanCode
        expectedBuffer  []byte
        expectedBitBuf  uint64
        expectedBufSize int
    }{
        {nil, 0, 0, huffmanCode{Bits: 0b101, Depth: 3}, nil, 0b101, 3},                             // Basic 3-bit code
        {nil, 0, 0, huffmanCode{Bits: 0b10, Depth: 2}, nil, 0b01, 2},                               // 2-bit code, reversed
        {nil, 0, 0, huffmanCode{Bits: 0b1011, Depth: 4}, nil, 0b1101, 4},                           // 4-bit code, reversed
        {nil, 0b1, 1, huffmanCode{Bits: 0b10, Depth: 2}, nil, 0b011, 3},                            // Append 2 bits to existing buffer
        {nil, 0, 0, huffmanCode{Bits: 0, Depth: 0}, nil, 0, 0},                                     // Zero-Depth: code, no operation
        {nil, 0b10101010, 8, huffmanCode{Bits: 0b1111, Depth: 4}, []byte{0b10101010}, 0b1111, 4},   // Flush full byte, 4 bits remaining
        {nil, 0, 0, huffmanCode{Bits: 0b10011, Depth: 5}, nil, 0b11001, 5},                         // 5-bit code, reversed
        {nil, 0, 0, huffmanCode{Bits: 0b1, Depth: -1}, nil, 0, 0},                                  // Negative Depth:, no operation
    } {
        buffer := &bytes.Buffer{}
        buffer.Write(tt.initialBuffer)
        writer := bitWriter{
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
            writer.writeCode(tt.code)
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
        writer := bitWriter{
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
        writer := bitWriter{
            Buffer:        buffer,
            BitBuffer:     tt.initialBitBuf,
            BitBufferSize: tt.initialBufSize,
        }

        writer.alignByte()

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

```

### File: `go.sum`

```go
golang.org/x/image v0.24.0 h1:AN7zRgVsbvmTfNyqIbbOraYL8mSwcKncEj8ofjgzcMQ=
golang.org/x/image v0.24.0/go.mod h1:4b/ITuLfqYq1hqZcjofwctIhi7sZh2WaCjvsBNjjya8=

```

### File: `huffman.go`

```go
package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "container/heap"
    "sort"
)

type huffmanCode struct {
    Symbol  int
    Bits    int
    Depth   int
}

type node struct {
    IsBranch    bool
    Weight      int
    Symbol      int
    BranchLeft  *node
    BranchRight *node
}

type nodeHeap []*node
func (h nodeHeap) Len() int             { return len(h) }
func (h nodeHeap) Less(i, j int) bool   { return h[i].Weight < h[j].Weight }
func (h nodeHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *nodeHeap) Push(x interface{})  { *h = append(*h, x.(*node)) }
func (h *nodeHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

func buildHuffmanTree(histo []int, maxDepth int) *node {
    sum := 0
    for _, x := range histo {
        sum += x
    }

    minWeight := sum >> (maxDepth - 2)

    nHeap := &nodeHeap{}
    heap.Init(nHeap)

    for s, w := range histo {
        if w > 0 {
            if w < minWeight {
                w = minWeight
            }

            heap.Push(nHeap, &node{
                Weight: w, 
                Symbol: s,
            })
        }
    }
    
    for nHeap.Len() < 1 {
        heap.Push(nHeap, &node{
            Weight: minWeight, 
            Symbol: 0,
        })
    }
    
    for nHeap.Len() > 1 {
        n1 := heap.Pop(nHeap).(*node)
        n2 := heap.Pop(nHeap).(*node)
        heap.Push(nHeap, &node{
            IsBranch: true, 
            Weight: n1.Weight + n2.Weight, 
            BranchLeft: n1, 
            BranchRight: n2,
        })
    }

    return heap.Pop(nHeap).(*node)
}

func buildhuffmanCodes(histo []int, maxDepth int) []huffmanCode {
    codes := make([]huffmanCode, len(histo))

    tree := buildHuffmanTree(histo, maxDepth)
    if !tree.IsBranch {
        codes[tree.Symbol] = huffmanCode{tree.Symbol, 0, -1}
        return codes
    }
    
    var symbols []huffmanCode
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

func setBitDepths(node *node, codes *[]huffmanCode, level int) {
    if node == nil {
        return
    }

    if !node.IsBranch {
        *codes = append(*codes, huffmanCode{
            Symbol: node.Symbol,
            Depth: level,
        })

        return
    }

    setBitDepths(node.BranchLeft, codes, level + 1)
    setBitDepths(node.BranchRight, codes, level + 1)
}

func writehuffmanCodes(w *bitWriter, codes []huffmanCode) {
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
        writeFullhuffmanCode(w, codes)
    }
}

func writeFullhuffmanCode(w *bitWriter, codes []huffmanCode) {
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

    lengths := buildhuffmanCodes(histo, 7)
    for i := 0; i < cnt; i++ {
        w.writeBits(uint64(lengths[lengthCodeOrder[i]].Depth), 3)
    }

    w.writeBits(0, 1)

    for _, c := range codes {
        w.writeCode(lengths[c.Depth])
    }
}
```

### File: `huffman_test.go`

```go
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
        expectedTree *node // Expected structure of the Huffman tree
    }{
        // Simple case with 2 symbols
        {
            histo:     []int{5, 10},
            maxDepth: 4,
            expectedTree: &node{
                IsBranch: true,
                Weight:   15,
                BranchLeft: &node{
                    IsBranch: false,
                    Weight:   5,
                    Symbol:   0,
                },
                BranchRight: &node{
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
            expectedTree: &node{
                IsBranch: true,
                Weight:   39,
                BranchLeft: &node{
                    IsBranch: true,
                    Weight:   14,
                    BranchLeft: &node{
                        IsBranch: false,
                        Weight:   5,
                        Symbol:   0,
                    },
                    BranchRight: &node{
                        IsBranch: false,
                        Weight:   9,
                        Symbol:   1,
                    },
                },
                BranchRight: &node{
                    IsBranch: true,
                    Weight:   25,
                    BranchLeft: &node{
                        IsBranch: false,
                        Weight:   12,
                        Symbol:   2,
                    },
                    BranchRight: &node{
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
            expectedTree: &node{
                IsBranch: false,
                Weight:   0,
                Symbol:   0,
            },
        },
        // Test case with all zero weights
        {
            histo:     []int{0, 0, 0},
            maxDepth: 4,
            expectedTree: &node{
                IsBranch: false,
                Weight:   0,
                Symbol:   0,
            },
        },
    } {
        resultTree := buildHuffmanTree(tt.histo, tt.maxDepth)

        var compareTrees func(a, b *node) bool
        compareTrees = func(a, b *node) bool {
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

func TestBuildhuffmanCodes(t *testing.T) {
    for id, tt := range []struct {
        histo        []int
        maxDepth    int
        expectedBits map[int]huffmanCode // Expected results as a map for clarity
    }{
        // Test case with a single symbol
        {
            histo:     []int{10},
            maxDepth: 4,
            expectedBits: map[int]huffmanCode{
                0: {Symbol: 0, Bits: 0, Depth: -1}, // Single symbol, no actual code assigned
            },
        },
        // Test case with two symbols
        {
            histo:     []int{5, 15},
            maxDepth: 4,
            expectedBits: map[int]huffmanCode{
                0: {Symbol: 0, Bits: 0b0, Depth: 1}, // Symbol 0 gets code '0'
                1: {Symbol: 1, Bits: 0b1, Depth: 1}, // Symbol 1 gets code '1'
            },
        },
        // Test case with symbols requiring different depthss
        {
            histo:     []int{5, 9, 12, 13, 1}, // Fifth symbol has lower weight, longer code
            maxDepth: 4,
            expectedBits: map[int]huffmanCode{
                0: {Symbol: 0, Bits: 0b110, Depth: 3}, // Symbol 0 gets code '110'
                1: {Symbol: 1, Bits: 0b0, Depth: 2},   // Symbol 1 gets code '0'
                2: {Symbol: 2, Bits: 0b1, Depth: 2},   // Symbol 2 gets code '1'
                3: {Symbol: 3, Bits: 0b10, Depth: 2},  // Symbol 3 gets code '10'
                4: {Symbol: 4, Bits: 0b111, Depth: 3}, // Symbol 4 gets code '111'
            },
        },
    } {
        resultCodes := buildhuffmanCodes(tt.histo, tt.maxDepth)

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
        tree           *node
        expectedCodes  []huffmanCode
    }{
        // Test case with a nil node
        {
            tree:          nil, // Nil node
            expectedCodes: []huffmanCode{}, // No codes generated
        },
        // Test case with a single node (no branches)
        {
            tree: &node{
                IsBranch: false,
                Weight:   5,
                Symbol:   0,
            },
            expectedCodes: []huffmanCode{
                {Symbol: 0, Depth: 0}, // Root node has depth 0
            },
        },
        // Test case with a simple binary tree
        {
            tree: &node{
                IsBranch: true,
                Weight:   15,
                BranchLeft: &node{
                    IsBranch: false,
                    Weight:   5,
                    Symbol:   0,
                },
                BranchRight: &node{
                    IsBranch: false,
                    Weight:   10,
                    Symbol:   1,
                },
            },
            expectedCodes: []huffmanCode{
                {Symbol: 0, Depth: 1}, // Left branch depth = 1
                {Symbol: 1, Depth: 1}, // Right branch depth = 1
            },
        },
        // Test case with a more complex tree
        {
            tree: &node{
                IsBranch: true,
                Weight:   30,
                BranchLeft: &node{
                    IsBranch: true,
                    Weight:   15,
                    BranchLeft: &node{
                        IsBranch: false,
                        Weight:   5,
                        Symbol:   0,
                    },
                    BranchRight: &node{
                        IsBranch: false,
                        Weight:   10,
                        Symbol:   1,
                    },
                },
                BranchRight: &node{
                    IsBranch: false,
                    Weight:   15,
                    Symbol:   2,
                },
            },
            expectedCodes: []huffmanCode{
                {Symbol: 0, Depth: 2},
                {Symbol: 1, Depth: 2}, 
                {Symbol: 2, Depth: 1},
            },
        },
    } {
        var codes []huffmanCode
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

func TestWritehuffmanCodes(t *testing.T) {
    for id, tt := range []struct {
        codes          []huffmanCode
        expectedBits   []byte
        expectedBitBuf uint64
        expectedBufSize int
    }{
        // No codes present
        {
            codes: []huffmanCode{},
            expectedBits: []byte{},
            expectedBitBuf: 0b0001,       
            expectedBufSize: 4,
        },
        // Single symbol, symbol[0] <= 1
        {
            codes: []huffmanCode{
                {Symbol: 0, Bits: 0, Depth: 1},
            },
            expectedBits: []byte{},       
            expectedBitBuf: 0b0001,       
            expectedBufSize: 4,           
        },
        // Single symbol, symbol[0] > 1
        {
            codes: []huffmanCode{
                {Symbol: 3, Bits: 0b11, Depth: 1},
            },
            expectedBits: []byte{0b00011101},       
            expectedBitBuf: 0b0000,       
            expectedBufSize: 3,           
        },
        // Two symbols, symbol[0] > 1
        {
            codes: []huffmanCode{
                {Symbol: 2, Bits: 0b10, Depth: 1},
                {Symbol: 3, Bits: 0b11, Depth: 1},
            },
            expectedBits: []byte{0b00010111, 0b00011000},
            expectedBitBuf: 0b00,
            expectedBufSize: 3,    
        },
        // Write full Huffman code (trigger writeFullhuffmanCode)
        {
            codes: []huffmanCode{
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
        writer := &bitWriter{
            Buffer:        buffer,
            BitBuffer:     0,
            BitBufferSize: 0,
        }

        writehuffmanCodes(writer, tt.codes)

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

```

### File: `reader.go`

```go
package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "io"
    "bytes"
    "encoding/binary"
    //------------------------------
    //imaging
    //------------------------------
    "image"
    //------------------------------
    //errors
    //------------------------------
    decoderWebP "golang.org/x/image/webp"
)

// registers the webp decoder so image.Decode can detect and use it.
func init() {
    image.RegisterFormat("webp", "RIFF", Decode, DecodeConfig)
}

// Decode reads a WebP image from the provided io.Reader and returns it as an image.Image.
//
// This function is a wrapper around the underlying WebP decode package (golang.org/x/image/webp).
// It supports both lossy and lossless WebP formats, decoding the image accordingly.
//
// Parameters:
//   r - The source io.Reader containing the WebP encoded image.
//
// Returns:
//   The decoded image as image.Image or an error if the decoding fails.
func Decode(r io.Reader) (image.Image, error) {
    return decoderWebP.Decode(r)
}

// DecodeConfig reads the image configuration from the provided io.Reader without fully decoding the image.
//
// This function is a wrapper around the underlying WebP decode package (golang.org/x/image/webp) and
// provides access to the image's metadata, such as its dimensions and color model.
// It is useful for obtaining image information before performing a full decode.
//
// Parameters:
//   r - The source io.Reader containing the WebP encoded image.
//
// Returns:
//   An image.Config containing the image's dimensions and color model, or an error if the configuration cannot be retrieved
func DecodeConfig(r io.Reader) (image.Config, error) {
    return decoderWebP.DecodeConfig(r)
}

// DecodeIgnoreAlphaFlag reads a WebP image from the provided io.Reader and returns it as an image.Image.
//
// This function fixes x/image/webp rejecting VP8L images with the VP8X alpha flag, expecting an ALPHA chunk.  
// VP8L handles transparency internally, and the WebP spec requires the flag for transparency.
//
// This function is a wrapper around the underlying WebP decode package (golang.org/x/image/webp).
// It supports both lossy and lossless WebP formats, decoding the image accordingly.
//
// Parameters:
//   r - The source io.Reader containing the WebP encoded image.
//
// Returns:
//   The decoded image as image.Image or an error if the decoding fails.
func DecodeIgnoreAlphaFlag(r io.Reader) (image.Image, error) {
    data, err := io.ReadAll(r)
    if err != nil {
        return nil, err
    }

    if len(data) >= 30 && string(data[8:16]) == "WEBPVP8X" {
        for i := 30; i + 8 < len(data); {
            // Detect VP8L chunk, which handles transparency internally.
            // The x/image/webp package misinterprets this, so we clear the alpha flag.
            if string(data[i: i + 4]) == "VP8L" {
                flags := binary.LittleEndian.Uint32(data[20:24])
                flags &^= 0x00000010
                binary.LittleEndian.PutUint32(data[20:24], flags)
                break
            }

            i += 8 + int(binary.LittleEndian.Uint32(data[i + 4: i + 8]))
        }
    }

    return decoderWebP.Decode(bytes.NewReader(data))
}
```

### File: `reader_test.go`

```go
package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "bytes"
    "encoding/binary"
    //------------------------------
    //imaging
    //------------------------------
    "image"
    "image/color"
    "image/draw"
    //------------------------------
    //testing
    //------------------------------
    "testing"
)

func TestImageDecodeRegistration(t *testing.T) {
    img := generateTestImageNRGBA(8, 16, 64, true)
    buf := new(bytes.Buffer)

    if err := Encode(buf, img, nil); err != nil {
        t.Fatalf("Encode failed: %v", err)
    }

    img, format, err := image.Decode(buf)
    if err != nil {
        t.Errorf("image.Decode error %v", err)
        return
    }

    if format != "webp" {
        t.Errorf("expected format as webp got %v", format)
        return
    }
}

func TestDecode(t *testing.T) {
    img := generateTestImageNRGBA(8, 8, 64, true)
    
    for id, tt := range []struct {
        img                 image.Image
        UseExtendedFormat   bool
        expectedErr         string
    }{
        {
            img,
            false,
            "",
        },
        {
            img,
            true,
            "webp: invalid format",
        },
        {
            nil,    // if nil is used create a non-webp buffer
            false,
            "riff: missing RIFF chunk header",
        },
    }{

        input := new(bytes.Buffer)
        
        if tt.img != nil {
            err := Encode(input, tt.img, &Options{UseExtendedFormat: tt.UseExtendedFormat})
            if err != nil {
                t.Errorf("test %d: expected err as nil got %v", id, err)
                continue
            }
        } else {
            input.Write([]byte("not a WebP file!"))
        }

        buf := bytes.NewBuffer(input.Bytes())
        result, err := Decode(buf)

        if err == nil && tt.expectedErr != "" {
            t.Errorf("test %d: expected err as %v got nil", id, tt.expectedErr)
            continue
        } 

        if err != nil {
            if tt.expectedErr == "" {
                t.Errorf("test %d: expected err as nil got %v", id, err)
                continue
            }

            if tt.expectedErr != err.Error() {
                t.Errorf("test %d: expected err as %v got %v", id, tt.expectedErr, err)
                continue
            }

            continue
        }

        img1, ok := tt.img.(*image.NRGBA)
        if !ok {
            t.Errorf("test: unsupported image format for img1")
            return
        }

        img2 := image.NewNRGBA(result.Bounds())
        draw.Draw(img2, result.Bounds(), result, result.Bounds().Min, draw.Src)

        if !img1.Rect.Eq(img2.Rect) || img1.Stride != img2.Stride {
            t.Errorf("test %d: expected image dimensions as %v got %v", id, img1.Rect, img2.Rect)
            continue
        }
        
        if !bytes.Equal(img1.Pix, img2.Pix) {
            t.Errorf("test %d: expected image to be equal", id)
            continue
        }
    }
}

func TestDecodeConfig(t *testing.T) {
    img := generateTestImageNRGBA(8, 16, 64, true)
    buf := new(bytes.Buffer)

    err := Encode(buf, img, nil)
    if err != nil {
         t.Errorf("Encode: expected err as nil got %v", err)
         return
    }

    for id, tt := range []struct {
        input               []byte
        expectedColorModel  color.Model
        expectedWidth       int
        expectedHeight      int
        expectedErr         string
    }{
        {
            buf.Bytes(),
            color.NRGBAModel,
            8,
            16,
            "",
        },
        {
            []byte("invalid WebP data"),
            color.GrayModel,
            0,
            0,
            "riff: missing RIFF chunk header",
        },
    }{

        buf := bytes.NewBuffer(tt.input)
        result, err := DecodeConfig(buf)

        if err == nil && tt.expectedErr != "" {
            t.Errorf("test %d: expected err as %v got nil", id, tt.expectedErr)
            continue
        } 

        if err != nil {
            if tt.expectedErr == "" {
                t.Errorf("test %d: expected err as nil got %v", id, err)
                continue
            }

            if tt.expectedErr != err.Error() {
                t.Errorf("test %d: expected err as %v got %v", id, tt.expectedErr, err)
                continue
            }

            continue
        }

        if result.ColorModel != tt.expectedColorModel {
            t.Errorf("test %d: expected color model as %v got %v", id, tt.expectedColorModel, result.ColorModel)
            continue
        }

        if result.Width != tt.expectedWidth {
            t.Errorf("test %d: expected width as %v got %v", id, tt.expectedWidth, result.Width)
            continue
        }

        if result.Height != tt.expectedHeight {
            t.Errorf("test %d: expected height as %v got %v", id, tt.expectedHeight, result.Height)
            continue
        }
    }
}

func TestDecodeIgnoreAlphaFlag(t *testing.T) {
    for id, tt := range []struct {
        useExtendedFormat       bool
        useAlpha                bool
        expectedErrorDecode     string
    }{
        {
            false,
            false,
            "",
        },
        {
            false,
            true,
            "",
        },
        {
            true,
            false,
            "",
        },
        {
            true,
            true,
            "webp: invalid format",
        },
    }{
        img := generateTestImageNRGBA(8, 8, 64, tt.useAlpha)

        buf := new(bytes.Buffer)
        err := Encode(buf, img, &Options{UseExtendedFormat: tt.useExtendedFormat})
        if err != nil {
            t.Errorf("test %v: expected err as nil got %v", id, err)
            continue
        }

        // TEST A: we expect the default Decode to give an error for VP8X with Alpha flag set
        _, err = Decode(bytes.NewReader(buf.Bytes()))
        if err == nil && tt.expectedErrorDecode != "" {
            t.Errorf("test %v: expected err as %v got %v", id, tt.expectedErrorDecode, err)
            continue
        }

        if err != nil && err.Error() != tt.expectedErrorDecode {
            t.Errorf("test %v: expected err as %v got %v", id, tt.expectedErrorDecode, err)
            continue
        }
    
        // TEST B: we expect the DecodeIgnoreAlphaFlag to correctly read VP8X with Alpha flag set
        _, err = DecodeIgnoreAlphaFlag(bytes.NewReader(buf.Bytes()))
        if err != nil {
            t.Errorf("test %v: expected err as nil got %v", id, err)
            continue
        }
    }
}


func TestDecodeIgnoreAlphaFlagSearchChunk(t *testing.T) {
    img := generateTestImageNRGBA(8, 8, 64, true)

    buf := new(bytes.Buffer)
    err := Encode(buf, img, &Options{UseExtendedFormat: true})
    if err != nil {
        t.Errorf("expected err as nil got %v", err)
        return
    }

    data := buf.Bytes()
    data[20] |= 0x08 // set EXIF flag in VP8X header
    
    var exif bytes.Buffer
    exif.Write([]byte("EXIF"))
    binary.Write(&exif, binary.LittleEndian, uint32(6))
    exif.Write([]byte("Hello!"))
 
    //TEST: test what happens if VP8L is not directly after VP8X chunk
    data = append(data[:30], append(exif.Bytes(), data[30:]...)...)
    binary.LittleEndian.PutUint32(data[4: 8], uint32(len(data) - 8))
    
    _, err = DecodeIgnoreAlphaFlag(bytes.NewReader(data))
    if err != nil {
        t.Errorf("expected err as nil got %v", err)
        return
    }
}
```

### File: `transform.go`

```go
package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "math"
    "slices"
    //------------------------------
    //imaging
    //------------------------------
    "image/color"
    //------------------------------
    //errors
    //------------------------------
    //"log"
    "errors"
)

type transform int

const (
    transformPredict        = transform(0)
    transformColor          = transform(1)
    transformSubGreen       = transform(2)
    transformColorIndexing  = transform(3)     
)

func applyPredictTransform(pixels []color.NRGBA, width, height int) (int, int, int, []color.NRGBA) {
    tileBits := 4
    tileSize := 1 << tileBits
    bw := (width + tileSize - 1) / tileSize
    bh := (height + tileSize - 1) / tileSize

    blocks := make([]color.NRGBA, bw * bh)
    deltas := make([]color.NRGBA, width * height)
    
    accum := [][]int{
        make([]int, 256),
        make([]int, 256),
        make([]int, 256),
        make([]int, 256),
        make([]int, 40),
    }

    histos := make([][]int, len(accum))
    for i := range accum {
        histos[i] = make([]int, len(accum[i]))
    }

    for y := 0; y < bh; y++ {
        for x := 0; x < bw; x++ {
            mx := min((x + 1) << tileBits, width)
            my := min((y + 1) << tileBits, height)

            var best int
            var bestEntropy float64
            for i := 0; i < 14; i++ {
                for j := range accum {
                    copy(histos[j], accum[j])
                }

                for tx := x << tileBits; tx < mx; tx++ {
                    for ty := y << tileBits; ty < my; ty++ {
                        d := applyFilter(pixels, width, tx, ty, i)

                        off := ty * width + tx
                        histos[0][int(uint8(pixels[off].R - d.R))]++
                        histos[1][int(uint8(pixels[off].G - d.G))]++
                        histos[2][int(uint8(pixels[off].B - d.B))]++
                        histos[3][int(uint8(pixels[off].A - d.A))]++
                    }
                }

                var total float64
                for _, histo := range histos {
                    sum := 0
                    sumSquares := 0
                
                    for _, count := range histo {
                        sum += count
                        sumSquares += count * count
                    }
                
                    if sum == 0 {
                        continue
                    }
                
                    total += 1.0 - float64(sumSquares) / (float64(sum) * float64(sum))    
                }

                if i == 0 || total < bestEntropy {
                    bestEntropy = total
                    best = i
                }
            }

            for tx := x << tileBits; tx < mx; tx++ {
                for ty := y << tileBits; ty < my; ty++ {
                    d := applyFilter(pixels, width, tx, ty, best)
                    
                    off := ty * width + tx
                    deltas[off] = color.NRGBA{
                        R: uint8(pixels[off].R - d.R),
                        G: uint8(pixels[off].G - d.G),
                        B: uint8(pixels[off].B - d.B),
                        A: uint8(pixels[off].A - d.A),
                    }

                    accum[0][int(uint8(pixels[off].R - d.R))]++
                    accum[1][int(uint8(pixels[off].G - d.G))]++
                    accum[2][int(uint8(pixels[off].B - d.B))]++
                    accum[3][int(uint8(pixels[off].A - d.A))]++
                }
            }

            blocks[y * bw + x] = color.NRGBA{0, byte(best), 0, 255}
        }
    }
    
    copy(pixels, deltas)
    
    return tileBits, bw, bh, blocks
}

func applyFilter(pixels []color.NRGBA, width, x, y, prediction int) color.NRGBA {
    if x == 0 && y == 0 {
        return color.NRGBA{0, 0, 0, 255}
    } else if x == 0 {
        return pixels[(y - 1) * width + x]
    } else if y == 0 {
        return pixels[y * width + (x - 1)]
    }
    
    t := pixels[(y - 1) * width + x]
    l := pixels[y * width + (x - 1)]

    tl := pixels[(y - 1) * width + (x - 1)]
    tr := pixels[(y - 1) * width + (x + 1)]

    avarage2 := func(a, b color.NRGBA) color.NRGBA {
        return color.NRGBA {
            uint8((int(a.R) + int(b.R)) / 2), 
            uint8((int(a.G) + int(b.G)) / 2),  
            uint8((int(a.B) + int(b.B)) / 2),  
            uint8((int(a.A) + int(b.A)) / 2),
        }
    }

    filters := []func(t, l, tl, tr color.NRGBA) color.NRGBA {
        func(t, l, tl, tr color.NRGBA) color.NRGBA { return color.NRGBA{0, 0, 0, 255} },
        func(t, l, tl, tr color.NRGBA) color.NRGBA { return l },
        func(t, l, tl, tr color.NRGBA) color.NRGBA { return t },
        func(t, l, tl, tr color.NRGBA) color.NRGBA { return tr },
        func(t, l, tl, tr color.NRGBA) color.NRGBA { return tl },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            return avarage2(avarage2(l, tr), t)
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            return avarage2(l, tl)
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            return avarage2(l, t)
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            return avarage2(tl, t)
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            return avarage2(t, tr)
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            return avarage2(avarage2(l, tl), avarage2(t, tr))
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA { 
            pr := float64(l.R) + float64(t.R) - float64(tl.R)
            pg := float64(l.G) + float64(t.G) - float64(tl.G)
            pb := float64(l.B) + float64(t.B) - float64(tl.B)
            pa := float64(l.A) + float64(t.A) - float64(tl.A)

            // Manhattan distances to estimates for left and top pixels.
            pl := math.Abs(pa - float64(l.A)) + math.Abs(pr - float64(l.R)) + 
                  math.Abs(pg - float64(l.G)) + math.Abs(pb - float64(l.B))
            pt := math.Abs(pa - float64(t.A)) + math.Abs(pr - float64(t.R)) + 
                  math.Abs(pg - float64(t.G)) + math.Abs(pb - float64(t.B))

            if pl < pt {
                return l
            }

            return t
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            return color.NRGBA{
                uint8(max(min(int(l.R) + int(t.R) - int(tl.R), 255), 0)),
                uint8(max(min(int(l.G) + int(t.G) - int(tl.G), 255), 0)),
                uint8(max(min(int(l.B) + int(t.B) - int(tl.B), 255), 0)),
                uint8(max(min(int(l.A) + int(t.A) - int(tl.A), 255), 0)),
            }
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            a := avarage2(l, t)

            return color.NRGBA{
                uint8(max(min(int(a.R) + (int(a.R) - int(tl.R)) / 2, 255), 0)),
                uint8(max(min(int(a.G) + (int(a.G) - int(tl.G)) / 2, 255), 0)),
                uint8(max(min(int(a.B) + (int(a.B) - int(tl.B)) / 2, 255), 0)),
                uint8(max(min(int(a.A) + (int(a.A) - int(tl.A)) / 2, 255), 0)),
            }
        },
    }
    
    return filters[prediction](t, l, tl, tr)
}

func applyColorTransform(pixels []color.NRGBA, width, height int) (int, int, int, []color.NRGBA) {
    tileBits := 4
    tileSize := 1 << tileBits
    bw := (width + tileSize - 1) / tileSize
    bh := (height + tileSize - 1) / tileSize

    blocks := make([]color.NRGBA, bw * bh)
    deltas := make([]color.NRGBA, width * height)
    
    //TODO: analyze block and pick best Color transform Element (CTE)
    cte := color.NRGBA {
        R: 1,   //red to blue
        G: 2,   //green to blue
        B: 3,   //green to red
        A: 255,
    }
    
    for y := 0; y < bh; y++ {
        for x := 0; x < bw; x++ {
            mx := min((x + 1) << tileBits, width)
            my := min((y + 1) << tileBits, height)

            for tx := x << tileBits; tx < mx; tx++ {
                for ty := y << tileBits; ty < my; ty++ {
                    off := ty * width + tx

                    r := int(int8(pixels[off].R))
                    g := int(int8(pixels[off].G))
                    b := int(int8(pixels[off].B))
                
                    b -= int(int8((int16(int8(cte.G)) * int16(g)) >> 5))
                    b -= int(int8((int16(int8(cte.R)) * int16(r)) >> 5))
                    r -= int(int8((int16(int8(cte.B)) * int16(g)) >> 5))
                    
                    pixels[off].R = uint8(r & 0xff)
                    pixels[off].B = uint8(b & 0xff)

                    deltas[off] = pixels[off]
                }
            }

            blocks[y * bw + x] = cte
        }
    }
    
    copy(pixels, deltas)
    
    return tileBits, bw, bh, blocks
}

func applySubtractGreenTransform(pixels []color.NRGBA) {
    for i, _ := range pixels {
        pixels[i].R = pixels[i].R - pixels[i].G
        pixels[i].B = pixels[i].B - pixels[i].G
    }
}

func applyPaletteTransform(pixels *[]color.NRGBA, width, height int) ([]color.NRGBA, int, error) {
    var pal []color.NRGBA
    for _, p := range (*pixels) {
        if !slices.Contains(pal, p) {
            pal = append(pal, p)
        }
   
        if len(pal) > 256 {
            return nil, 0, errors.New("palette exceeds 256 colors")
        }
    }

    size := 1
    if len(pal) <= 2 {
        size = 8
    } else if len(pal) <= 4 {
        size = 4
    } else if len(pal) <= 16 {
        size = 2
    }
    
    pw := (width + size - 1) / size

    packed := make([]color.NRGBA, pw * height)
    for y := 0; y < height; y++ {
        for x := 0; x < pw; x++ {
            pack := 0
            for i := 0; i < size; i++ {
                px := x * size + i
                if px >= width {
                    break
                }

                idx := slices.Index(pal, (*pixels)[y * width + px])
                pack |= int(idx) << (i * (8 / size))
            }

            packed[y * pw + x] = color.NRGBA{G: uint8(pack), A: 255}
        }
    }

    *pixels = packed
    
    for i := len(pal) - 1; i > 0; i-- {
        pal[i] = color.NRGBA{
            R: pal[i].R - pal[i - 1].R,
            G: pal[i].G - pal[i - 1].G,
            B: pal[i].B - pal[i - 1].B,
            A: pal[i].A - pal[i - 1].A,
        }
    }

    return pal, pw, nil
}

```

### File: `transform_test.go`

```go
package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "reflect"
    "encoding/hex"
    "crypto/sha256"
    //------------------------------
    //imaging
    //------------------------------
    "image/color"
    //------------------------------
    //testing
    //------------------------------
    "testing"
)

func TestApplyPredictTransform(t *testing.T) {
    for id, tt := range []struct {
        width                   int
        height                  int
        expectedBlockWidth      int
        expectedBlockHeight     int
        expectedHash            string
        expectedBlocks          []color.NRGBA
        expectedBit             int
    }{
        {   // default case
            32,
            32,
            2,
            2,
            "d333d3e3bea7503db703dc5608240d7919b584cfa113bb655444c3547a6b8457",
            []color.NRGBA{
                {0, 4, 0, 255}, 
                {0, 4, 0, 255}, 
                {0, 4, 0, 255}, 
                {0, 4, 0, 255},
            }, 
            4,
        },
        {   // not power of 2 image res
            33,
            33,
            3,
            3,
            "a92e9e0413411cff17aec2abe8adf17c38149bd28ed3230c96ac6379e7055038",
            []color.NRGBA{
                {0, 4, 0, 255}, 
                {0, 4, 0, 255}, 
                {0, 4, 0, 255}, 
                {0, 4, 0, 255}, 
                {0, 4, 0, 255}, 
                {0, 4, 0, 255}, 
                {0, 4, 0, 255}, 
                {0, 4, 0, 255}, 
                {0, 3, 0, 255},
            }, 
            4,
        },
    }{
        img := generateTestImageNRGBA(tt.width, tt.height, 64, true)
        pixels, err := flatten(img)
        if err != nil {
            t.Errorf("test %v: unexpected error %v", id, err)
            continue
        }

        tileBit, bw, bh, blocks := applyPredictTransform(pixels, tt.width, tt.height)

        if bw != tt.expectedBlockWidth {
            t.Errorf("test %v: expected block width as %v got %v", id, tt.expectedBlockWidth, bw)
            continue
        }

        if bh != tt.expectedBlockHeight {
            t.Errorf("test %v: expected block height as %v got %v", id, tt.expectedBlockHeight, bh)
            continue
        }

        if !reflect.DeepEqual(blocks, tt.expectedBlocks) {
            t.Errorf("test %v: expected blocks as %v got %v", id, tt.expectedBlocks, blocks)
            continue
        }

        if tileBit != tt.expectedBit {
            t.Errorf("test %v: expected tile bit as %v got %v", id, tt.expectedBit, tileBit)
            continue
        }

        data := make([]byte, len(pixels) * 4)
        for j := 0; j < len(pixels); j++ {
            data[j * 4 + 0] = byte(pixels[j].R)
            data[j * 4 + 1] = byte(pixels[j].G)
            data[j * 4 + 2] = byte(pixels[j].B)
            data[j * 4 + 3] = byte(pixels[j].A)
        }

        hash := sha256.Sum256(data)
        if hex.EncodeToString(hash[:]) != tt.expectedHash {
            t.Errorf("test %v: expected hash as %v got %v", id, tt.expectedHash, hash)
            continue
        }
    }
}

func TestApplyFilter(t *testing.T) {
    pixels := []color.NRGBA{
        {R: 100, G: 100, B: 100, A: 255}, {R: 50, G: 50, B: 50, A: 255}, {R: 25, G: 25, B: 25, A: 255},
        {R: 200, G: 200, B: 200, A: 255}, {R: 75, G: 75, B: 75, A: 255}, {R: 0, G: 0, B: 0, A: 0}, 
        //added extra row for filter 11 if statement check
        {R: 100, G: 100, B: 100, A: 255}, {R: 250, G: 250, B: 250, A: 255}, {R: 225, G: 225, B: 225, A: 255},
        {R: 200, G: 200, B: 200, A: 255}, {R: 75, G: 75, B: 75, A: 255}, {R: 0, G: 0, B: 0, A: 0},
    }

    width := 3

    for id, tt := range []struct {
        prediction int
        x int
        y int
        expected   color.NRGBA
    }{
        // x y edge cases
        {prediction: 0, x: 0, y: 0, expected: color.NRGBA{R: 0, G: 0, B: 0, A: 255}},
        {prediction: 0, x: 0, y: 1, expected: color.NRGBA{R: 100, G: 100, B: 100, A: 255}},
        {prediction: 0, x: 1, y: 0, expected: color.NRGBA{R: 100, G: 100, B: 100, A: 255}},
        //filter predictions
        {prediction: 0, x: 1, y: 1, expected: color.NRGBA{R: 0, G: 0, B: 0, A: 255}},
        {prediction: 1, x: 1, y: 1, expected: color.NRGBA{R: 200, G: 200, B: 200, A: 255}},
        {prediction: 2, x: 1, y: 1, expected: color.NRGBA{R: 50, G: 50, B: 50, A: 255}},
        {prediction: 3, x: 1, y: 1, expected: color.NRGBA{R: 25, G: 25, B: 25, A: 255}},
        {prediction: 4, x: 1, y: 1, expected: color.NRGBA{R: 100, G: 100, B: 100, A: 255}},
        {prediction: 5, x: 1, y: 1, expected: color.NRGBA{R: 81, G: 81, B: 81, A: 255}},
        {prediction: 6, x: 1, y: 1, expected: color.NRGBA{R: 150, G: 150, B: 150, A: 255}},
        {prediction: 7, x: 1, y: 1, expected: color.NRGBA{R: 125, G: 125, B: 125, A: 255}},
        {prediction: 8, x: 1, y: 1, expected: color.NRGBA{R: 75, G: 75, B: 75, A: 255}},
        {prediction: 9, x: 1, y: 1, expected: color.NRGBA{R: 37, G: 37, B: 37, A: 255}},
        {prediction: 10, x: 1, y: 1, expected: color.NRGBA{R: 93, G: 93, B: 93, A: 255}},
        {prediction: 11, x: 1, y: 1, expected: color.NRGBA{R: 200, G: 200, B: 200, A: 255}},
        {prediction: 11, x: 1, y: 3, expected: color.NRGBA{R: 250, G: 250, B: 250, A: 255}}, // diff Manhattan distances
        {prediction: 12, x: 1, y: 1, expected: color.NRGBA{R: 150, G: 150, B: 150, A: 255}},
        {prediction: 13, x: 1, y: 1, expected: color.NRGBA{R: 137, G: 137, B: 137, A: 255}},
    } {
        got := applyFilter(pixels, width, tt.x, tt.y, tt.prediction)

        if !reflect.DeepEqual(got, tt.expected) {
            t.Errorf("test %d: mismatch\nexpected: %+v\n     got: %+v", id, tt.expected, got)
        }
    }
}

func TestApplyColorTransform(t *testing.T) {
    for id, tt := range []struct {
        width                   int
        height                  int
        expectedBlockWidth      int
        expectedBlockHeight     int
        expectedHash            string
        expectedBlocks          []color.NRGBA
        expectedBit             int
    }{
        {   // default case
            32,
            32,
            2,
            2,
            "7d2e490f816b7abe5f0f3dde85435a95da2a4295636cbc338689739fb1d936aa",
            []color.NRGBA{
                {1, 2, 3, 255},
                {1, 2, 3, 255},
                {1, 2, 3, 255},
                {1, 2, 3, 255},
            },
            4,
        },
        {   // non-power-of-2 dimensions
            33,
            33,
            3,
            3,
            "be8a424305cc8e044a6fbb16c2d3a14c2ece1fd2733d41f6f9b452790c22ccb8",
            []color.NRGBA{
                {1, 2, 3, 255},
                {1, 2, 3, 255},
                {1, 2, 3, 255},
                {1, 2, 3, 255},
                {1, 2, 3, 255},
                {1, 2, 3, 255},
                {1, 2, 3, 255},
                {1, 2, 3, 255},
                {1, 2, 3, 255},
            },
            4,
        },
    } {
        img := generateTestImageNRGBA(tt.width, tt.height, 128, true)
        pixels, err := flatten(img)
        if err != nil {
            t.Errorf("test %v: unexpected error %v", id, err)
            continue
        }

        tileBit, bw, bh, blocks := applyColorTransform(pixels, tt.width, tt.height)

        if bw != tt.expectedBlockWidth {
            t.Errorf("test %v: expected block width as %v got %v", id, tt.expectedBlockWidth, bw)
            continue
        }

        if bh != tt.expectedBlockHeight {
            t.Errorf("test %v: expected block height as %v got %v", id, tt.expectedBlockHeight, bh)
            continue
        }

        if !reflect.DeepEqual(blocks, tt.expectedBlocks) {
            t.Errorf("test %v: expected blocks as %v got %v", id, tt.expectedBlocks, blocks)
            continue
        }

        if tileBit != tt.expectedBit {
            t.Errorf("test %v: expected tile bit as %v got %v", id, tt.expectedBit, tileBit)
            continue
        }

        data := make([]byte, len(pixels)*4)
        for j := 0; j < len(pixels); j++ {
            data[j*4+0] = byte(pixels[j].R)
            data[j*4+1] = byte(pixels[j].G)
            data[j*4+2] = byte(pixels[j].B)
            data[j*4+3] = byte(pixels[j].A)
        }

        hash := sha256.Sum256(data)
        hashString := hex.EncodeToString(hash[:])

        if hashString != tt.expectedHash {
            t.Errorf("test %v: expected hash as %v got %v", id, tt.expectedHash, hashString)
            continue
        }
    }
}

func TestApplySubtractGreenTransform(t *testing.T) {
    for id, tt := range []struct {
        inputPixels    []color.NRGBA
        expectedPixels []color.NRGBA
    }{
        {
            inputPixels: []color.NRGBA{
                {R: 100, G: 50, B: 150},
            },
            expectedPixels: []color.NRGBA{
                {R: 50, G: 50, B: 100},
            },
        },
        {
            inputPixels: []color.NRGBA{
                {R: 200, G: 200, B: 150},
            },
            expectedPixels: []color.NRGBA{
                {R: 0, G: 200, B: 206},
            },
        },
        {
            inputPixels: []color.NRGBA{
                {R: 0, G: 128, B: 150},
            },
            expectedPixels: []color.NRGBA{
                {R: 128, G: 128, B: 22},
            },
        },
    }{
        pixels := make([]color.NRGBA, len(tt.inputPixels))
        copy(pixels, tt.inputPixels)

        applySubtractGreenTransform(pixels)

        if !reflect.DeepEqual(pixels, tt.expectedPixels) {
            t.Errorf("test %d: pixel mismatch\nexpected: %+v\n     got: %+v", id, tt.expectedPixels, pixels)
            continue
        }
    }
}

func TestApplyPaletteTransform(t *testing.T) {
    //check for too many colors error
    pixels := make([]color.NRGBA, 257)
    for i := 0; i < 257; i++ {
        pixels[i] = color.NRGBA{
            R: uint8(i % 16 * 16),
            G: uint8((i / 16) % 16 * 16),
            B: uint8((i / 256) % 16 * 16),
            A: 255,
        }
    }

    _, _, err := applyPaletteTransform(&pixels, 4, 4)

    msg := "palette exceeds 256 colors"
    if err == nil || err.Error() != msg {
        t.Errorf("test: expected error %v got %v", msg, err)
    }

    for id, tt := range []struct {
        width           int
        height          int
        pixels          []color.NRGBA
        expectedPalette []color.NRGBA
        expectedPixels  []color.NRGBA
        expectedWidth   int
    }{
        {
            //2 color pal - pack size = 8
            width: 3,
            height: 2,
            pixels: []color.NRGBA{
                {R: 255, G: 0, B: 0, A: 255}, 
                {R: 0, G: 255, B: 0, A: 255}, 
                {R: 255, G: 0, B: 0, A: 255}, 
                {R: 0, G: 255, B: 0, A: 255}, 
                {R: 255, G: 0, B: 0, A: 255}, 
                {R: 0, G: 255, B: 0, A: 255}, 
            },
            expectedPalette: []color.NRGBA{
                {R: 255, G: 0, B: 0, A: 255},
                {R: 1, G: 255, B: 0, A: 0},
            },
            expectedPixels: []color.NRGBA{
                {R: 0, G: 2, B: 0, A: 255}, 
                {R: 0, G: 5, B: 0, A: 255}, 
            },
            expectedWidth: 1,
        },
        {
            //4 color pal - pack size = 4
            width: 3,
            height: 2,
            pixels: []color.NRGBA{
                {R: 255, G: 0, B: 0, A: 255}, 
                {R: 0, G: 255, B: 0, A: 255}, 
                {R: 0, G: 0, B: 255, A: 255}, 
                {R: 255, G: 255, B: 0, A: 255}, 
                {R: 255, G: 0, B: 0, A: 255}, 
                {R: 0, G: 255, B: 0, A: 255}, 
            },
            expectedPalette: []color.NRGBA{
                {R: 255, G: 0, B: 0, A: 255},
                {R: 1, G: 255, B: 0, A: 0},
                {R: 0, G: 1, B: 255, A: 0},
                {R: 255, G: 255, B: 1, A: 0},

            },
            expectedPixels: []color.NRGBA{
                {R: 0, G: 36, B: 0, A: 255}, 
                {R: 0, G: 19, B: 0, A: 255}, 
            },
            expectedWidth: 1,
        },
        {
            //5 color pal - pack size = 2
            width: 3,
            height: 2,
            pixels: []color.NRGBA{
                {R: 255, G: 0, B: 0, A: 255}, 
                {R: 0, G: 255, B: 0, A: 255}, 
                {R: 0, G: 0, B: 255, A: 255}, 
                {R: 255, G: 255, B: 0, A: 255}, 
                {R: 255, G: 0, B: 255, A: 255}, 
                {R: 0, G: 255, B: 0, A: 255}, 
            },
            expectedPalette: []color.NRGBA{
                {R: 255, G: 0, B: 0, A: 255},
                {R: 1, G: 255, B: 0, A: 0},
                {R: 0, G: 1, B: 255, A: 0},
                {R: 255, G: 255, B: 1, A: 0},
                {R: 0, G: 1, B: 255, A: 0},
            },
            expectedPixels: []color.NRGBA{
                {R: 0, G: 16, B: 0, A: 255}, 
                {R: 0, G: 2, B: 0, A: 255}, 
                {R: 0, G: 67, B: 0, A: 255}, 
                {R: 0, G: 1, B: 0, A: 255},
            },
            expectedWidth: 2,
        },
        {
            // 16 color palette - pack size = 1
            width: 4,
            height: 5,
            pixels: []color.NRGBA{
                {R: 255, G: 0, B: 0, A: 255},   {R: 0, G: 255, B: 0, A: 255},   {R: 0, G: 0, B: 255, A: 255},   {R: 255, G: 255, B: 0, A: 255},  
                {R: 255, G: 0, B: 255, A: 255}, {R: 0, G: 255, B: 255, A: 255}, {R: 128, G: 128, B: 128, A: 255}, {R: 255, G: 128, B: 0, A: 255},
                {R: 128, G: 0, B: 255, A: 255}, {R: 255, G: 128, B: 128, A: 255}, {R: 0, G: 128, B: 128, A: 255}, {R: 128, G: 255, B: 0, A: 255}, 
                {R: 128, G: 0, B: 128, A: 255}, {R: 0, G: 128, B: 0, A: 255}, {R: 255, G: 255, B: 255, A: 255}, {R: 0, G: 0, B: 0, A: 255},
                {R: 128, G: 0, B: 128, A: 255}, {R: 0, G: 128, B: 0, A: 255}, {R: 255, G: 255, B: 255, A: 255}, {R: 0, G: 13, B: 37, A: 255},
            },
            expectedPalette: []color.NRGBA{
                {R: 255, G: 0, B: 0, A: 255},  
                {R: 1, G: 255, B: 0, A: 0},  
                {R: 0, G: 1, B: 255, A: 0},  
                {R: 255, G: 255, B: 1, A: 0},  
                {R: 0, G: 1, B: 255, A: 0},  
                {R: 1, G: 255, B: 0, A: 0},  
                {R: 128, G: 129, B: 129, A: 0},  
                {R: 127, G: 0, B: 128, A: 0},  
                {R: 129, G: 128, B: 255, A: 0},  
                {R: 127, G: 128, B: 129, A: 0},  
                {R: 1, G: 0, B: 0, A: 0},  
                {R: 128, G: 127, B: 128, A: 0},  
                {R: 0, G: 1, B: 128, A: 0},  
                {R: 128, G: 128, B: 128, A: 0},  
                {R: 255, G: 127, B: 255, A: 0},  
                {R: 1, G: 1, B: 1, A: 0},  
                {R: 0, G: 13, B: 37, A: 0},
            },
            expectedPixels: []color.NRGBA{
                {R: 0, G: 0, B: 0, A: 255},  
                {R: 0, G: 1, B: 0, A: 255},  
                {R: 0, G: 2, B: 0, A: 255},  
                {R: 0, G: 3, B: 0, A: 255},  
                {R: 0, G: 4, B: 0, A: 255},  
                {R: 0, G: 5, B: 0, A: 255},  
                {R: 0, G: 6, B: 0, A: 255},  
                {R: 0, G: 7, B: 0, A: 255},  
                {R: 0, G: 8, B: 0, A: 255},  
                {R: 0, G: 9, B: 0, A: 255},  
                {R: 0, G: 10, B: 0, A: 255},  
                {R: 0, G: 11, B: 0, A: 255},  
                {R: 0, G: 12, B: 0, A: 255},  
                {R: 0, G: 13, B: 0, A: 255},  
                {R: 0, G: 14, B: 0, A: 255},  
                {R: 0, G: 15, B: 0, A: 255},  
                {R: 0, G: 12, B: 0, A: 255},  
                {R: 0, G: 13, B: 0, A: 255},  
                {R: 0, G: 14, B: 0, A: 255},  
                {R: 0, G: 16, B: 0, A: 255}, 
            },
            expectedWidth: 4,
        },
    } {
        // Copy inputPixels to avoid modifying the test case
        pixels := make([]color.NRGBA, len(tt.pixels))
        copy(pixels, tt.pixels)

        pal, pw, err := applyPaletteTransform(&pixels, tt.width, tt.height)
        if err != nil {
            t.Errorf("test %d: unexpected error %v", id, err)
            continue
        }

        if pw != tt.expectedWidth {
            t.Errorf("test %d: expected width %v got %v", id, tt.expectedWidth, pw)
            continue
        }

        if !reflect.DeepEqual(pal, tt.expectedPalette) {
            t.Errorf("test %d: palette mismatch expected %+v got %+v", id, tt.expectedPalette, pal)
            continue
        }

        if !reflect.DeepEqual(pixels, tt.expectedPixels) {
            t.Errorf("test %d: pixel mismatch expected %+v got %+v", id, tt.expectedPixels, pixels)
            continue
        }
    }
}
```

### File: `writer.go`

```go
package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "io"
    "bytes"
    "encoding/binary"
    //------------------------------
    //imaging
    //------------------------------
    "image"
    "image/draw"
    "image/color"
    //------------------------------
    //errors
    //------------------------------
    "errors"
)

// Options holds configuration settings for WebP encoding.
//
// Currently, it provides a flag to enable the extended WebP format (VP8X),
// which allows for metadata support such as EXIF, ICC color profiles, and XMP.
//
// Fields:
//   - UseExtendedFormat: If true, wraps the VP8L frame inside a VP8X container
//     to enable metadata support. This does not affect image compression or
//     encoding itself, as VP8L remains the encoding format.
type Options struct {
    UseExtendedFormat   bool
}

// Animation holds configuration settings for WebP animations.
//
// It allows encoding a sequence of frames with individual timing and disposal options,
// supporting features like looping and background color settings.
//
// Fields:
//   - Images: A list of frames to be displayed in sequence.
//   - Durations: Timing for each frame in milliseconds, matching the Images slice.
//   - Disposals: Disposal methods for frames after display; 0 = keep, 1 = clear to background.
//   - LoopCount: Number of times the animation should repeat; 0 means infinite looping.
//   - BackgroundColor: Canvas background color in BGRA order, used for clear operations.
type Animation struct {
	Images              []image.Image
	Durations           []uint
	Disposals           []uint
	LoopCount           uint16
	BackgroundColor     uint32
}

// Encode writes the provided image.Image to the specified io.Writer in WebP format.
//
// This function always encodes the image using VP8L (lossless WebP). If `UseExtendedFormat`
// is enabled, it wraps the VP8L frame inside a VP8X container, allowing the use of metadata
// such as EXIF, ICC color profiles, or XMP metadata.
//
// Note: VP8L already supports transparency, so VP8X is **not required** for alpha support.
//
// Parameters:
//   w   - The destination writer where the encoded WebP image will be written.
//   img - The input image to be encoded.
//   o   - Pointer to Options containing encoding settings:
//         - UseExtendedFormat: If true, wraps the image in a VP8X container to enable 
//           extended WebP features like metadata.
//
// Returns:
//   An error if encoding fails or writing to the io.Writer encounters an issue.
func Encode(w io.Writer, img image.Image, o *Options) error {
    stream, hasAlpha, err := writeBitStream(img)
    if err != nil {
        return err
    }

    buf := &bytes.Buffer{}

    if o != nil && o.UseExtendedFormat {
        writeChunkVP8X(buf, img.Bounds(), hasAlpha, false)
    }

    buf.Write([]byte("VP8L"))
    binary.Write(buf, binary.LittleEndian, uint32(stream.Len()))
    buf.Write(stream.Bytes())

    w.Write([]byte("RIFF"))
    binary.Write(w, binary.LittleEndian, uint32(4 + buf.Len()))

    w.Write([]byte("WEBP"))
    w.Write(buf.Bytes())

    return nil
}

// EncodeAll writes the provided animation sequence to the specified io.Writer in WebP format.
//
// This function encodes a list of frames as a WebP animation using the VP8X container, which
// supports features like looping, frame timing, disposal methods, and background color settings.
// Each frame is individually compressed using the VP8L (lossless) format.
//
// Note: Even if `UseExtendedFormat` is not explicitly set, animations always use the VP8X container
// because it is required for WebP animation support.
//
// Parameters:
//   w   - The destination writer where the encoded WebP animation will be written.
//   ani - Pointer to Animation containing the frames and animation settings:
//         - Images: List of frames to encode.
//         - Durations: Display times for each frame in milliseconds.
//         - Disposals: Disposal methods after frame display (keep or clear).
//         - LoopCount: Number of times the animation should loop (0 = infinite).
//         - BackgroundColor: Background color for the canvas, used when clearing.
//   o   - Pointer to Options containing additional encoding settings:
//         - UseExtendedFormat: Currently unused for animations, but accepted for consistency.
//
// Returns:
//   An error if encoding fails or writing to the io.Writer encounters an issue.
func EncodeAll(w io.Writer, ani *Animation, o *Options) error {
    frames, alpha, err := writeFrames(ani)
    if err != nil {
        return err
    }

    var bounds image.Rectangle
    for _, img := range ani.Images {
        bounds.Max.X = max(img.Bounds().Max.X, bounds.Max.X)
        bounds.Max.Y = max(img.Bounds().Max.Y, bounds.Max.Y)
    }

    buf := &bytes.Buffer{}

    writeChunkVP8X(buf, bounds, alpha, true)

    buf.Write([]byte("ANIM"))
    binary.Write(buf, binary.LittleEndian, uint32(6))
    binary.Write(buf, binary.LittleEndian, uint32(ani.BackgroundColor))
    binary.Write(buf, binary.LittleEndian, uint16(ani.LoopCount))

    buf.Write(frames.Bytes())

    w.Write([]byte("RIFF"))
    binary.Write(w, binary.LittleEndian, uint32(4 + buf.Len()))

    w.Write([]byte("WEBP"))
    w.Write(buf.Bytes())

    return nil
}

func writeChunkVP8X(buf *bytes.Buffer, bounds image.Rectangle, flagAlpha, flagAni bool) {
    buf.Write([]byte("VP8X"))
    binary.Write(buf, binary.LittleEndian, uint32(10))

    var flags byte
    if flagAni {
        flags |= 1 << 1
    }

    if flagAlpha {
        flags |= 1 << 4
    }

    binary.Write(buf, binary.LittleEndian, flags)
    buf.Write([]byte{0x00, 0x00, 0x00})

    dx := bounds.Dx() - 1
    dy := bounds.Dy() - 1

    buf.Write([]byte{byte(dx), byte(dx >> 8), byte(dx >> 16)})
    buf.Write([]byte{byte(dy), byte(dy >> 8), byte(dy >> 16)})
}

func writeFrames(ani *Animation) (*bytes.Buffer, bool, error) {
    if len(ani.Images) == 0 {
        return nil, false, errors.New("must provide at least one image")
    }

    if len(ani.Images) != len(ani.Durations) {
        return nil, false, errors.New("mismatched image and durations lengths")
    }

    if len(ani.Images) != len(ani.Disposals) {
        return nil, false, errors.New("mismatched image and disposals lengths")
    }

    for i := 0; i < len(ani.Images); i++ {
        ani.Durations[i] = min(ani.Durations[i], 1 << 24 - 1)
        ani.Disposals[i] = min(ani.Disposals[i], 1)
    }

    buf := &bytes.Buffer{}
    
    var hasAlpha bool
    for i, img := range ani.Images {
        stream, alpha, err := writeBitStream(img)
        if err != nil {
            return nil, false, err
        }
    
        hasAlpha = hasAlpha || alpha

        w := &bitWriter{Buffer: buf}
        w.writeBytes([]byte("ANMF"))
        w.writeBits(uint64(16 + 8 + stream.Len()), 32)
    
        // WebP specs requires frame offsets to be divided by 2
        w.writeBits(uint64(img.Bounds().Min.X / 2), 24)
        w.writeBits(uint64(img.Bounds().Min.Y / 2), 24)
    
        w.writeBits(uint64(img.Bounds().Dx() - 1), 24)
        w.writeBits(uint64(img.Bounds().Dy() - 1), 24)
    
        w.writeBits(uint64(ani.Durations[i]), 24)
        w.writeBits(uint64(ani.Disposals[i]), 1)
        w.writeBits(uint64(0), 1)
        w.writeBits(uint64(0), 6)
    
        w.writeBytes([]byte("VP8L"))
        w.writeBits(uint64(stream.Len()), 32)
        w.Buffer.Write(stream.Bytes())
    }

    return buf, hasAlpha, nil
}

func writeBitStream(img image.Image) (*bytes.Buffer, bool, error) {
    if img == nil {
        return nil, false, errors.New("image is nil")
    }

    if img.Bounds().Dx() < 1 || img.Bounds().Dy() < 1 {
        return nil, false, errors.New("invalid image size")
    }

    if img.Bounds().Dx() > 1 << 14 || img.Bounds().Dy() > 1 << 14 {
        return nil, false, errors.New("invalid image size")
    }

    _, isIndexed := img.(*image.Paletted)

    rgba := image.NewNRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
    draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)

    b := &bytes.Buffer{}
    s := &bitWriter{Buffer: b}

    writeBitStreamHeader(s, rgba.Bounds(), !rgba.Opaque())

    var transforms [4]bool
    transforms[transformPredict] = !isIndexed
    transforms[transformColor] = false
    transforms[transformSubGreen] = !isIndexed
    transforms[transformColorIndexing] = isIndexed

    err := writeBitStreamData(s, rgba, 4, transforms)
    if err != nil {
        return nil, false, err
    }
    
    s.alignByte()

    if b.Len() % 2 != 0 {
        b.Write([]byte{0x00})
    }

    return b, !rgba.Opaque(), nil
}

func writeBitStreamHeader(w *bitWriter, bounds image.Rectangle, hasAlpha bool) {
    w.writeBits(0x2f, 8)

    w.writeBits(uint64(bounds.Dx() - 1), 14)
    w.writeBits(uint64(bounds.Dy() - 1), 14)

    if hasAlpha {
        w.writeBits(1, 1)
    } else {
        w.writeBits(0, 1)
    }

    w.writeBits(0, 3)
}

func writeBitStreamData(w *bitWriter, img image.Image, colorCacheBits int, transforms [4]bool) error {
    pixels, err := flatten(img)
    if err != nil {
        return err
    }

    width := img.Bounds().Dx()
    height := img.Bounds().Dy()

    if transforms[transformColorIndexing] {
        w.writeBits(1, 1)
        w.writeBits(3, 2)
       
        pal, pw, err := applyPaletteTransform(&pixels, width, height)
        if err != nil {
            return err
        }

        width = pw
       
        w.writeBits(uint64(len(pal) - 1), 8);
        writeImageData(w, pal, len(pal), 1, false, colorCacheBits);
    }

    if transforms[transformSubGreen] {
        w.writeBits(1, 1)
        w.writeBits(2, 2)

        applySubtractGreenTransform(pixels)
    }

    if transforms[transformColor] {
        w.writeBits(1, 1)
        w.writeBits(1, 2)

        bits, bw, bh, blocks := applyColorTransform(pixels, width, height)

        w.writeBits(uint64(bits - 2), 3);
        writeImageData(w, blocks, bw, bh, false, colorCacheBits)
    }

    if transforms[transformPredict] {
        w.writeBits(1, 1)
        w.writeBits(0, 2)

        bits, bw, bh, blocks := applyPredictTransform(pixels, width, height)

        w.writeBits(uint64(bits - 2), 3);
        writeImageData(w, blocks, bw, bh, false, colorCacheBits)
    }

    w.writeBits(0, 1) // end of transform
    writeImageData(w, pixels, width, height, true, colorCacheBits)

    return nil
}

func writeImageData(w *bitWriter, pixels []color.NRGBA, width, height int, isRecursive bool, colorCacheBits int) {
    if colorCacheBits > 0 {
        w.writeBits(1, 1)
        w.writeBits(uint64(colorCacheBits), 4) 
    } else {
        w.writeBits(0, 1)
    }

    if isRecursive {
        w.writeBits(0, 1)
    }

    encoded := encodeImageData(pixels, width, height, colorCacheBits)
    histos := computeHistograms(encoded, colorCacheBits)

    var codes [][]huffmanCode
    for i := 0; i < 5; i++ {
        // WebP specs requires Huffman codes with maximum depth of 15
        c := buildhuffmanCodes(histos[i], 15)
        codes = append(codes, c)

        writehuffmanCodes(w, c)
    }

    for i := 0; i < len(encoded); i ++ {
        w.writeCode(codes[0][encoded[i + 0]])
        if encoded[i + 0] < 256 {
            w.writeCode(codes[1][encoded[i + 1]])
            w.writeCode(codes[2][encoded[i + 2]])
            w.writeCode(codes[3][encoded[i + 3]])
            i += 3
        } else if encoded[i + 0] < 256 + 24 {
            cnt := prefixEncodeBits(int(encoded[i + 0]) - 256)
            w.writeBits(uint64(encoded[i + 1]), cnt);

            w.writeCode(codes[4][encoded[i + 2]])

            cnt = prefixEncodeBits(int(encoded[i + 2]))
            w.writeBits(uint64(encoded[i + 3]), cnt);
            i += 3
        }
    }
}

func encodeImageData(pixels []color.NRGBA, width, height, colorCacheBits int) []int {
    head := make([]int, 1 << 14)
    prev := make([]int, len(pixels))
    cache := make([]color.NRGBA, 1 << colorCacheBits)

    encoded := make([]int, len(pixels) * 4)
    cnt := 0

    var distances = []int {
        96,   73,  55,  39,  23,  13,   5,  1,  255, 255, 255, 255, 255, 255, 255, 255,
        101,  78,  58,  42,  26,  16,   8,  2,    0,   3,  9,   17,  27,  43,  59,  79,
        102,  86,  62,  46,  32,  20,  10,  6,    4,   7,  11,  21,  33,  47,  63,  87,
        105,  90,  70,  52,  37,  28,  18,  14,  12,  15,  19,  29,  38,  53,  71,  91,
        110,  99,  82,  66,  48,  35,  30,  24,  22,  25,  31,  36,  49,  67,  83, 100,
        115, 108,  94,  76,  64,  50,  44,  40,  34,  41,  45,  51,  65,  77,  95, 109,
        118, 113, 103,  92,  80,  68,  60,  56,  54,  57,  61,  69,  81,  93, 104, 114,
        119, 116, 111, 106,  97,  88,  84,  74,  72,  75,  85,  89,  98, 107, 112, 117,
    }

    for i := 0; i < len(pixels); i++ {
        if i + 2 < len(pixels) {
            h := hash(pixels[i + 0], 14)
            h ^= hash(pixels[i + 1], 14) * 0x9e3779b9
            h ^= hash(pixels[i + 2], 14) * 0x85ebca6b
            h = h % (1 << 14)

            cur := head[h] - 1
            prev[i] = head[h]
            head[h] = i + 1

            dis := 0
            streak := 0
            for j := 0; j < 8; j++ {
                // 1 << 20: sliding window size is 2^20 (1,048,576) per WebP specs.
                // 120: reserved margin for offset adjustments.
                if cur == -1 || i - cur >= 1 << 20 - 120 {
                    break
                }

                l := 0
                // Limit the maximum match length to 4096 pixels per WebP specs.
                for i + l < len(pixels) && l < 4096 {
                    if pixels[i + l] != pixels[cur + l] {
                        break
                    }
                    l++
                }

                if l > streak {
                    streak = l
                    dis = i - cur
                }

                cur = prev[cur] - 1
            }

            // Only use the match if it is at least 3 pixels long per WebP specs.
            if streak >= 3 {
                for j := 0; j < streak; j++ {
                    h := hash(pixels[i + j], colorCacheBits)
                    cache[h] = pixels[i + j]
                }
                
                y := dis / width
                x := dis - y * width
            
                code := dis + 120
                if x <= 8 && y < 8 {
                    code = distances[y * 16 + 8 - x] + 1
                } else if x > width - 8 && y < 7 {
                    code = distances[(y + 1) * 16 + 8 + (width - x)] + 1
                }

                s, l := prefixEncodeCode(streak)
                encoded[cnt + 0] = int(s + 256)
                encoded[cnt + 1] = int(l)

                s, l = prefixEncodeCode(code)
                encoded[cnt + 2] = int(s)
                encoded[cnt + 3] = int(l)
                cnt += 4
    
                i += streak - 1
                continue
            }
        }

        p := pixels[i]
        if colorCacheBits > 0 {
            hash := hash(p, colorCacheBits)

            if cache[hash] == p {
                encoded[cnt] = int(hash + 256 + 24)
                cnt++
                continue
            }

            cache[hash] = p
        }

        encoded[cnt+0] = int(p.G)
        encoded[cnt+1] = int(p.R)
        encoded[cnt+2] = int(p.B)
        encoded[cnt+3] = int(p.A)
        cnt += 4
    }

    return encoded[:cnt]
}

func prefixEncodeCode(n int) (int, int) {
    if n <= 5 {
        return max(0, n - 1), 0
    }

    shift := 0
    rem := n - 1
    for rem > 3 {
        rem >>= 1
        shift += 1
    }

    if rem == 2 {
        return 2 + 2 * shift, n - (2 << shift) - 1
    }

    return 3 + 2 * shift, n - (3 << shift) - 1
}

func prefixEncodeBits(prefix int) int {
    if prefix < 4 {
        return 0
    }

    return (prefix - 2) >> 1
}

func hash(c color.NRGBA, shifts int) uint32 {
    //hash formula including magic number 0x1e35a7bd comes directly from WebP specs!
    x := uint32(c.A) << 24 | uint32(c.R) << 16 | uint32(c.G) << 8 | uint32(c.B)
    return (x * 0x1e35a7bd) >> (32 - min(shifts, 32))
}

func computeHistograms(pixels []int, colorCacheBits int) [][]int {
    c := 0
    if colorCacheBits > 0 {
        c = 1 << colorCacheBits
    }

    histos := [][]int{
        make([]int, 256 + 24 + c),
        make([]int, 256),
        make([]int, 256),
        make([]int, 256),
        make([]int, 40),
    }

    for i := 0; i < len(pixels); i++ {
        histos[0][pixels[i]]++
        if(pixels[i] < 256) {
            histos[1][pixels[i + 1]]++
            histos[2][pixels[i + 2]]++
            histos[3][pixels[i + 3]]++
            i += 3
        } else if pixels[i] < 256 + 24 {
            histos[4][pixels[i + 2]]++
            i += 3
        }
    }

    return histos
}

func flatten(img image.Image) ([]color.NRGBA, error) {
    w := img.Bounds().Dx()
    h := img.Bounds().Dy()

    rgba, ok := img.(*image.NRGBA)
    if !ok {
        return nil, errors.New("unsupported image format")
    }

    pixels := make([]color.NRGBA, w * h)
    for y := 0; y < h; y++ {
        for x := 0; x < w; x++ {
            i := rgba.PixOffset(x, y)
            s := rgba.Pix[i : i + 4 : i + 4]

            pixels[y * w + x].R = uint8(s[0])
            pixels[y * w + x].G = uint8(s[1])
            pixels[y * w + x].B = uint8(s[2])
            pixels[y * w + x].A = uint8(s[3])
        }
    }

    return pixels, nil
}
```

### File: `writer_test.go`

```go
package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "bytes"
    "reflect"
    //------------------------------
    //imaging
    //------------------------------
    "image"
    "image/color"
    //------------------------------
    //testing
    //------------------------------
    "testing"
)

func generateTestImageNRGBA(width int, height int, brightness float64, hasAlpha bool) image.Image {
    dest := image.NewNRGBA(image.Rect(0, 0, width, height))
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            n := uint8(float64(x ^ y) * brightness)
            var c color.Color

            a := uint8(255)
            if hasAlpha {
                a = n
            }
            if y < height / 2 {
                if x < width / 2 {
                    c = color.RGBA{n, 0, 0, a}
                } else {
                    c = color.RGBA{0, n, 0, a}
                }
            } else {
                if x < width / 2 {
                    c = color.RGBA{0, 0, n, a}
                } else {
                    c = color.RGBA{n, n, 0, a}
                }
            }
            dest.Set(x, y, c)
        }
    }
    return dest
}

func TestEncodeErrors(t *testing.T) {
    for id, tt := range []struct {
        img         image.Image
        expectedMsg string
    }{
        {
            nil,
            "image is nil",
        },
        {
            image.NewNRGBA(image.Rectangle{}),
            "invalid image size",
        },
        {
            image.NewNRGBA(image.Rect(0, 0, 1 << 14 + 1, 1 << 14 + 1)),
            "invalid image size",
        },
    }{
        b := &bytes.Buffer{}

        err := Encode(b, tt.img, nil)
        if err == nil {
            t.Errorf("test %v: expected error %v got nil", id, tt.expectedMsg)
            continue
        }

        if err != nil && err.Error() != tt.expectedMsg {
            t.Errorf("test %v: expected error %v got %v", id, tt.expectedMsg, err)
            continue
        }
    }
}

func TestEncode(t *testing.T) {
    for id, tt := range []struct {
        img                 image.Image
        UseExtendedFormat   bool
        expectedBytes       []byte
    }{
        {
            generateTestImageNRGBA(8, 8, 64, true),
            false,
            []byte {
                0x52, 0x49, 0x46, 0x46, 0xd0, 0x00, 0x00, 0x00, 
                0x57, 0x45, 0x42, 0x50, 0x56, 0x50, 0x38, 0x4c, 
                0xc4, 0x00, 0x00, 0x00, 0x2f, 0x07, 0xc0, 0x01, 
                0x10, 0x8d, 0x52, 0x09, 0x22, 0xfa, 0x1f, 0x12, 
                0x06, 0x04, 0x1b, 0x89, 0x09, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x50, 0xee, 0x15, 0x00, 
                0x80, 0xb2, 0x3e, 0x37, 0x78, 0x04, 0xc8, 0x34, 
                0xeb, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x70, 0x02, 0x64, 0x9a, 0x75, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x38, 0x01, 0x32, 0xcd, 
                0x0f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0xe0, 0x04, 0x08, 0x70, 0x2e, 0xa8, 0x24, 0x55, 
                0xed, 0x0d, 0x88, 0x96, 0xf9, 0x6e, 0x56, 0x6b, 
                0xf3, 0x35, 0x1e, 0x1d, 0x7d, 0x5f, 0x38, 0xdc, 
                0x7e, 0xbc, 0x41, 0xc6, 0x5a, 0x36, 0xeb, 0x03,
            },
        },
        {
            generateTestImageNRGBA(8, 8, 64, true),
            true,
            []byte {
                0x52, 0x49, 0x46, 0x46, 0xe2, 0x00, 0x00, 0x00, 
                0x57, 0x45, 0x42, 0x50, 0x56, 0x50, 0x38, 0x58, 
                0x0a, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 
                0x07, 0x00, 0x00, 0x07, 0x00, 0x00, 0x56, 0x50, 
                0x38, 0x4c, 0xc4, 0x00, 0x00, 0x00, 0x2f, 0x07, 
                0xc0, 0x01, 0x10, 0x8d, 0x52, 0x09, 0x22, 0xfa, 
                0x1f, 0x12, 0x06, 0x04, 0x1b, 0x89, 0x09, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50, 0xee, 
                0x15, 0x00, 0x80, 0xb2, 0x3e, 0x37, 0x78, 0x04, 
                0xc8, 0x34, 0xeb, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x70, 0x02, 0x64, 0x9a, 0x75, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x38, 0x01, 
                0x32, 0xcd, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0xe0, 0x04, 0x08, 0x70, 0x2e, 0xa8, 
                0x24, 0x55, 0xed, 0x0d, 0x88, 0x96, 0xf9, 0x6e, 
                0x56, 0x6b, 0xf3, 0x35, 0x1e, 0x1d, 0x7d, 0x5f, 
                0x38, 0xdc, 0x7e, 0xbc, 0x41, 0xc6, 0x5a, 0x36, 
                0xeb, 0x03,
            },
        },
    }{
        b := &bytes.Buffer{}
        Encode(b, tt.img, &Options{UseExtendedFormat: tt.UseExtendedFormat})

        result := b.Bytes()
        
        if !bytes.Equal(result, tt.expectedBytes) {
            t.Errorf("test %v: BitStream mismatch. Got %s, expected %s", id, result, tt.expectedBytes)
        }
    }
}

func TestEncodeAllErrors(t *testing.T) {
    frame := generateTestImageNRGBA(0, 0, 64, true)

    for id, tt := range []struct {
        ani             *Animation
        expectedMsg     string
    }{
        {
            &Animation {
                Images: []image.Image{},
            },
            "must provide at least one image",
        },
        {
            &Animation {
                Images: []image.Image{
                    frame,
                },
            },
            "mismatched image and durations lengths",
        },
        {
            &Animation {
                Images: []image.Image{
                    frame,
                },
                Durations: []uint {
                    100,
                },
            },
            "mismatched image and disposals lengths",
        },
    }{
        b := &bytes.Buffer{}

        err := EncodeAll(b, tt.ani, nil)
        if err == nil {
            t.Errorf("test %v: expected error %v got nil", id, tt.expectedMsg)
            continue
        }

        if err != nil && err.Error() != tt.expectedMsg {
            t.Errorf("test %v: expected error %v got %v", id, tt.expectedMsg, err)
            continue
        }
    }
}

func TestEncodeAll(t *testing.T) {
    frame1 := generateTestImageNRGBA(4, 4, 64, true)
    frame2 := generateTestImageNRGBA(8, 8, 64, true)

    for id, tt := range []struct {
        ani                 *Animation
        expectedBytes       []byte
    }{
        {
            &Animation {
                Images: []image.Image{
                    frame1,
                },
                Durations: []uint {
                    100,
                },
                Disposals: []uint {
                    1,
                },
            },
            []byte {
                0x52, 0x49, 0x46, 0x46, 0xf0, 0x00, 0x00, 0x00, 
                0x57, 0x45, 0x42, 0x50, 0x56, 0x50, 0x38, 0x58, 
                0x0a, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x00, 
                0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 0x41, 0x4e, 
                0x49, 0x4d, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x41, 0x4e, 0x4d, 0x46, 
                0xc4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 
                0x64, 0x00, 0x00, 0x01, 0x56, 0x50, 0x38, 0x4c, 
                0xac, 0x00, 0x00, 0x00, 0x2f, 0x03, 0xc0, 0x00, 
                0x10, 0x8d, 0x52, 0x17, 0x22, 0xfa, 0x1f, 0x12, 
                0x04, 0x64, 0xd8, 0x26, 0x05, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x80, 
                0x1e, 0xd8, 0x21, 0x40, 0xa6, 0x59, 0x07, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x0b, 0x20, 
                0x92, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x0b, 0x20, 0x12, 0x03, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x04, 0x72, 0xba, 0xc0, 
                0x93, 0x87, 0xb5, 0xe5, 0xab, 0xec, 0x7e, 0x3c,
            },
        },
        {
            &Animation {
                Images: []image.Image{
                    frame1,
                    frame2,
                },
                Durations: []uint {
                    200,
                    100,
                },
                Disposals: []uint {
                    0,
                    1,
                },
            },
            []byte {
                0x52, 0x49, 0x46, 0x46, 0xd4, 0x01, 0x00, 0x00, 
                0x57, 0x45, 0x42, 0x50, 0x56, 0x50, 0x38, 0x58, 
                0x0a, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x00, 
                0x07, 0x00, 0x00, 0x07, 0x00, 0x00, 0x41, 0x4e, 
                0x49, 0x4d, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x41, 0x4e, 0x4d, 0x46, 
                0xc4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x03, 0x00, 0x00, 0x03, 0x00, 0x00, 
                0xc8, 0x00, 0x00, 0x00, 0x56, 0x50, 0x38, 0x4c, 
                0xac, 0x00, 0x00, 0x00, 0x2f, 0x03, 0xc0, 0x00, 
                0x10, 0x8d, 0x52, 0x17, 0x22, 0xfa, 0x1f, 0x12, 
                0x04, 0x64, 0xd8, 0x26, 0x05, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x80, 
                0x1e, 0xd8, 0x21, 0x40, 0xa6, 0x59, 0x07, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x0b, 0x20, 
                0x92, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x0b, 0x20, 0x12, 0x03, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x04, 0x72, 0xba, 0xc0, 
                0x93, 0x87, 0xb5, 0xe5, 0xab, 0xec, 0x7e, 0x3c, 
                0x41, 0x4e, 0x4d, 0x46, 0xdc, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07, 0x00, 
                0x00, 0x07, 0x00, 0x00, 0x64, 0x00, 0x00, 0x01, 
                0x56, 0x50, 0x38, 0x4c, 0xc4, 0x00, 0x00, 0x00, 
                0x2f, 0x07, 0xc0, 0x01, 0x10, 0x8d, 0x52, 0x09, 
                0x22, 0xfa, 0x1f, 0x12, 0x06, 0x04, 0x1b, 0x89, 
                0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x50, 0xee, 0x15, 0x00, 0x80, 0xb2, 0x3e, 0x37, 
                0x78, 0x04, 0xc8, 0x34, 0xeb, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x70, 0x02, 0x64, 0x9a, 
                0x75, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x38, 0x01, 0x32, 0xcd, 0x0f, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0xe0, 0x04, 0x08, 0x70, 
                0x2e, 0xa8, 0x24, 0x55, 0xed, 0x0d, 0x88, 0x96, 
                0xf9, 0x6e, 0x56, 0x6b, 0xf3, 0x35, 0x1e, 0x1d, 
                0x7d, 0x5f, 0x38, 0xdc, 0x7e, 0xbc, 0x41, 0xc6, 
                0x5a, 0x36, 0xeb, 0x03, 
            },
        },
    }{
        b := &bytes.Buffer{}
        EncodeAll(b, tt.ani, nil)

        result := b.Bytes()

        if !bytes.Equal(result, tt.expectedBytes) {
            t.Errorf("test %v: BitStream mismatch. Got %s, expected %s", id, result, tt.expectedBytes)
        }
    }
}

func TestWriteChunkVP8X(t *testing.T) {
    for id, tt := range []struct {
        bounds       image.Rectangle
        flagAlpha    bool
        flagAni      bool
        expectedBits []byte
    }{
        {
            bounds:   image.Rect(0, 0, 16, 16),
            flagAlpha: false,
            flagAni:   false,
            expectedBits: []byte{
                'V', 'P', '8', 'X',
                0x0a, 0x00, 0x00, 0x00,     // Chunk size (10)
                0x00,                       // Flags
                0x00, 0x00, 0x00,           // Reserved
                0x0f, 0x00, 0x00,           // Width - 1 = 15
                0x0f, 0x00, 0x00,           // Height - 1 = 15
            },
        },
        {
            bounds:   image.Rect(0, 0, 32, 32),
            flagAlpha: true,
            flagAni:   false,
            expectedBits: []byte{
                'V', 'P', '8', 'X',
                0x0a, 0x00, 0x00, 0x00,
                0x10,                       // Flags (alpha bit set)
                0x00, 0x00, 0x00,
                0x1f, 0x00, 0x00,
                0x1f, 0x00, 0x00,
            },
        },
        {
            bounds:   image.Rect(0, 0, 64, 128),
            flagAlpha: false,
            flagAni:   true,
            expectedBits: []byte{
                'V', 'P', '8', 'X',
                0x0a, 0x00, 0x00, 0x00,
                0x02,                       // Flags (animation bit set)
                0x00, 0x00, 0x00,
                0x3f, 0x00, 0x00,           // Width - 1 = 63
                0x7f, 0x00, 0x00,           // Height - 1 = 127
            },
        },
        {
            bounds:   image.Rect(0, 0, 256, 256),
            flagAlpha: true,
            flagAni:   true,
            expectedBits: []byte{
                'V', 'P', '8', 'X',
                0x0a, 0x00, 0x00, 0x00,
                0x12,                       // Flags (alpha + animation bits set)
                0x00, 0x00, 0x00,
                0xff, 0x00, 0x00,           // Width - 1 = 255
                0xff, 0x00, 0x00,           // Height - 1 = 255
            },
        },
    }{
        buffer := &bytes.Buffer{}
        writeChunkVP8X(buffer, tt.bounds, tt.flagAlpha, tt.flagAni)

        if !bytes.Equal(buffer.Bytes(), tt.expectedBits) {
            t.Errorf("test %d: buffer mismatch expected: %v got: %v\n", id, tt.expectedBits, buffer.Bytes())
            continue
        }
    }
}

func TestWriteFramesErrors(t *testing.T) {
    frame := generateTestImageNRGBA(0, 0, 64, true)

    for id, tt := range []struct {
        ani             *Animation
        expectedMsg     string
    }{
        {
            &Animation {
                Images: []image.Image{},
            },
            "must provide at least one image",
        },
        {
            &Animation {
                Images: []image.Image{
                    frame,
                },
            },
            "mismatched image and durations lengths",
        },
        {
            &Animation {
                Images: []image.Image{
                    frame,
                },
                Durations: []uint {
                    100,
                },
            },
            "mismatched image and disposals lengths",
        },
        {
            // Note: although this test is grouped with writeFrames error tests,
            // it specifically targets an error inside writeBitStream, which is called by writeFrames
            &Animation {
                Images: []image.Image{
                    frame,
                },
                Durations: []uint {
                    100,
                },
                Disposals: []uint {
                    1,
                },
            },
            "invalid image size",
        },
    }{
        _, _, err := writeFrames(tt.ani)
        if err == nil {
            t.Errorf("test %v: expected error %v got nil", id, tt.expectedMsg)
            continue
        }

        if err != nil && err.Error() != tt.expectedMsg {
            t.Errorf("test %v: expected error %v got %v", id, tt.expectedMsg, err)
            continue
        }
    }
}


func TestWriteFrames(t *testing.T) {
    frame1 := generateTestImageNRGBA(12, 12, 64, false)
    frame2 := generateTestImageNRGBA(16, 16, 64, true)

    for id, tt := range []struct {
        ani             *Animation
        expectedAlpha   bool
        expectedBits    []byte
    }{
        {
            ani: &Animation {
                Images: []image.Image{
                    frame1,
                },
                Durations: []uint {
                    100,
                },
                Disposals: []uint {
                    1,
                },
            },
            expectedAlpha: false,
            expectedBits: []byte{
                0x41, 0x4e, 0x4d, 0x46, 0xd2, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0b, 0x00, 
                0x00, 0x0b, 0x00, 0x00, 0x64, 0x00, 0x00, 0x01, 
                0x56, 0x50, 0x38, 0x4c, 0xba, 0x00, 0x00, 0x00, 
                0x2f, 0x0b, 0xc0, 0x02, 0x00, 0x8d, 0x52, 0x09, 
                0x22, 0xfa, 0x1f, 0x12, 0x06, 0x04, 0xd8, 0x86, 
                0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0xa0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0xa0, 0xc9, 0x06, 0x00, 0x20, 0x07, 0xce, 
                0xbf, 0x22, 0x40, 0xa6, 0x19, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0xe0, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x80, 0x03, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x20, 0xd3, 
                0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x80, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0xc2, 0x00, 0xe1, 0x28, 0xde, 0x8e, 0x02, 
                0x00, 0x00, 0x00, 0x04, 0x4b, 0xe3, 0x54, 0xd7, 
                0x26, 0x78, 0xed, 0xb2, 0xf0, 0xda, 0x51, 0xd7, 
                0xfa, 0x9d, 0xf2, 0x23, 0x44, 0xf7, 0x86, 0xbf, 
                0x11, 0xf4, 0x60, 0xc3, 0xa1, 0x87, 0xb9, 0x9c, 
                0x7c, 0xb0, 0x80, 0xb6, 0x14, 0xd2, 0xbe, 0xea, 
                0xa1, 0xd4, 0x38, 0x4b, 0x47, 0xac, 0x0d, 0x7f, 
                0x03, 0x00, 
            },
        },
        {
            ani: &Animation {
                Images: []image.Image{
                    frame1,
                    frame2,
                },
                Durations: []uint {
                    100,
                    20,
                },
                Disposals: []uint {
                    0,
                    0,
                },
            },
            expectedAlpha: true,
            expectedBits: []byte{
                0x41, 0x4e, 0x4d, 0x46, 0xd2, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0b, 0x00, 
                0x00, 0x0b, 0x00, 0x00, 0x64, 0x00, 0x00, 0x00, 
                0x56, 0x50, 0x38, 0x4c, 0xba, 0x00, 0x00, 0x00, 
                0x2f, 0x0b, 0xc0, 0x02, 0x00, 0x8d, 0x52, 0x09, 
                0x22, 0xfa, 0x1f, 0x12, 0x06, 0x04, 0xd8, 0x86, 
                0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0xa0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0xa0, 0xc9, 0x06, 0x00, 0x20, 0x07, 0xce, 
                0xbf, 0x22, 0x40, 0xa6, 0x19, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0xe0, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x80, 0x03, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x20, 0xd3, 
                0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x80, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0xc2, 0x00, 0xe1, 0x28, 0xde, 0x8e, 0x02, 
                0x00, 0x00, 0x00, 0x04, 0x4b, 0xe3, 0x54, 0xd7, 
                0x26, 0x78, 0xed, 0xb2, 0xf0, 0xda, 0x51, 0xd7, 
                0xfa, 0x9d, 0xf2, 0x23, 0x44, 0xf7, 0x86, 0xbf, 
                0x11, 0xf4, 0x60, 0xc3, 0xa1, 0x87, 0xb9, 0x9c, 
                0x7c, 0xb0, 0x80, 0xb6, 0x14, 0xd2, 0xbe, 0xea, 
                0xa1, 0xd4, 0x38, 0x4b, 0x47, 0xac, 0x0d, 0x7f, 
                0x03, 0x00, 0x41, 0x4e, 0x4d, 0x46, 0xfc, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x0f, 0x00, 0x00, 0x0f, 0x00, 0x00, 0x14, 0x00, 
                0x00, 0x00, 0x56, 0x50, 0x38, 0x4c, 0xe4, 0x00, 
                0x00, 0x00, 0x2f, 0x0f, 0xc0, 0x03, 0x10, 0x8d, 
                0x52, 0x09, 0x22, 0xfa, 0x1f, 0x12, 0x06, 0x04, 
                0x1b, 0x89, 0xc9, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x40, 0xf9, 0xff, 0x3a, 0x36, 0x00, 
                0x38, 0xd6, 0xe7, 0x07, 0x43, 0x80, 0x4c, 0xb3, 
                0x0e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x27, 0x40, 0xa6, 0x59, 0x07, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x80, 0x13, 0x20, 0xd3, 0xfc, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x0c, 0x10, 0x68, 0x83, 0x0f, 0xe1, 0x09, 0x00, 
                0x00, 0x00, 0xce, 0x01, 0x21, 0x26, 0x57, 0x49, 
                0xae, 0xdf, 0xe0, 0xbb, 0xa1, 0x65, 0xb5, 0x32, 
                0x41, 0xd6, 0x3b, 0x5a, 0x33, 0x41, 0xca, 0xf9, 
                0xcd, 0x62, 0xef, 0xce, 0x03, 0x98, 0x1e, 0x7d, 
                0x6f, 0x4f, 0xce, 0xc5, 0xed, 0xeb, 0xc1, 0x71, 
                0x66, 0x93, 0x55, 0x76, 0xcb, 0x56, 0x76, 0xfb, 
                0x20, 0x65, 0xf7, 0xfd, 0x18, 0x00,
            },
        },
    }{
        
        buffer, alpha, err := writeFrames(tt.ani)
        if err != nil {
            t.Errorf("test %v: unexpected error %v", id, err)
            continue
        }

        if alpha != tt.expectedAlpha {
            t.Errorf("test %v: expected alpha as %v got %v", id, tt.expectedAlpha, alpha)
            continue
        }
        
        if !bytes.Equal(buffer.Bytes(), tt.expectedBits) {
            t.Errorf("test %d: buffer mismatch expected: %v got: %v\n", id, tt.expectedBits, buffer.Bytes())
            continue
        }
    }
}

func TestWriteBitStreamHeader(t *testing.T) {
    for id, tt := range []struct {
        bounds       image.Rectangle
        hasAlpha     bool
        expectedBits []byte
    }{
        // Test case with no alpha channel
        {
            bounds:   image.Rect(0, 0, 16, 16),
            hasAlpha: false,
            expectedBits: []byte{
                0x2f,       // Header prefix
                0x0f, 0xc0, // Width - 1 (14 bits: 15) + first 6 bits of Height - 1
                0x03, 0x00, // Remaining bits of Height - 1 (14 bits: 15) + no alpha + padding
            },
        },
        // Test case with alpha channel
        {
            bounds:   image.Rect(0, 0, 32, 32),
            hasAlpha: true,
            expectedBits: []byte{
                0x2f,       // Header prefix
                0x1f, 0xc0, // Width - 1 (14 bits: 31) + first 6 bits of Height - 1
                0x07, 0x10, // Remaining bits of Height - 1 (14 bits: 31) + alpha + padding
            },
        },
        // Larger rectangle with no alpha
        {
            bounds:   image.Rect(0, 0, 128, 64),
            hasAlpha: false,
            expectedBits: []byte{
                0x2f,       // Header prefix
                0x7f, 0xc0, // Width - 1 (14 bits: 127) + first 6 bits of Height - 1
                0x0f, 0x00, // Remaining bits of Height - 1 (14 bits: 63) + no alpha + padding
            },
        },
    }{
        buffer := &bytes.Buffer{}
        writer := &bitWriter{
            Buffer:        buffer,
            BitBuffer:     0,
            BitBufferSize: 0,
        }

        writeBitStreamHeader(writer, tt.bounds, tt.hasAlpha)

        if !bytes.Equal(buffer.Bytes(), tt.expectedBits) {
            t.Errorf("test %d: buffer mismatch expected: %v got: %v\n", id, tt.expectedBits, buffer.Bytes())
            continue
        }
    }
}

func TestWritBitStreamErrors(t *testing.T) {
    for id, tt := range []struct {
        img         image.Image
        expectedMsg string
    }{
        {
            nil,
            "image is nil",
        },
        {
            image.NewNRGBA(image.Rectangle{}),
            "invalid image size",
        },
        {
            image.NewNRGBA(image.Rect(0, 0, 1 << 14 + 1, 1 << 14 + 1)),
            "invalid image size",
        },
    }{
        _, _, err := writeBitStream(tt.img)
        if err == nil {
            t.Errorf("test %v: expected error %v got nil", id, tt.expectedMsg)
            continue
        }

        if err != nil && err.Error() != tt.expectedMsg {
            t.Errorf("test %v: expected error %v got %v", id, tt.expectedMsg, err)
            continue
        }
    }
}

func TestWriteBitStream(t *testing.T) {
    for id, tt := range []struct {
        img                 image.Image
        expectedAlpha       bool
        expectedBytes       []byte
    }{
        {
            generateTestImageNRGBA(8, 8, 64, true),
            true,
            []byte {
                0x2f, 0x07, 0xc0, 0x01, 0x10, 0x8d, 0x52, 0x09, 
                0x22, 0xfa, 0x1f, 0x12, 0x06, 0x04, 0x1b, 0x89, 
                0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x50, 0xee, 0x15, 0x00, 0x80, 0xb2, 0x3e, 0x37, 
                0x78, 0x04, 0xc8, 0x34, 0xeb, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x70, 0x02, 0x64, 0x9a, 
                0x75, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x38, 0x01, 0x32, 0xcd, 0x0f, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0xe0, 0x04, 0x08, 0x70, 
                0x2e, 0xa8, 0x24, 0x55, 0xed, 0x0d, 0x88, 0x96, 
                0xf9, 0x6e, 0x56, 0x6b, 0xf3, 0x35, 0x1e, 0x1d, 
                0x7d, 0x5f, 0x38, 0xdc, 0x7e, 0xbc, 0x41, 0xc6, 
                0x5a, 0x36, 0xeb, 0x03,
            },
        },
        {
            generateTestImageNRGBA(8, 8, 64, false),
            false,
            []byte {
                0x2f, 0x07, 0xc0, 0x01, 0x00, 0x8d, 0x52, 0x09, 
                0x22, 0xfa, 0x1f, 0x12, 0x04, 0x04, 0xdb, 0xa6, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 
                0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x0a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x80, 0x0f, 0x00, 0x00, 0x72, 0xe0, 0x58, 0x87, 
                0x00, 0x99, 0x66, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x80, 0x03, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x1c, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x40, 0x80, 0x4c, 0x13, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xe0, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x01, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0e, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 
                0x25, 0x40, 0xeb, 0xaa, 0x3f, 0x59, 0x54, 0x94, 
                0xe3, 0xb0, 0x21, 0x66, 0x98, 0xee, 0x58, 0xa7, 
                0x4d, 0xd1, 0x8f, 0x2f, 0x1e, 0x19, 0x82, 0x12, 
                0x86, 0xff, 0x78, 0xd6, 0xb2, 0xdd, 0xd1, 0x0d,
            },
        },
    }{
        b, alpha, err := writeBitStream(tt.img)
        if err != nil {
            t.Errorf("test %v: unexpected error %v", id, err)
            continue
        }

        if alpha != tt.expectedAlpha {
            t.Errorf("test %v: expected alpha as %v got %v", id, tt.expectedAlpha, alpha)
            continue 
        }

        result := b.Bytes()

        if !bytes.Equal(result, tt.expectedBytes) {
            t.Errorf("test %v: BitStream mismatch. Got %s, expected %s", id, result, tt.expectedBytes)
            continue
        }
    }
}

func TestWriteBitStreamDataErrors(t *testing.T) {
    imgpal := image.NewNRGBA(image.Rect(0, 0, 257, 1))
    for i := 0; i < 257; i++ {
        imgpal.Set(i, 0, color.NRGBA{
            R: uint8(i % 16 * 16),
            G: uint8((i / 16) % 16 * 16),
            B: uint8((i / 256) % 16 * 16),
            A: 255,
        })
    }

    for id, tt := range []struct {
        img         image.Image
        transforms  [4]bool
        expectedMsg string
    }{
        {
            image.NewRGBA(image.Rectangle{}),
            [4]bool{ false, false, false, false, },
            "unsupported image format",
        },
        {
            imgpal,
            [4]bool{ false, false, false, true, },
            "palette exceeds 256 colors",
        },
    }{
        b := &bytes.Buffer{}
        s := &bitWriter{Buffer: b}

        err := writeBitStreamData(s, tt.img, 0, tt.transforms)
        if err == nil {
            t.Errorf("test %v: expected error %v got nil", id, tt.expectedMsg)
            continue
        }

        if err != nil && err.Error() != tt.expectedMsg {
            t.Errorf("test %v: expected error %v got %v", id, tt.expectedMsg, err)
            continue
        }
    }
}

func TestWriteBitStreamData(t *testing.T) {
    img := generateTestImageNRGBA(8, 8, 64, true)

    for id, tt := range []struct {
        transforms      [4]bool
        colorCacheBits  int
        expectedBytes   []byte
    }{
        { 
            [4]bool{
                false,  //transformPredict
                false,  //transformColor
                true,   //transformSubGreen
                false,  //transformColorIndexing
            },
            0,
            []byte{
                0x85, 0x00, 0x22, 0x09, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x98, 0x01, 0x00, 0x80, 0x00, 
                0x22, 0x69, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0xb0, 0x00, 0x22, 0x69, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0xb0, 0x00, 0x82, 0x08, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xe8, 
                0x02, 0x30, 0x2d, 0x1b, 0x54, 0x56, 0x55, 0x9b, 
                0xc0, 0xb6, 0x2a, 0x41, 0x75, 0xd5, 0x6a, 0x56, 
                0x55, 0x83, 0x4a, 0xdb, 0x32, 0xd7, 0x4a, 0x00, 
                0x58, 0x8e, 0x07, 0xc9, 0x54, 0x9a, 0x05, 0x3c, 
                0x97, 0x04, 0xe9, 0xd4, 0xca, 0xa6, 0xd2, 0x20, 
                0xc9, 0x73, 0xec, 0x9a, 0x04,                
            },
        },
        {
            [4]bool{
                false,  //transformPredict
                false,  //transformColor
                true,   //transformSubGreen
                false,  //transformColorIndexing
            },
            8,
            []byte{
                0x15, 0x21, 0x20, 0xd8, 0x36, 0x03, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0xca, 0x00, 0x00, 
                0x40, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x1c, 0x66, 0x00, 0x00, 0x00, 
                0x00, 0x03, 0x00, 0x00, 0x80, 0x3b, 0x00, 0x00, 
                0x00, 0xc0, 0x01, 0x00, 0x00, 0x83, 0x3b, 0x00, 
                0x00, 0x00, 0xc0, 0x01, 0x02, 0x88, 0xa4, 0x01, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x02, 
                0x88, 0xe4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0xc0, 0x02, 0x88, 0x04, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x60, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x80, 0x01, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x5d, 0xc0, 0x2e, 
                0x3b, 0x76, 0xa3, 0xa4, 0x50, 0x0e, 0xee, 0x2a, 
                0xfe, 0x71, 0x8d, 0xf3, 0xa9, 0xbb, 0xc2, 0xb5, 
                0xc0, 0x9c, 0x99, 0x79, 0x44, 0x82, 0x38, 0x79, 
                0xbb, 0x99, 0xc3, 0x35, 0xc7, 0xa4, 0xdf, 0x4e, 
                0xd7,                
            },
        },
        { 
            [4]bool{
                false,  //transformPredict
                true,   //transformColor
                false,  //transformSubGreen
                false,  //transformColorIndexing
            },
            0,
            []byte{
                0x93, 0x0a, 0x64, 0x07, 0xfa, 0x1f, 0x10, 0x40, 
                0x24, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x33, 0x00, 0x00, 0x10, 0x40, 0x24, 0x0d, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x26, 
                0x40, 0xa6, 0x69, 0x07, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x80, 0x0b, 0x20, 0x88, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x2e, 0x10, 
                0xa6, 0x65, 0x13, 0xc5, 0x52, 0xd9, 0x24, 0x6c, 
                0xab, 0x48, 0x94, 0x4b, 0xab, 0x59, 0x2a, 0x13, 
                0x45, 0xdb, 0x32, 0xd7, 0x22, 0x41, 0x70, 0x79, 
                0x7c, 0x22, 0x33, 0x2b, 0x9b, 0x4b, 0xf0, 0x79, 
                0x99, 0x44, 0x76, 0xd6, 0xca, 0xcd, 0xca, 0x26, 
                0x32, 0xf9, 0x3c, 0xee, 0x9a, 0x49,
            },
        },
        {
            [4]bool{
                false,  //transformPredict
                true,   //transformColor
                false,  //transformSubGreen
                false,  //transformColorIndexing
            },
            8,
            []byte{
                0x53, 0xac, 0x40, 0x76, 0xa0, 0xff, 0x21, 0x42, 
                0x40, 0xb0, 0x6d, 0x06, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x94, 0x01, 0x00, 0x80, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x80, 0x03, 0x00, 0x00, 
                0x00, 0x00, 0x0c, 0x00, 0x00, 0x38, 0x0c, 0x06, 
                0x00, 0x00, 0x00, 0x1c, 0x00, 0x00, 0x00, 0x87, 
                0x03, 0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x06, 
                0x87, 0x03, 0x04, 0x10, 0x49, 0x03, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x80, 0x05, 0x10, 0x89, 
                0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 
                0x05, 0x10, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0xba, 0x80, 0x2d, 0x1b, 0xdb, 
                0x28, 0x29, 0x94, 0x93, 0xb7, 0x8b, 0x7f, 0x5c, 
                0xf3, 0x7c, 0xea, 0xed, 0x74, 0x2d, 0x30, 0x67, 
                0x66, 0x1e, 0x29, 0x89, 0x74, 0x70, 0x57, 0x33, 
                0x87, 0x6b, 0x8c, 0x49, 0xdf, 0x15, 0xae,            
            },
        },
        {
            [4]bool{
                true,   //transformPredict
                false,  //transformColor
                false,  //transformSubGreen
                false,  //transformColorIndexing
            },
            0,
            []byte{
                0x91, 0x12, 0x44, 0xf4, 0x3f, 0x60, 0x80, 0x6c, 
                0x9b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x65, 0x77, 0x00, 0x00, 0x04, 0x10, 0x49, 
                0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 
                0x05, 0x10, 0x49, 0x03, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x80, 0x09, 0x90, 0x69, 0x7e, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x60, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x27, 
                0x40, 0xc0, 0x0d, 0x80, 0x08, 0x00, 0x06, 0x43, 
                0x58, 0x69, 0x0f, 0x08, 0x43, 0x65, 0x7b, 0x78, 
                0x08, 0x9d, 0x9e, 0x0e, 0xb3, 0x6c, 0x16, 0x1b, 
                0x5d, 0xaf, 0xbb, 0xbd, 0x3e,
            },
        },
        {
            [4]bool{
                true,   //transformPredict
                false,  //transformColor
                false,  //transformSubGreen
                false,  //transformColorIndexing
            },
            8,
            []byte{
                0x51, 0x2c, 0x41, 0x44, 0xff, 0x43, 0xc4, 0x80, 
                0x60, 0x23, 0x31, 0x01, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0xca, 0xbd, 0x02, 0x00, 0x50, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x60, 
                0x05, 0x00, 0x80, 0x07, 0x00, 0x00, 0x00, 0x00, 
                0xe0, 0x00, 0x00, 0x60, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1e, 0x00, 
                0x00, 0x40, 0x00, 0x91, 0x1c, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x58, 0x00, 0x91, 0x34, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x98, 
                0x00, 0x99, 0xe6, 0x07, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x70, 0x02, 0x04, 0x72, 0x16, 
                0xa9, 0x90, 0x52, 0x7b, 0x43, 0x44, 0x98, 0xef, 
                0x66, 0xc5, 0xe6, 0x6b, 0x3c, 0x0c, 0xef, 0x0b, 
                0xc3, 0xf6, 0xd3, 0x0d, 0x3a, 0xd6, 0xba, 0x59, 
                0x1f,
            },
        },
        {
            [4]bool{
                true,   //transformPredict
                false,  //transformColor
                true,   //transformSubGreen
                false,  //transformColorIndexing
            },
            0,
            []byte{
                0x8d, 0x94, 0x20, 0xa2, 0xff, 0x01, 0x03, 0x64, 
                0xdb, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x28, 0xbb, 0x03, 0x00, 0x40, 0x80, 0x4c, 
                0xb3, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x27, 0x40, 0xa6, 0xd9, 0x0f, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x20, 0xd3, 
                0xfc, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x4e, 0x80, 0x80, 0x3b, 0x00, 0xa2, 0x06, 
                0xc0, 0xc1, 0x10, 0xd6, 0xd6, 0x1e, 0x10, 0x86, 
                0xda, 0xb6, 0x87, 0x87, 0x50, 0xa9, 0xa9, 0x30, 
                0x97, 0x9b, 0x8b, 0x8c, 0xbc, 0xd7, 0xf9, 0x5e, 
                0x1f,
            },
        },
        {
            [4]bool{
                true,   //transformPredict
                false,  //transformColor
                true,   //transformSubGreen
                false,  //transformColorIndexing
            },
            8,
            []byte{
                0x8d, 0x62, 0x09, 0x22, 0xfa, 0x1f, 0x22, 0x06, 
                0x04, 0x1b, 0x89, 0x09, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x50, 0xee, 0x15, 0x00, 0x80, 
                0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x2b, 0x00, 0x00, 0x3c, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x07, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf0, 0x00, 
                0x00, 0x00, 0x04, 0xc8, 0x34, 0xeb, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x70, 0x02, 0x64, 
                0x9a, 0x75, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x38, 0x01, 0x32, 0xcd, 0x0f, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0xe0, 0x04, 0x08, 
                0x70, 0x2e, 0xa8, 0x24, 0x55, 0xed, 0x0d, 0x88, 
                0x96, 0xf9, 0x6e, 0x56, 0x6b, 0xf3, 0x35, 0x1e, 
                0x1d, 0x7d, 0x5f, 0x38, 0xdc, 0x7e, 0xbc, 0x41, 
                0xc6, 0x5a, 0x36, 0xeb,
            },
        },
        {   // paletted image
            [4]bool{
                false,   //transformPredict
                false,  //transformColor
                false,   //transformSubGreen
                true,   //transformColorIndexing
            },
            4,
            []byte{
                0x67, 0x48, 0x06, 0xc8, 0xb6, 0xd9, 0x01, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x01, 0x00, 
                0x00, 0x8e, 0x00, 0x40, 0x00, 0x91, 0x34, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x58, 0x00, 
                0x91, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0xf8, 0x40, 0x80, 0xf1, 0x9b, 0x41, 0xc9, 
                0x39, 0x5d, 0x24, 0x08, 0x08, 0x00, 0x89, 0x99, 
                0x18, 0x04, 0x00, 0x40, 0x00, 0xc0, 0x00, 0x60, 
                0x00, 0x00, 0x30, 0x00, 0x30, 0x00, 0x01, 0x00, 
                0x00, 0x0c, 0x00, 0x0c, 0x08, 0x00, 0x00, 0x00, 
                0x03, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x80, 0xf8, 0xa7, 0xe7, 0x88, 
                0xfe, 0x07, 0xc6, 0x54, 0x1c, 0xf9, 0x7c, 0x71, 
                0x1b, 0x9d, 0xa9, 0xdc, 0x64, 0x40, 0x0d, 0x5d, 
                0xf9, 0xc9, 0x07, 0x5b, 
            },
        },
    }{
        b := &bytes.Buffer{}
        s := &bitWriter{Buffer: b}

        err := writeBitStreamData(s, img, tt.colorCacheBits, tt.transforms)
        if err != nil {
            t.Fatalf("test %v: writeBitStreamData returned error: %v", id, err)
        }

        result := b.Bytes()

        if !bytes.Equal(result, tt.expectedBytes) {
            t.Errorf("test %v: BitStream mismatch. Got %s, expected %s", id, result, tt.expectedBytes)
        }
    }
}

func TestWriteImageData(t *testing.T) {
    for id, tt := range []struct {
        inputPixels     []color.NRGBA
        width           int
        height          int
        isRecursive     bool
        colorCacheBits  int
        expectedBits    []byte
    }{
        {
            inputPixels: []color.NRGBA{
                {R: 100, G: 50, B: 150, A: 255},
                {R: 200, G: 100, B: 50, A: 255},
                {R: 100, G: 50, B: 150, A: 255}, // Same as the first pixel
            },
            width: 3,
            height: 1,
            isRecursive: false,
            colorCacheBits: 2,
            expectedBits: []byte{
                0x45, 0x00, 0x91, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x4f, 
                0x86, 0x7c, 0x19, 0xcb, 0xfe, 0x47,
            },
        },
        {
            inputPixels: []color.NRGBA{
                {R: 100, G: 50, B: 150, A: 255},
                {R: 200, G: 100, B: 50, A: 255},
                {R: 100, G: 50, B: 150, A: 255}, // Same as the first pixel
            },
            width: 3,
            height: 1,
            isRecursive: false,
            colorCacheBits: 0,
            expectedBits: []byte{
                0x2e, 0x43, 0x76, 0x32, 0xe4, 0xcb, 0x58, 0xf6, 
                0x3f, 0x38,
            },
        },
        {
            inputPixels: []color.NRGBA{
                {R: 100, G: 50, B: 150, A: 255},
                {R: 200, G: 100, B: 50, A: 255},
                {R: 100, G: 50, B: 150, A: 255}, // Same as the first pixel
            },
            width: 3,
            height: 1,
            isRecursive: true,
            colorCacheBits: 2,
            expectedBits: []byte{
                0x85, 0x00, 0x22, 0x01, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x9f, 
                0x0c, 0xf9, 0x32, 0x96, 0xfd, 0x8f, 0xd4,
            },
        },
        {
            inputPixels: []color.NRGBA{
                {R: 100, G: 50, B: 150, A: 255},
                {R: 200, G: 100, B: 50, A: 255},
                {R: 100, G: 50, B: 150, A: 255}, // Same as the first pixel
            },
            width: 3,
            height: 1,
            isRecursive: true,
            colorCacheBits: 0,
            expectedBits: []byte{
                0x5c, 0x86, 0xec, 0x64, 0xc8, 0x97, 0xb1, 0xec, 
                0x7f, 0x70,
            },
        },
    } {
        buffer := &bytes.Buffer{}
        writer := &bitWriter{
            Buffer:        buffer,
            BitBuffer:     0,
            BitBufferSize: 0,
        }

        writeImageData(writer, tt.inputPixels, tt.width, tt.height, tt.isRecursive, tt.colorCacheBits)

        if !bytes.Equal(buffer.Bytes(), tt.expectedBits) {
            t.Errorf("test %d: buffer mismatch\nexpected: %v got: %v", id, tt.expectedBits, buffer.Bytes())
            continue
        }
    }
}

func TestEncodeImageData(t *testing.T) {
    for id, tt := range []struct {
        inputPixels     []color.NRGBA
        width           int
        height          int
        colorCacheBits  int
        expectedEncoded []int
    }{
        {   //cached encoding
            inputPixels: []color.NRGBA{
                {R: 100, G: 50, B: 150, A: 255},
                {R: 200, G: 100, B: 50, A: 255},
                {R: 100, G: 50, B: 150, A: 255}, // Same as the first pixel
            },
            width: 3,
            height: 1,
            colorCacheBits: 2,
            expectedEncoded: []int{
                50, 100, 150, 255, // First pixel
                100, 200, 50, 255, // Second pixel
                256 + 24 + 3,       // Cached first pixel (hash index 0)
            },
        },
        {   //full RGBA encoding
            inputPixels: []color.NRGBA{
                {R: 100, G: 50, B: 150, A: 255},
                {R: 200, G: 100, B: 50, A: 255},
                {R: 100, G: 50, B: 150, A: 255}, // Same as the first pixel
            },
            width: 3,
            height: 1,
            colorCacheBits: 0,
            expectedEncoded: []int{
                50, 100, 150, 255,
                100, 200, 50, 255, 
                50, 100, 150, 255,  
            },
        },
    } {
        encoded := encodeImageData(tt.inputPixels, tt.width, tt.height, tt.colorCacheBits)

        if !reflect.DeepEqual(encoded, tt.expectedEncoded) {
            t.Errorf("test %d: encoded data mismatch\nexpected: %+v\n     got: %+v", id, tt.expectedEncoded, encoded)
            continue
        }
    }
}

func TestPrefixEncodeCode(t *testing.T) {
    tests := []struct {
        n                 int  // input value
        expectedCode      int  // expected prefix code
        expectedRemainder int  // expected remainder value
    }{
        // n <= 5: code should be max(0, n-1) and remainder 0.
        {-1, 0, 0}, // even negative numbers fall in this branch
        {0, 0, 0},
        {1, 0, 0},
        {2, 1, 0},
        {3, 2, 0},
        {4, 3, 0},
        {5, 4, 0},

        // n > 5: calculations using shifts.
        // For n = 6: n-1 = 5, loop runs once (5 >> 1 = 2) → shift=1, rem=2,
        // so returns (2 + 2*1, 6 - (2<<1) - 1) = (4, 1).
        {6, 4, 1},

        // For n = 7: n-1 = 6, loop: 6 >> 1 = 3 → shift=1, rem=3,
        // so returns (3 + 2*1, 7 - (3<<1) - 1) = (5, 0).
        {7, 5, 0},

        // For n = 8: n-1 = 7, loop: 7 >> 1 = 3 → shift=1, rem=3,
        // returns (3 + 2*1, 8 - (3<<1) - 1) = (5, 1).
        {8, 5, 1},

        // For n = 9: n-1 = 8, loop:
        // 8 >> 1 = 4, shift becomes 1; then 4 >> 1 = 2, shift becomes 2;
        // rem == 2 so returns (2 + 2*2, 9 - (2<<2) - 1) = (6, 0).
        {9, 6, 0},

        // For n = 10: returns (6, 1)
        {10, 6, 1},

        // For n = 11: returns (6, 2)
        {11, 6, 2},

        // For n = 12: returns (6, 3)
        {12, 6, 3},

        // For n = 13: n-1 = 12, loop: 12 >> 1 = 6 (shift=1),
        // then 6 >> 1 = 3 (shift=2), rem becomes 3 so returns (3+2*2, 13 - (3<<2) -1) = (7, 0).
        {13, 7, 0},

        // For n = 14: returns (7, 1)
        {14, 7, 1},

        // For n = 15: returns (7, 2)
        {15, 7, 2},

        // For n = 16: returns (7, 3)
        {16, 7, 3},
    }
    for idx, tt := range tests {
        code, remainder := prefixEncodeCode(tt.n)

        if code != tt.expectedCode {
            t.Errorf("Test %d: expected code %d, got %d", idx, tt.expectedCode, code)
            continue
        }
        if remainder != tt.expectedRemainder {
            t.Errorf("Test %d: expected remainder %d, got %d", idx, tt.expectedRemainder, remainder)
            continue
        }
    }
}

func TestPrefixEncodeBits(t *testing.T) {
    tests := []struct {
        prefix   int
        expected int
    }{
        // For prefix values less than 4, the function returns 0.
        {-10, 0},
        {-1, 0},
        {0, 0},
        {1, 0},
        {2, 0},
        {3, 0},
        // For prefix values 4 and above, the function computes (prefix-2) >> 1.
        // Example: For prefix = 4, (4-2) >> 1 = 2 >> 1 = 1.
        {4, 1},
        // For prefix = 5, (5-2) >> 1 = 3 >> 1 = 1.
        {5, 1},
        // For prefix = 6, (6-2) >> 1 = 4 >> 1 = 2.
        {6, 2},
        // For prefix = 7, (7-2) >> 1 = 5 >> 1 = 2.
        {7, 2},
        // For prefix = 8, (8-2) >> 1 = 6 >> 1 = 3.
        {8, 3},
        // For prefix = 9, (9-2) >> 1 = 7 >> 1 = 3.
        {9, 3},
        // For prefix = 10, (10-2) >> 1 = 8 >> 1 = 4.
        {10, 4},
        // Additional test cases
        {11, 4}, // (11-2)=9, 9 >> 1 = 4 (integer division)
        {12, 5}, // (12-2)=10, 10 >> 1 = 5
    }

    for idx, tt := range tests {
        result := prefixEncodeBits(tt.prefix)
        if result != tt.expected {
            t.Errorf("Test %d: expected %d got %d", idx, tt.expected, result)
        }
    }
}

func TestHash(t *testing.T) {
    tests := []struct {
        c        color.NRGBA
        shifts   int
        expected uint32
    }{
        {
            c:        color.NRGBA{R: 0, G: 0, B: 0, A: 0},
            shifts:   8,
            expected: 0,
        },
        {
            // Note: hash uses c.A as the most significant byte.
            // This test uses A=0, R=0, G=0, B=1 so that:
            //   x = 0<<24 | 0<<16 | 0<<8 | 1 = 1,
            //   then hash = (1 * 0x1e35a7bd) >> (32-8) = 0x1e35a7bd >> 24.
            // 0x1e35a7bd in hex is: 0x1e 0x35 0xa7 0xbd, so shifting right 24 bits yields 0x1e (30 in decimal).
            c:        color.NRGBA{R: 0, G: 0, B: 1, A: 0},
            shifts:   8,
            expected: 30,
        },
        {
            // Here x = 2 and so hash = (2*0x1e35a7bd) >> 24.
            // Since 0x1e35a7bd >> 24 is 30, doubling gives 60.
            c:        color.NRGBA{R: 0, G: 0, B: 2, A: 0},
            shifts:   8,
            expected: 60,
        },
        {
            // For c = {255,255,255,255} we have:
            //   x = 0xFF<<24 | 0xFF<<16 | 0xFF<<8 | 0xFF = 0xFFFFFFFF.
            // In 32-bit arithmetic, multiplying by 0x1e35a7bd gives:
            //   0xFFFFFFFF * 0x1e35a7bd ≡ -0x1e35a7bd (mod 2^32)
            // which equals 0x100000000 - 0x1e35a7bd = 0xE1CA5823 = 3788134467.
            c:        color.NRGBA{R: 255, G: 255, B: 255, A: 255},
            shifts:   32,
            expected: 3788134467,
        },
        {
            // Here x = 1<<24 = 0x01000000.
            // Multiplying by 0x1e35a7bd is equivalent to shifting the magic left 24 bits:
            //   (0x1e35a7bd << 24) mod 2^32.
            // Only the lower 8 bits of the magic survive in the final result,
            // so expected = (0x1e35a7bd & 0xFF) << 24 = 0xbd << 24 = 0xbd000000.
            c:        color.NRGBA{R: 0, G: 0, B: 0, A: 1},
            shifts:   32,
            expected: 0xbd000000,
        },
        {
            // With c = {R:0, G:0, B:1, A:0}, x = 1.
            // Then hash = (0x1e35a7bd) >> (32-16) = (0x1e35a7bd) >> 16.
            // Shifting 0x1e35a7bd right 16 bits yields 0x1e35, which is 7733 in decimal.
            c:        color.NRGBA{R: 0, G: 0, B: 1, A: 0},
            shifts:   16,
            expected: 7733,
        },
        {
            // case where shift is higher than maximum of 32 (should be set back to 32)
            c:        color.NRGBA{R: 255, G: 255, B: 255, A: 255},
            shifts:   33,
            expected: 3788134467,
        },
    }

    for id, tt := range tests {
        result := hash(tt.c, tt.shifts)
        if result != tt.expected {
            t.Errorf("test %v: expected hash as %v got %v", id, tt.expected, result)
        }
    }
}

func TestComputeHistograms(t *testing.T) {
    for id, tt := range []struct {
        pixels         []int
        colorCacheBits int
        expectedSizes  []int
        expectedCounts []map[int]int
    }{
        {
            pixels: []int{
                0xff, 0x01, 0x00, 0xff,
                0x00, 0xff, 0x00, 0xff,
                0x01, 0x01, 0xff, 0xff,
            },
            colorCacheBits: 0,
            expectedSizes:  []int{256 + 24, 256, 256, 256, 40},
            expectedCounts: []map[int]int{
                {0: 1, 1: 1, 255: 1}, // histos[0]
                {0: 0, 1: 2, 255: 1}, // histos[1]
                {0: 2, 1: 0, 255: 1}, // histos[2]
                {0: 0, 1: 0, 255: 3}, // histos[3]
                {},                   // histos[4] (unused in this case)
            },
        },
        {
            pixels: []int{
                0xff, 0x01, 0x00, 0xff, 
                0x00, 0xff, 0x00, 0xff,
                0x01, 0x01, 0xff, 0xff,
            },
            colorCacheBits: 4,
            expectedSizes:  []int{256 + 24 + (1 << 4), 256, 256, 256, 40},
            expectedCounts: []map[int]int{
                {0: 1, 1: 1, 255: 1}, // histos[0]
                {0: 0, 1: 2, 255: 1}, // histos[1]
                {0: 2, 1: 0, 255: 1}, // histos[2]
                {0: 0, 1: 0, 255: 3}, // histos[3]
                {},                   // histos[4] (unused in this case)
            },
        },
        {
            pixels: []int{
                0x104, 0x01, 0x02, 0x03, // over 256
                0xff, 0x01, 0x00, 0xff, 
                0x00, 0xff, 0x00, 0xff,
                0x01, 0x01, 0xff, 0xff,
            },
            colorCacheBits: 4,
            expectedSizes:  []int{256 + 24 + (1 << 4), 256, 256, 256, 40},
            expectedCounts: []map[int]int{
                {0: 1, 1: 1, 255: 1}, // histos[0]
                {0: 0, 1: 2, 255: 1}, // histos[1]
                {0: 2, 1: 0, 255: 1}, // histos[2]
                {0: 0, 1: 0, 255: 3}, // histos[3]
                {2: 1},               // histos[4] (unused in this case)
            },
        },
    }{
        histos := computeHistograms(tt.pixels, tt.colorCacheBits)

        for i, histo := range histos {
            if len(histo) != tt.expectedSizes[i] {
                t.Errorf("test %d: histos[%d] size mismatch\nexpected: %d\ngot: %d", id, i, tt.expectedSizes[i], len(histo))
                continue
            }
        }

        for histoIdx, expectedCounts := range tt.expectedCounts {
            for value, expectedCount := range expectedCounts {
                if histos[histoIdx][value] != expectedCount {
                    t.Errorf("test %d: histos[%d][%d] count mismatch\nexpected: %d\ngot: %d", id, histoIdx, value, expectedCount, histos[histoIdx][value])
                    continue
                }
            }
        }
    }
}

func TestFlatten(t *testing.T) {
    for id, tt := range []struct {
        width       int
        height      int
        brightness  float64
        hasAlpha    bool
        expectError bool
        expectedErrorMsg string
    }{
        // Valid NRGBA image with alpha
        {
            width:       16,
            height:      16,
            brightness:  64,
            hasAlpha:    true,
            expectError: false,
            expectedErrorMsg: "",
        },
        // Valid NRGBA image without alpha
        {
            width:       16,
            height:      16,
            brightness:  64,
            hasAlpha:    false,
            expectError: false,
            expectedErrorMsg: "",
        },
        // Unsupported image format
        {
            width:       16,
            height:      16,
            brightness:  64,
            hasAlpha:    true,
            expectError: true, // Will convert to an unsupported format
            expectedErrorMsg: "unsupported image format",
        },
    }{
        img := generateTestImageNRGBA(tt.width, tt.height, tt.brightness, tt.hasAlpha)

        var testImage image.Image = img
        if tt.expectError {
            testImage = image.NewGray(img.Bounds())
        }

        pixels, err := flatten(testImage)

        if tt.expectError {
            if err == nil {
                t.Errorf("test %d: expected error but got nil", id)
                continue
            }
            
            if err.Error() != tt.expectedErrorMsg {
                t.Errorf("test %d: expected error %v got %v", id, tt.expectedErrorMsg, err)
                continue
            }

            continue
        }

        if err != nil {
            t.Errorf("test %d: unexpected error: %v", id, err)
            continue
        }

        for y := 0; y < tt.height; y++ {
            for x := 0; x < tt.width; x++ {
                index := y*tt.width + x
                expected := img.At(x, y).(color.NRGBA)
                actual := pixels[index]

                if expected != actual {
                    t.Errorf("test %d: pixel mismatch at (%d, %d): expected %+v, got %+v", id, x, y, expected, actual)
                    continue
                }
            }
        }
    }
}
```

---
**End of Code Context Bundle**
