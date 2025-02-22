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
    "image/draw"
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
    buf := new(bytes.Buffer)

    err := Encode(buf, img, nil)
    if err != nil {
         t.Errorf("Encode: expected err as nil got %v", err)
         return
    }

    img1, ok := img.(*image.NRGBA)
    if !ok {
        t.Errorf("test: unsupported image format for img1")
        return
    }

    for id, tt := range []struct {
        input       []byte 
        expected    *image.NRGBA
        expectedErr string
    }{
        {
            buf.Bytes(),
            img1,
            "",
        },
        {
            []byte("invalid WebP data"),
            nil,
            "riff: missing RIFF chunk header",
        },
    }{

        buf := bytes.NewBuffer(tt.input)
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

        img2 := image.NewNRGBA(result.Bounds())
        draw.Draw(img2, result.Bounds(), result, result.Bounds().Min, draw.Src)

        if !tt.expected.Rect.Eq(img2.Rect) || tt.expected.Stride != img2.Stride {
            t.Errorf("test %d: expected image dimensions as %v got %v", id, tt.expected.Rect, img2.Rect)
            continue
        }
        
        if !bytes.Equal(tt.expected.Pix, img2.Pix) {
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
        img             image.Image
        expectedBytes   []byte
    }{
        {
            generateTestImageNRGBA(8, 8, 64, true),
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
    }{
        b := &bytes.Buffer{}
        Encode(b, tt.img, nil)

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