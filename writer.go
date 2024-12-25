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
    //log
    "errors"
)

type Transform int

const (
    TransformPredict    = Transform(0)
    TransformSubGreen   = Transform(2)    
)

func Encode(w io.Writer, img image.Image) error {
    if img == nil {
        return errors.New("image is nil")
    }

    if img.Bounds().Dx() < 1 || img.Bounds().Dy() < 1 {
        return errors.New("invalid image size")
    }

    rgba := image.NewNRGBA(img.Bounds())
    draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)

    b := &bytes.Buffer{}
    s := &BitWriter{Buffer: b}

    writeBitStreamHeader(s, rgba.Bounds(), !rgba.Opaque())

    var transforms [4]bool
    transforms[TransformPredict] = true
    transforms[TransformSubGreen] = true

    WriteBitStreamData(s, rgba, 4, transforms)

    s.AlignByte()

    if b.Len() % 2 != 0 {
        b.Write([]byte{0x00})
    }

    writeWebPHeader(w, b)

    data := b.Bytes()
    w.Write(data)

    return nil
}

func writeWebPHeader(w io.Writer, b *bytes.Buffer) {
    w.Write([]byte("RIFF"))

    tmp := make([]byte, 4)
    binary.LittleEndian.PutUint32(tmp, uint32(12 + b.Len()))
    w.Write(tmp)

    w.Write([]byte("WEBP"))
    w.Write([]byte("VP8L"))

    tmp = make([]byte, 4)
    binary.LittleEndian.PutUint32(tmp, uint32(b.Len()))
    w.Write(tmp)
}

func writeBitStreamHeader(w *BitWriter, bounds image.Rectangle, hasAlpha bool) {
    w.WriteBits(0x2f, 8)

    w.WriteBits(uint64(bounds.Dx() - 1), 14)
    w.WriteBits(uint64(bounds.Dy() - 1), 14)

    if hasAlpha {
        w.WriteBits(1, 1)
    } else {
        w.WriteBits(0, 1)
    }

    w.WriteBits(0, 3)
}

func WriteBitStreamData(w *BitWriter, img image.Image, colorCacheBits int, transforms [4]bool) error {
    pixels, err := flatten(img)
    if err != nil {
        return err
    }

    if transforms[TransformSubGreen] {
        w.WriteBits(1, 1)
        w.WriteBits(2, 2)

        ApplySubtractGreenTransform(pixels)
    }

    if transforms[TransformPredict] {
        w.WriteBits(1, 1)
        w.WriteBits(0, 2)

        bits, blocks := ApplyPredictTransform(pixels, img.Bounds().Dx(), img.Bounds().Dy())

        w.WriteBits(uint64(bits - 2), 3);
        WriteImageData(w, blocks, false, colorCacheBits)
    }

    w.WriteBits(0, 1) // end of transform
    WriteImageData(w, pixels, true, colorCacheBits)

    return nil
}

func WriteImageData(w *BitWriter, pixels []color.NRGBA, isRecursive bool, colorCacheBits int) {
    if colorCacheBits > 0 {
        w.WriteBits(1, 1)
        w.WriteBits(uint64(colorCacheBits), 4) 
    } else {
        w.WriteBits(0, 1)
    }

    if isRecursive {
        w.WriteBits(0, 1)
    }

    encoded := encodeImageData(pixels, colorCacheBits)
    histos := computeHistograms(encoded, colorCacheBits)

    var codes [][]HuffmanCode
    for i := 0; i < 5; i++ {
        c := buildHuffmanCodes(histos[i], 16)
        codes = append(codes, c)

        writeHuffmanCodes(w, c)
    }

    for i := 0; i < len(encoded); i ++ {
        w.WriteCode(codes[0][encoded[i + 0]])
        if encoded[i + 0] < 256 {
            w.WriteCode(codes[1][encoded[i + 1]])
            w.WriteCode(codes[2][encoded[i + 2]])
            w.WriteCode(codes[3][encoded[i + 3]])
            i += 3
        }
    }
}

func encodeImageData(pixels []color.NRGBA, colorCacheBits int) []uint16 {
    cache := make([]color.NRGBA, 1 << colorCacheBits)
    encoded := make([]uint16, len(pixels) * 4)

    cnt := 0
    for _, p := range pixels {
        hash := 0
        if colorCacheBits > 0 {
            //hash formula including magic number 0x1e35a7bd comes directly from WebP specs!
            pack := uint32(p.A) << 24 | uint32(p.R) << 16 | uint32(p.G) << 8 | uint32(p.B)
            hash = int((pack * 0x1e35a7bd) >> (32 - colorCacheBits))

            if cache[hash] == p {
                encoded[cnt] = uint16(hash + 256 + 24)
                cnt++
                continue
            }

            cache[hash] = p
        }

        encoded[cnt + 0] = uint16(p.G)
        encoded[cnt + 1] = uint16(p.R)
        encoded[cnt + 2] = uint16(p.B)
        encoded[cnt + 3] = uint16(p.A)
        cnt += 4
    }

    return encoded[:cnt]
}

func computeHistograms(pixels []uint16, colorCacheBits int) [][]int {
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

func ApplySubtractGreenTransform(pixels []color.NRGBA) {
    for i, _ := range pixels {
        pixels[i].R = pixels[i].R - pixels[i].G
        pixels[i].B = pixels[i].B - pixels[i].G
    }
}

func ApplyPredictTransform(pixels []color.NRGBA, width, height int) (int, []color.NRGBA) {
    tileBits := 4
    tileSize := 1 << tileBits
    bw := (width + tileSize - 1) / tileSize
    bh := (height + tileSize - 1) / tileSize

    blocks := make([]color.NRGBA, bw * bh)
    deltas := make([]color.NRGBA, width * height)
    
    best := 1
    for y := 0; y < bh; y++ {
        for x := 0; x < bw; x++ {
            mx := min((x + 1) << tileBits, width)
            my := min((y + 1) << tileBits, height)

            for tx := x << tileBits; tx < mx; tx++ {
                for ty := y << tileBits; ty < my; ty++ {
                    var d color.NRGBA
                    if tx == 0 && ty == 0 {
                        d = color.NRGBA{0, 0, 0, 255}
                    } else if tx == 0 {
                        d = pixels[(ty - 1) * width + tx]
                    } else if ty == 0 {
                        d = pixels[ty * width + (tx - 1)]
                    } else { 
                        // left prediction 
                        d = pixels[ty * width + (tx - 1)]
                    }
                    
                    off := ty * width + tx

                    deltas[off] = color.NRGBA{
                        R: uint8(pixels[off].R - d.R),
                        G: uint8(pixels[off].G - d.G),
                        B: uint8(pixels[off].B - d.B),
                        A: uint8(pixels[off].A - d.A),
                    }
                }
            }

            blocks[y * bw + x] = color.NRGBA{0, byte(best), 0, 255}
        }
    }
    
    copy(pixels, deltas)
    
    return tileBits, blocks
}