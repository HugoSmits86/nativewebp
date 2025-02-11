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
    //"log"
    "errors"
)

// Encode writes the provided image.Image to the specified io.Writer in WebP VP8L format.
//
// This function supports VP8L (lossless WebP) encoding and can handle color-indexed images
// when img is provided as image.Paletted.
//
// Parameters:
//   w   - The destination writer where the encoded WebP image will be written.
//   img - The input image to be encoded.
//
// Returns:
//   An error if encoding fails or writing to the io.Writer encounters an issue.
func Encode(w io.Writer, img image.Image) error {
    if img == nil {
        return errors.New("image is nil")
    }

    if img.Bounds().Dx() < 1 || img.Bounds().Dy() < 1 {
        return errors.New("invalid image size")
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
        return err
    }
    
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

    if transforms[transformColorIndexing] {
        w.writeBits(1, 1)
        w.writeBits(3, 2)
       
        pal, err := applyPaletteTransform(pixels)
        if err != nil {
            return err
        }
       
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

        bits, bw, bh, blocks := applyColorTransform(pixels, img.Bounds().Dx(), img.Bounds().Dy())

        w.writeBits(uint64(bits - 2), 3);
        writeImageData(w, blocks, bw, bh, false, colorCacheBits)
    }

    if transforms[transformPredict] {
        w.writeBits(1, 1)
        w.writeBits(0, 2)

        bits, bw, bh, blocks := applyPredictTransform(pixels, img.Bounds().Dx(), img.Bounds().Dy())

        w.writeBits(uint64(bits - 2), 3);
        writeImageData(w, blocks, bw, bh, false, colorCacheBits)
    }

    w.writeBits(0, 1) // end of transform
    writeImageData(w, pixels, img.Bounds().Dx(), img.Bounds().Dy(), true, colorCacheBits)

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

    encoded := encodeImageData(pixels, colorCacheBits)
    histos := computeHistograms(encoded, colorCacheBits)

    var codes [][]huffmanCode
    for i := 0; i < 5; i++ {
        c := buildhuffmanCodes(histos[i], 16)
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