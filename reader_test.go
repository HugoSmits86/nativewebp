package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "bytes"
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
            "",
        },
        {
            nil,    // if nil is used create a non-webp buffer
            false,
            "webp: invalid format",
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

func TestDecodeVP8XHeader(t *testing.T) {
    img := generateTestImageNRGBA(8, 8, 64, true)
    buf := new(bytes.Buffer)
    
    err := Encode(buf, img, &Options{UseExtendedFormat: true})
    if err != nil {
        t.Errorf("expected err as nil got %v", err)
        return
    }

    data := buf.Bytes()

    for id, tt := range []struct {
        ICCProfileFlag      bool
        AnimationFlag       bool
        expectedErr         string
    }{
        {
            true, 
            false,
            "webp: invalid format",
        },
        {
            false, 
            true,
            "webp: invalid format",
        },
        {
            true, 
            true,
            "webp: invalid format",
        },
        {
            false, 
            false,
            "",
        },
    }{
        modifiedData := make([]byte, len(data))
        copy(modifiedData, data)

        if tt.ICCProfileFlag {
            modifiedData[20] |= (1 << 1) // Set color profile flag
        }

        if tt.AnimationFlag {
            modifiedData[20] |= (1 << 5) // Set animation flag
        }

        _, err = Decode(bytes.NewReader(modifiedData))
        if err == nil && tt.expectedErr != "" {
            t.Errorf("test %d: expected err as %v got nil", id, tt.expectedErr)
            continue
        }

        if err != nil && tt.expectedErr == "" {
            t.Errorf("test %d: expected err as nil got %v", id, err)
            continue
        }

        if err != nil && tt.expectedErr != err.Error() {
            t.Errorf("test %d: expected err as %v got %v", id, tt.expectedErr, err)
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