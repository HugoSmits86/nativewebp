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
    "errors"
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
    data, err := io.ReadAll(r)
    if err != nil {
        return nil, err
    }

    if len(data) < 24 || string(data[:4]) != "RIFF" || string(data[8:12]) != "WEBP" {
        return nil, errors.New("webp: invalid format")
    }

    if string(data[12:16]) == "VP8X" && string(data[30:34]) == "VP8L" {
        if data[20] & (1 << 1) == 1 || data[20] & (1 << 5) == 1 {
            return nil, errors.New("webp: invalid format")
        }

        size := int(binary.LittleEndian.Uint32(data[16:20]))
        data = append(data[:12], data[20 + size:]...)
        binary.LittleEndian.PutUint32(data[4:8], uint32(len(data) - 8))
    }
   

    return decoderWebP.Decode(bytes.NewReader(data))
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