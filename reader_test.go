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