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
            3,
            3,
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