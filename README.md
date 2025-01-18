[![Codecov Coverage](https://codecov.io/gh/HugoSmits86/nativewebp/branch/main/graph/badge.svg)](https://codecov.io/gh/HugoSmits86/nativewebp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# Native WebP for Go

This is a native WebP encoder written entirely in Go, with **no dependencies on libwebp** or other external libraries. Designed for performance and efficiency, this encoder generates smaller files than the standard Go PNG encoder and is approximately **40% faster** in execution.

Currently, the encoder supports only WebP lossless images (VP8L).

## Benchmark

We conducted a quick benchmark to showcase file size reduction and encoding performance. Using an image from Google’s WebP Lossless and Alpha Gallery, we compared the results of our nativewebp encoder with the standard PNG decoder.

<p align="center">
  <img src="https://www.gstatic.com/webp/gallery3/4.png">
  <br/>
  <sub>image source: https://developers.google.com/speed/webp/gallery2</sub>
  <br/><br/>
  <table align="center">
    <tr>
      <th></th>
      <th>PNG encoder</th>
      <th>nativeWebP encoder</th>
      <th>reduction</th>
    </tr>
    <tr>
      <td>file size</td>
      <td>54kb</td>
      <td>50kb</td>
      <td>7% smaller</td>
    </tr>
    <tr>
      <td>encoding time</td>
      <td>5651958 ns/op</td>
      <td>1950541 ns/op</td>
      <td>65% faster</td>
    </tr>
  </table>
</p>

## Installation

To install the nativewebp package, use the following command:
```Bash
go get github.com/HugoSmits86/nativewebp
```
## Usage

Here’s a simple example of how to encode an image:
```Go
file, err := os.Create(name)
if err != nil {
  log.Fatalf("Error creating file %s: %v", name, err)
}
defer file.Close()

err = nativewebp.Encode(file, img)
if err != nil {
  log.Fatalf("Error encoding image to WebP: %v", err)
}
```
