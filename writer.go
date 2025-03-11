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
    stream, hasAlpha, err := EncodeVP8L(img, o)
    if err != nil {
        return err
    }

    tmp := make([]byte, 4)
    buf := &bytes.Buffer{}

    if o != nil && o.UseExtendedFormat {
        buf.Write([]byte("VP8X"))
        binary.LittleEndian.PutUint32(tmp, 10)
        buf.Write(tmp)
    
        var flags uint32
        if hasAlpha {
            flags |= 0x00000010
        }
    
        binary.Write(buf, binary.LittleEndian, flags)
    
        dx := img.Bounds().Dx() - 1
        dy := img.Bounds().Dy() - 1
    
        buf.Write([]byte{byte(dx), byte(dx >> 8), byte(dx >> 16)})
        buf.Write([]byte{byte(dy), byte(dy >> 8), byte(dy >> 16)})
    }

    buf.Write([]byte("VP8L"))
    binary.LittleEndian.PutUint32(tmp, uint32(stream.Len()))
    buf.Write(tmp)
    buf.Write(stream.Bytes())

    w.Write([]byte("RIFF"))
    binary.LittleEndian.PutUint32(tmp, uint32(4 + buf.Len()))
    w.Write(tmp)

    w.Write([]byte("WEBP"))
    w.Write(buf.Bytes())

    return nil
}

func EncodeVP8L(img image.Image, o *Options) (*bytes.Buffer, bool, error) {
    if img == nil {
        return nil, false, errors.New("image is nil")
    }

    if img.Bounds().Dx() < 1 || img.Bounds().Dy() < 1 {
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
    
    s.AlignByte()

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