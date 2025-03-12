[![Codecov Coverage](https://codecov.io/gh/HugoSmits86/nativewebp/branch/main/graph/badge.svg)](https://codecov.io/gh/HugoSmits86/nativewebp)
[![Go Reference](https://pkg.go.dev/badge/github.com/HugoSmits86/nativewebp.svg)](https://pkg.go.dev/github.com/HugoSmits86/nativewebp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# Native WebP for Go

This is a native WebP encoder written entirely in Go, with **no dependencies on libwebp** or other external libraries. Designed for performance and efficiency, this encoder generates smaller files than the standard Go PNG encoder and is approximately **50% faster** in execution.

Currently, the encoder supports only WebP lossless images (VP8L).

## Decoding Support

We provide WebP decoding through a wrapper around `golang.org/x/image/webp`, with an additional `DecodeIgnoreAlphaFlag` function to handle VP8X images where the alpha flag causes decoding issues.
## Benchmark

We conducted a quick benchmark to showcase file size reduction and encoding performance. Using an image from Google’s WebP Lossless and Alpha Gallery, we compared the results of our nativewebp encoder with the standard PNG encoder. <br/><br/>
For the PNG encoder, we applied the `png.BestCompression` setting to achieve the most competitive compression outcomes.
<br/><br/>

<table align="center">
  <tr>
    <th></th>
    <th></th>
    <th>PNG encoder</th>
    <th>nativeWebP encoder</th>
    <th>reduction</th>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/1.png" height="100px"></p></td>
    <td>file size</td>
    <td>120 kb</td>
    <td>96 kb</td>
    <td>20% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>42945049 ns/op</td>
    <td>27716447 ns/op</td>
    <td>35% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/2.png" height="100px"></p></td>
    <td>file size</td>
    <td>46 kb</td>
    <td>36 kb</td>
    <td>22% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>98509399 ns/op</td>
    <td>31461759 ns/op</td>
    <td>68% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/3.png" height="100px"></p></td>
    <td>file size</td>
    <td>236 kb</td>
    <td>194 kb</td>
    <td>18% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>178205535 ns/op</td>
    <td>102454192 ns/op</td>
    <td>43% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/4.png" height="60px"></p></td>
    <td>file size</td>
    <td>53 kb</td>
    <td>41 kb</td>
    <td>23% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>29088555 ns/op</td>
    <td>14959849 ns/op</td>
    <td>49% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/5.png" height="100px"></p></td>
    <td>file size</td>
    <td>139 kb</td>
    <td>123 kb</td>
    <td>12% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>63423995 ns/op</td>
    <td>21717392 ns/op</td>
    <td>66% faster</td>
  </tr>
</table>
<p align="center">
<sub>image source: https://developers.google.com/speed/webp/gallery2</sub>
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

err = nativewebp.Encode(file, img, nil)
if err != nil {
  log.Fatalf("Error encoding image to WebP: %v", err)
}
```
