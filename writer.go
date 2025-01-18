package nativewebp

import (
	//------------------------------
	//general
	//------------------------------
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"slices"

	//------------------------------
	//imaging
	//------------------------------
	"image"
	"image/color"
	"image/draw"

	//------------------------------
	//errors
	//------------------------------
	//"log"
	"errors"
)

type Transform int

const (
	TransformPredict       = Transform(0)
	TransformSubGreen      = Transform(2)
	TransformColorIndexing = Transform(3)
)

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

	s := getBitWriter()
	defer putBitWriter(s)

	writeBitStreamHeader(s, rgba.Bounds(), !rgba.Opaque())

	var transforms [4]bool
	transforms[TransformPredict] = !isIndexed
	transforms[TransformSubGreen] = !isIndexed
	transforms[TransformColorIndexing] = isIndexed

	err := writeBitStreamData(s, rgba, 4, transforms)
	if err != nil {
		return err
	}

	s.AlignByte()

	if s.Buffer.Len()%2 != 0 {
		s.Buffer.Write([]byte{0x00})
	}

	writeWebPHeader(w, s.Buffer)

	if _, err := s.Buffer.WriteTo(w); err != nil {
		return fmt.Errorf("writing data: %w", err)
	}

	return nil
}

func writeWebPHeader(w io.Writer, b *bytes.Buffer) {
	w.Write([]byte("RIFF"))

	tmp := make([]byte, 4)
	binary.LittleEndian.PutUint32(tmp, uint32(12+b.Len()))
	w.Write(tmp)

	w.Write([]byte("WEBP"))
	w.Write([]byte("VP8L"))

	tmp = make([]byte, 4)
	binary.LittleEndian.PutUint32(tmp, uint32(b.Len()))
	w.Write(tmp)
}

func writeBitStreamHeader(w *BitWriter, bounds image.Rectangle, hasAlpha bool) {
	w.writeBits(0x2f, 8)

	w.writeBits(uint64(bounds.Dx()-1), 14)
	w.writeBits(uint64(bounds.Dy()-1), 14)

	if hasAlpha {
		w.writeBits(1, 1)
	} else {
		w.writeBits(0, 1)
	}

	w.writeBits(0, 3)
}

func writeBitStreamData(w *BitWriter, img image.Image, colorCacheBits int, transforms [4]bool) error {
	pixels, err := flatten(img)
	if err != nil {
		return err
	}

	if transforms[TransformColorIndexing] {
		w.writeBits(1, 1)
		w.writeBits(3, 2)

		pal, err := applyPaletteTransform(pixels)
		if err != nil {
			return err
		}

		w.writeBits(uint64(len(pal)-1), 8)
		writeImageData(w, pal, false, colorCacheBits)
	}

	if transforms[TransformSubGreen] {
		w.writeBits(1, 1)
		w.writeBits(2, 2)

		applySubtractGreenTransform(pixels)
	}

	if transforms[TransformPredict] {
		w.writeBits(1, 1)
		w.writeBits(0, 2)

		bits, blocks := applyPredictTransform(pixels, img.Bounds().Dx(), img.Bounds().Dy())

		w.writeBits(uint64(bits-2), 3)
		writeImageData(w, blocks, false, colorCacheBits)
	}

	w.writeBits(0, 1) // end of transform
	writeImageData(w, pixels, true, colorCacheBits)

	return nil
}

func writeImageData(w *BitWriter, pixels []color.NRGBA, isRecursive bool, colorCacheBits int) {
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

	var codes [][]HuffmanCode
	for i := 0; i < 5; i++ {
		c := buildHuffmanCodes(histos[i], 16)
		codes = append(codes, c)

		writeHuffmanCodes(w, c)
	}

	for i := 0; i < len(encoded); i++ {
		w.writeCode(codes[0][encoded[i+0]])
		if encoded[i+0] < 256 {
			w.writeCode(codes[1][encoded[i+1]])
			w.writeCode(codes[2][encoded[i+2]])
			w.writeCode(codes[3][encoded[i+3]])
			i += 3
		}
	}
}

func encodeImageData(pixels []color.NRGBA, colorCacheBits int) []uint16 {
	cache := make([]color.NRGBA, 1<<colorCacheBits)
	encoded := make([]uint16, len(pixels)*4)

	cnt := 0
	for _, p := range pixels {
		hash := 0
		if colorCacheBits > 0 {
			//hash formula including magic number 0x1e35a7bd comes directly from WebP specs!
			pack := uint32(p.A)<<24 | uint32(p.R)<<16 | uint32(p.G)<<8 | uint32(p.B)
			hash = int((pack * 0x1e35a7bd) >> (32 - colorCacheBits))

			if cache[hash] == p {
				encoded[cnt] = uint16(hash + 256 + 24)
				cnt++
				continue
			}

			cache[hash] = p
		}

		encoded[cnt+0] = uint16(p.G)
		encoded[cnt+1] = uint16(p.R)
		encoded[cnt+2] = uint16(p.B)
		encoded[cnt+3] = uint16(p.A)
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
		make([]int, 256+24+c),
		make([]int, 256),
		make([]int, 256),
		make([]int, 256),
		make([]int, 40),
	}

	for i := 0; i < len(pixels); i++ {
		histos[0][pixels[i]]++
		if pixels[i] < 256 {
			histos[1][pixels[i+1]]++
			histos[2][pixels[i+2]]++
			histos[3][pixels[i+3]]++
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

	pixels := make([]color.NRGBA, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := rgba.PixOffset(x, y)
			s := rgba.Pix[i : i+4 : i+4]

			pixels[y*w+x].R = uint8(s[0])
			pixels[y*w+x].G = uint8(s[1])
			pixels[y*w+x].B = uint8(s[2])
			pixels[y*w+x].A = uint8(s[3])
		}
	}

	return pixels, nil
}

func applySubtractGreenTransform(pixels []color.NRGBA) {
	for i, _ := range pixels {
		pixels[i].R = pixels[i].R - pixels[i].G
		pixels[i].B = pixels[i].B - pixels[i].G
	}
}

func applyPredictTransform(pixels []color.NRGBA, width, height int) (int, []color.NRGBA) {
	tileBits := 4
	tileSize := 1 << tileBits
	bw := (width + tileSize - 1) / tileSize
	bh := (height + tileSize - 1) / tileSize

	blocks := make([]color.NRGBA, bw*bh)
	deltas := make([]color.NRGBA, width*height)

	best := 1
	for y := 0; y < bh; y++ {
		for x := 0; x < bw; x++ {
			mx := min((x+1)<<tileBits, width)
			my := min((y+1)<<tileBits, height)

			for tx := x << tileBits; tx < mx; tx++ {
				for ty := y << tileBits; ty < my; ty++ {
					d := applyFilter(pixels, width, tx, ty, best)

					off := ty*width + tx
					deltas[off] = color.NRGBA{
						R: uint8(pixels[off].R - d.R),
						G: uint8(pixels[off].G - d.G),
						B: uint8(pixels[off].B - d.B),
						A: uint8(pixels[off].A - d.A),
					}
				}
			}

			blocks[y*bw+x] = color.NRGBA{0, byte(best), 0, 255}
		}
	}

	copy(pixels, deltas)

	return tileBits, blocks
}

func applyFilter(pixels []color.NRGBA, width, x, y, prediction int) color.NRGBA {
	if x == 0 && y == 0 {
		return color.NRGBA{0, 0, 0, 255}
	} else if x == 0 {
		return pixels[(y-1)*width+x]
	} else if y == 0 {
		return pixels[y*width+(x-1)]
	}

	t := pixels[(y-1)*width+x]
	l := pixels[y*width+(x-1)]

	tl := pixels[(y-1)*width+(x-1)]
	tr := pixels[(y-1)*width+(x+1)]

	avarage2 := func(a, b color.NRGBA) color.NRGBA {
		return color.NRGBA{
			uint8((int(a.R) + int(b.R)) / 2),
			uint8((int(a.G) + int(b.G)) / 2),
			uint8((int(a.B) + int(b.B)) / 2),
			uint8((int(a.A) + int(b.A)) / 2),
		}
	}

	filters := []func(t, l, tl, tr color.NRGBA) color.NRGBA{
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
			pl := math.Abs(pa-float64(l.A)) + math.Abs(pr-float64(l.R)) +
				math.Abs(pg-float64(l.G)) + math.Abs(pb-float64(l.B))
			pt := math.Abs(pa-float64(t.A)) + math.Abs(pr-float64(t.R)) +
				math.Abs(pg-float64(t.G)) + math.Abs(pb-float64(t.B))

			if pl < pt {
				return l
			}

			return t
		},
		func(t, l, tl, tr color.NRGBA) color.NRGBA {
			return color.NRGBA{
				uint8(max(min(int(l.R)+int(t.R)-int(tl.R), 255), 0)),
				uint8(max(min(int(l.G)+int(t.G)-int(tl.G), 255), 0)),
				uint8(max(min(int(l.B)+int(t.B)-int(tl.B), 255), 0)),
				uint8(max(min(int(l.A)+int(t.A)-int(tl.A), 255), 0)),
			}
		},
		func(t, l, tl, tr color.NRGBA) color.NRGBA {
			a := avarage2(l, t)

			return color.NRGBA{
				uint8(max(min(int(a.R)+(int(a.R)-int(tl.R))/2, 255), 0)),
				uint8(max(min(int(a.G)+(int(a.G)-int(tl.G))/2, 255), 0)),
				uint8(max(min(int(a.B)+(int(a.B)-int(tl.B))/2, 255), 0)),
				uint8(max(min(int(a.A)+(int(a.A)-int(tl.A))/2, 255), 0)),
			}
		},
	}

	return filters[prediction](t, l, tl, tr)
}

func applyPaletteTransform(pixels []color.NRGBA) ([]color.NRGBA, error) {
	var pal []color.NRGBA
	for _, p := range pixels {
		if !slices.Contains(pal, p) {
			pal = append(pal, p)
		}

		if len(pal) > 256 {
			return nil, errors.New("palette exceeds 256 colors")
		}
	}

	for i, p := range pixels {
		pixels[i] = color.NRGBA{G: uint8(slices.Index(pal, p)), A: 255}
	}

	for i := len(pal) - 1; i > 0; i-- {
		pal[i] = color.NRGBA{
			R: pal[i].R - pal[i-1].R,
			G: pal[i].G - pal[i-1].G,
			B: pal[i].B - pal[i-1].B,
			A: pal[i].A - pal[i-1].A,
		}
	}

	return pal, nil
}
