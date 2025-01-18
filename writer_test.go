package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "bytes"
    "reflect"
    "encoding/hex"
    "crypto/sha256"
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

func TestWriteWebPHeader(t *testing.T) {
    for id, tt := range []struct {
        inputData       []byte 
        expectedHeader  []byte
    }{
        // Test case with an empty 'b' buffer
        {
            inputData:     []byte{},
            expectedHeader: []byte{
                'R', 'I', 'F', 'F',         // RIFF
                0x0C, 0x00, 0x00, 0x00,     // 12 in little-endian (12 + 0)
                'W', 'E', 'B', 'P',         // WEBP
                'V', 'P', '8', 'L',         // VP8L
                0x00, 0x00, 0x00, 0x00,     // 0 in little-endian (size of 'b' buffer)
            },
        },
        // Test case with non-empty 'b' buffer
        {
            inputData:     []byte{1, 2, 3, 4, 5},
            expectedHeader: []byte{
                'R', 'I', 'F', 'F',         // RIFF
                0x11, 0x00, 0x00, 0x00,     // 12 in little-endian (12 + 5)
                'W', 'E', 'B', 'P',         // WEBP
                'V', 'P', '8', 'L',         // VP8L
                0x05, 0x00, 0x00, 0x00,     // 0 in little-endian (size of 'b' buffer)
            },
        },
    }{
        w := &bytes.Buffer{}
        b := bytes.NewBuffer(tt.inputData)

        writeWebPHeader(w, b)

        if !bytes.Equal(w.Bytes(), tt.expectedHeader) {
            t.Errorf("test %d: header mismatch expected: %v got: %v", id, tt.expectedHeader, w.Bytes())
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

func TestWriteEncodeErrors(t *testing.T) {
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
    }{
        b := &bytes.Buffer{}

        err := Encode(b, tt.img)
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
        img             image.Image
        expectedBytes   []byte
    }{
        {
            generateTestImageNRGBA(8, 8, 64, true),
            []byte {
                0x52, 0x49, 0x46, 0x46, 0xda, 0x00, 0x00, 0x00, 
                0x57, 0x45, 0x42, 0x50, 0x56, 0x50, 0x38, 0x4c, 
                0xce, 0x00, 0x00, 0x00, 0x2f, 0x07, 0xc0, 0x01, 
                0x10, 0x8d, 0x52, 0x46, 0xf4, 0x3f, 0x24, 0x0c, 
                0x08, 0x36, 0x12, 0x93, 0x03, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x98, 
                0xc9, 0x6d, 0x50, 0xff, 0x04, 0xc8, 0x34, 0xeb, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x70, 
                0x02, 0x64, 0x9a, 0xef, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x90, 0x00, 0x22, 0x31, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x60, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x60, 
                0xf1, 0xc1, 0x48, 0x93, 0xd0, 0xa5, 0xcd, 0x37, 
                0x76, 0xe6, 0x79, 0x27, 0x57, 0xcf, 0x66, 0x38, 
                0xcf, 0x6f, 0xfb, 0xdd, 0xf6, 0x3d, 0x88, 0x82, 
                0x40, 0x11, 0x18, 0xdd, 0xf8, 0xa3, 0xc3, 0x55, 
                0x80, 0xc4, 0x39, 0x12, 0xd8, 0xc2, 0xe1, 0xf6, 
                0x01, 0x00,
            },
        },
    }{
        b := &bytes.Buffer{}
        Encode(b, tt.img)

        result := b.Bytes()
        if !bytes.Equal(result, tt.expectedBytes) {
            t.Errorf("test %v: BitStream mismatch. Got %s, expected %s", id, result, tt.expectedBytes)
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
                0xc5, 0xfc, 0x0b, 0x20, 0x92, 0x06, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x0b, 0x20, 0x92, 
                0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x0b, 0x20, 0x88, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x80, 0x00, 0xa6, 0x65, 0x83, 0xe5, 
                0x6a, 0x6d, 0x02, 0xdb, 0x5a, 0x82, 0xf5, 0xca, 
                0xb2, 0x81, 0xb9, 0x5a, 0x83, 0xa5, 0x6d, 0x99, 
                0x60, 0xbd, 0x5a, 0x02, 0xc0, 0x72, 0x3c, 0x28, 
                0xab, 0x9a, 0x05, 0x3c, 0x57, 0x82, 0xba, 0xe2, 
                0x78, 0xc0, 0x56, 0x35, 0x28, 0x79, 0x8e, 0x05, 
                0x75, 0x55, 
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
                0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x00, 
                0x10, 0x00, 0xc0, 0x01, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x83, 0x39, 0x00, 0x00, 0x00, 
                0xc0, 0x01, 0x00, 0x00, 0x60, 0x06, 0x00, 0x00, 
                0x00, 0x30, 0x00, 0x00, 0xc0, 0xe1, 0x06, 0x00, 
                0x00, 0x00, 0x70, 0x80, 0x00, 0x22, 0x69, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xb0, 0x00, 
                0x22, 0x39, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0xb0, 0x00, 0x22, 0x01, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x60, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x40, 0x60, 0x97, 0x1d, 
                0xbb, 0x51, 0x52, 0x28, 0x07, 0x2f, 0x0f, 0xbf, 
                0xf3, 0x11, 0xfd, 0x39, 0x2f, 0x83, 0xdf, 0x07, 
                0xcc, 0x99, 0x99, 0x47, 0x24, 0x88, 0x37, 0x53, 
                0x17, 0x5f, 0x35, 0xec, 0xe5, 0xdc, 0xa9, 0x8d, 
                0xaf,                 
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
                0x93, 0x0a, 0x64, 0x07, 0xfa, 0x1f, 0x98, 0x7f, 
                0x01, 0x44, 0xd2, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x60, 0x02, 0x64, 0x9a, 0x76, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xb8, 0x00, 
                0x82, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x08, 0xc2, 0xb4, 0x6c, 0x62, 0x9c, 0x66, 
                0x93, 0xb0, 0xad, 0x91, 0x98, 0x27, 0xcb, 0x26, 
                0xcc, 0x69, 0x26, 0x46, 0xdb, 0x32, 0x89, 0x79, 
                0x1a, 0x09, 0x82, 0xcb, 0xe3, 0x13, 0xcd, 0x56, 
                0x9b, 0x4b, 0xf0, 0x79, 0x4d, 0xa2, 0xdd, 0xe2, 
                0xf1, 0x09, 0x6e, 0xab, 0x4d, 0x34, 0xf9, 0x3c, 
                0x2e, 0xd1, 0x6e, 0x35,
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
                0x00, 0x00, 0x00, 0x14, 0x00, 0x00, 0x20, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0xe0, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x03, 0x00, 0x00, 0x06, 0x87, 0x03, 
                0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0xc0, 0x60, 
                0x00, 0x00, 0x00, 0xc0, 0x01, 0x00, 0x80, 0xc3, 
                0xe0, 0x00, 0x01, 0x44, 0xd2, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x60, 0x01, 0x44, 0xa2, 
                0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x60, 
                0x01, 0x44, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x80, 0xc0, 0x96, 0x8d, 0x6d, 0x94, 
                0x14, 0xca, 0x9b, 0xa9, 0xc3, 0xef, 0x1a, 0x76, 
                0x7f, 0xce, 0xd4, 0xe6, 0xf7, 0x01, 0x73, 0x66, 
                0xe6, 0x91, 0x92, 0x48, 0x3b, 0x37, 0x1f, 0xb1, 
                0xf2, 0xe2, 0x2b, 0x78, 0x37, 0x9d, 0x58,             
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
                0x91, 0x8c, 0xe8, 0x7f, 0x40, 0x00, 0x91, 0x34, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 
                0x00, 0x00, 0x40, 0x00, 0x91, 0x34, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x58, 0x00, 0x91, 
                0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x58, 0x00, 0x91, 0x18, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x20, 0x90, 0x01, 0x10, 0x03, 
                0xc0, 0xc8, 0x36, 0x3e, 0xc7, 0x33, 0x06, 0x10, 
                0xe6, 0x80, 0x30, 0x00, 0x63, 0x79, 0xc6, 0x38, 
                0x16, 0x0c, 0x80, 0x3c, 0x00, 0x8c, 0xe6, 0xf1, 
                0x77, 0xf9, 0x8d, 0x01, 0x08, 0xfb, 0x40, 0xf1, 
                0x00, 0xc6, 0xf4, 0x1b, 0xe3,
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
                0x51, 0xcc, 0x88, 0xfe, 0x87, 0x88, 0x01, 0xc1, 
                0x46, 0x62, 0x7a, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0xa0, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x3c, 0x00, 0x08, 0x30, 
                0x00, 0x70, 0x1e, 0x78, 0x80, 0x01, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x80, 0x01, 0x00, 0xc0, 0x00, 
                0x40, 0x81, 0x07, 0x00, 0x00, 0x06, 0xf0, 0x00, 
                0x60, 0x00, 0x08, 0x20, 0x92, 0x06, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x0b, 0x20, 0x92, 
                0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x0b, 0x20, 0x12, 0x03, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x84, 0x48, 0x23, 0xa0, 0x30, 
                0xa0, 0x95, 0xb5, 0xc5, 0xce, 0xc5, 0xe6, 0xc0, 
                0x69, 0x07, 0x85, 0x97, 0xe4, 0x1e, 0xc9, 0x7f, 
                0x64, 0x20, 0xf1, 0xc0, 0x94, 0x96, 0x8f, 0xef, 
                0xca, 0xde, 0x81, 0x35, 0x7d, 0xa8, 0xd8, 0x24, 
                0xdf, 0x95, 0x3c, 
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
                0x8d, 0x64, 0x44, 0xff, 0x03, 0x02, 0x88, 0xa4, 
                0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 
                0x00, 0x00, 0x00, 0x04, 0xc8, 0x34, 0xeb, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x70, 0x02, 
                0x64, 0x9a, 0x75, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0xb8, 0x00, 0x22, 0x31, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x60, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x20, 0x07, 
                0x40, 0x5c, 0x00, 0x1c, 0xd9, 0x8d, 0x5f, 0xf7, 
                0x7f, 0x8d, 0x01, 0x84, 0xf3, 0x02, 0xfd, 0x2e, 
                0x80, 0xb1, 0x7c, 0x8d, 0x71, 0xff, 0x82, 0x03, 
                0x20, 0x0b, 0x00, 0x47, 0xf3, 0xf1, 0x76, 0xfe, 
                0x1a, 0x03, 0x10, 0x6e, 0x03, 0xb9, 0x05, 0x30, 
                0xa6, 0xd7, 0x18, 0xe7, 
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
                0x8d, 0x62, 0x46, 0xf4, 0x3f, 0x44, 0x0c, 0x08, 
                0x36, 0x12, 0xd3, 0x03, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 
                0x0f, 0x00, 0x00, 0x0f, 0x00, 0x00, 0x00, 0x02, 
                0x00, 0x18, 0xbc, 0x03, 0x06, 0x30, 0x00, 0x80, 
                0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x0a, 0x00, 0x78, 0x00, 0x30, 0x80, 0x07, 
                0x00, 0x18, 0x80, 0x00, 0x99, 0x66, 0x1d, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4e, 0x80, 
                0x4c, 0xd3, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x17, 0x40, 0x24, 0x06, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0x29, 0x3d, 
                0x12, 0x54, 0x1c, 0xb0, 0x54, 0xeb, 0x2d, 0xed, 
                0xb9, 0xfe, 0x9d, 0x03, 0x67, 0x79, 0x50, 0xbf, 
                0x97, 0xe4, 0xf9, 0x93, 0x79, 0xe4, 0x40, 0x30, 
                0xa0, 0xc5, 0xf2, 0x63, 0xbb, 0x7c, 0x73, 0xe0, 
                0x6e, 0x1b, 0xca, 0x6d, 0x92, 0xef, 0x4a,
            },
        },
        {   // paletted image
            [4]bool{
                true,   //transformPredict
                false,  //transformColor
                true,   //transformSubGreen
                true,   //transformColorIndexing
            },
            8,
            []byte{
                0x67, 0x88, 0x06, 0xc8, 0xb6, 0xd9, 0x01, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x01, 0x00, 
                0x00, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x40, 0x00, 0x91, 0x34, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x58, 0x00, 0x91, 0x34, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf8, 
                0x40, 0x80, 0xf1, 0x9b, 0x41, 0xc9, 0x39, 0x5d, 
                0x1a, 0xc5, 0x8c, 0xe8, 0x7f, 0x88, 0x1c, 0x10, 
                0x0c, 0x24, 0x36, 0xda, 0x56, 0x6b, 0x55, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x40, 0xb5, 0x3a, 0xad, 0x02, 
                0x00, 0x00, 0x1c, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x50, 0x80, 0x00, 0x00, 0x28, 0x14, 0xe0, 0x01, 
                0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0xa0, 
                0x00, 0x00, 0x80, 0x07, 0x28, 0x00, 0x00, 0xa0, 
                0x00, 0x00, 0xa0, 0x00, 0x00, 0x00, 0x04, 0x04, 
                0x80, 0x44, 0x5a, 0x39, 0x01, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x52, 0x4a, 0x15, 0x02, 0x02, 0x40, 0x22, 
                0xad, 0x9c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
                0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x29, 
                0xa5, 0x2a, 0x22, 0x8a, 0x3b, 0x5d, 0x11, 0xeb, 
                0x6e, 0x75, 0x03, 0xa8, 0x66, 0x28, 0x77, 0x97, 
                0xf7, 0x9a, 0x2a, 0x34, 0x99, 0xfd, 0x13, 0x65, 
                0x81, 0x0b, 0xbe, 0xd4, 0x66, 0xea, 0xf4, 0x7f, 
                0xe7, 0xae, 0x31, 0x33, 0xff, 0x22, 0x78, 0x47, 
                0x91, 0xde, 0x0c, 0x7d, 0x22, 0x3b, 0xa8, 0x9e, 
                0xc2, 0xbd, 0xc3,
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

        writeImageData(writer, tt.inputPixels, tt.isRecursive, tt.colorCacheBits)

        if !bytes.Equal(buffer.Bytes(), tt.expectedBits) {
            t.Errorf("test %d: buffer mismatch\nexpected: %v got: %v", id, tt.expectedBits, buffer.Bytes())
            continue
        }
    }
}

func TestEncodeImageData(t *testing.T) {
    for id, tt := range []struct {
        inputPixels     []color.NRGBA
        colorCacheBits  int
        expectedEncoded []uint16
    }{
        {   //cached encoding
            inputPixels: []color.NRGBA{
                {R: 100, G: 50, B: 150, A: 255},
                {R: 200, G: 100, B: 50, A: 255},
                {R: 100, G: 50, B: 150, A: 255}, // Same as the first pixel
            },
            colorCacheBits: 2,
            expectedEncoded: []uint16{
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
            colorCacheBits: 0,
            expectedEncoded: []uint16{
                50, 100, 150, 255,
                100, 200, 50, 255, 
                50, 100, 150, 255,  
            },
        },
    } {
        encoded := encodeImageData(tt.inputPixels, tt.colorCacheBits)

        if !reflect.DeepEqual(encoded, tt.expectedEncoded) {
            t.Errorf("test %d: encoded data mismatch\nexpected: %+v\n     got: %+v", id, tt.expectedEncoded, encoded)
            continue
        }
    }
}

func TestComputeHistograms(t *testing.T) {
    for id, tt := range []struct {
        pixels         []uint16
        colorCacheBits int
        expectedSizes  []int
        expectedCounts []map[int]int
    }{
        {
            pixels: []uint16{
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
            pixels: []uint16{
                0x101,                  // larger than 256
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

func TestApplyPredicttransform(t *testing.T) {
    for id, tt := range []struct {
        width           int
        height          int
        expectedHash    string
        expectedBlocks  []color.NRGBA
        expectedBit     int
    }{
        {   // default case
            32,
            32,
            "3c3a5319fe90b766abf54876f70f21f5322f2b1bad5884800529f082de30cfe1",
            []color.NRGBA{
                {0, 1, 0, 255},
                {0, 1, 0, 255},
                {0, 1, 0, 255},
                {0, 1, 0, 255},
            }, 
            4,
        },
        {   // not power of 2 image res
            33,
            33,
            "3812a5cd02c500ea176d2710521990796e149cf7b25ac0c9bd74b3a665d0637c",
            []color.NRGBA{
                {0, 1, 0, 255},
                {0, 1, 0, 255},
                {0, 1, 0, 255},
                {0, 1, 0, 255},
                {0, 1, 0, 255},
                {0, 1, 0, 255},
                {0, 1, 0, 255},
                {0, 1, 0, 255},
                {0, 1, 0, 255},
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

        tileBit, blocks := applyPredicttransform(pixels, tt.width, tt.height)

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

func TestApplyColortransform(t *testing.T) {
    for id, tt := range []struct {
        width          int
        height         int
        expectedHash   string
        expectedBlocks []color.NRGBA
        expectedBit    int
    }{
        {   // default case
            32,
            32,
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

        tileBit, blocks := applyColortransform(pixels, tt.width, tt.height)

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

func TestApplySubtractGreentransform(t *testing.T) {
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

        applySubtractGreentransform(pixels)

        if !reflect.DeepEqual(pixels, tt.expectedPixels) {
            t.Errorf("test %d: pixel mismatch\nexpected: %+v\n     got: %+v", id, tt.expectedPixels, pixels)
            continue
        }
    }
}

func TestApplyPalettetransformWithManualPixels(t *testing.T) {
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

    _, err := applyPalettetransform(pixels)

    msg := "palette exceeds 256 colors"
    if err == nil || err.Error() != msg {
        t.Errorf("test: expected error %v got %v", msg, err)
    }

    for id, tt := range []struct {
        pixels     []color.NRGBA
        expectedPalette []color.NRGBA
        expectedPixels  []color.NRGBA
    }{
        {
            pixels: []color.NRGBA{
                {R: 255, G: 0, B: 0, A: 255}, 
                {R: 0, G: 255, B: 0, A: 255}, 
                {R: 0, G: 0, B: 255, A: 255}, 
                {R: 255, G: 255, B: 0, A: 255},
                {R: 255, G: 0, B: 0, A: 255},  // repeat color 1
            },
            expectedPalette: []color.NRGBA{
                {R: 255, G: 0, B: 0, A: 255},
                {R: 1, G: 255, B: 0, A: 0},
                {R: 0, G: 1, B: 255, A: 0},
                {R: 255, G: 255, B: 1, A: 0},
            },
            expectedPixels: []color.NRGBA{
                {R: 0, G: 0, B:0, A: 255}, 
                {R: 0, G: 1, B:0, A: 255},  
                {R: 0, G: 2, B:0, A: 255}, 
                {R: 0, G: 3, B:0, A: 255}, 
                {R: 0, G: 0, B:0, A: 255}, 
            },
        },
    } {
        // Copy inputPixels to avoid modifying the test case
        pixels := make([]color.NRGBA, len(tt.pixels))
        copy(pixels, tt.pixels)

        pal, err := applyPalettetransform(pixels)

        if err != nil {
            t.Errorf("test %d: unexpected error %v", id, err)
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