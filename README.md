[![Codecov Coverage](https://codecov.io/gh/HugoSmits86/nativewebp/branch/main/graph/badge.svg)](https://codecov.io/gh/HugoSmits86/nativewebp)
[![Go Reference](https://pkg.go.dev/badge/github.com/HugoSmits86/nativewebp.svg)](https://pkg.go.dev/github.com/HugoSmits86/nativewebp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# Native WebP for Go

This is a native WebP encoder written entirely in Go, with **no dependencies on libwebp** or other external libraries. Designed for performance and efficiency, this encoder generates smaller files than the standard Go PNG encoder and is approximately **50% faster** in execution.

Currently, the encoder supports only WebP lossless images (VP8L).

## Benchmark

We conducted a quick benchmark to showcase file size reduction and encoding performance. Using an image from Google’s WebP Lossless and Alpha Gallery, we compared the results of our nativewebp encoder with the standard PNG decoder.
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
    <td>121kb</td>
    <td>105kb</td>
    <td>13% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>14170403 ns/op</td>
    <td>5389776 ns/op</td>
    <td>62% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/2.png" height="100px"></p></td>
    <td>file size</td>
    <td>48kb</td>
    <td>38kb</td>
    <td>21% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>10662832 ns/op</td>
    <td>3760902 ns/op</td>
    <td>65% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/3.png" height="100px"></p></td>
    <td>file size</td>
    <td>238</td>
    <td>215</td>
    <td>10% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>30952147 ns/op</td>
    <td>16371708 ns/op</td>
    <td>47% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/4.png" height="60px"></p></td>
    <td>file size</td>
    <td>53kb</td>
    <td>43kb</td>
    <td>19% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>4511737 ns/op</td>
    <td>2181801 ns/op</td>
    <td>52% faster</td>
  </tr>
  <tr>
    <td rowspan="2" height="110px"><p align="center"><img src="https://www.gstatic.com/webp/gallery3/5.png" height="100px"></p></td>
    <td>file size</td>
    <td>140kb</td>
    <td>137kb</td>
    <td>2% smaller</td>
  </tr>
  <tr>
    <td>encoding time</td>
    <td>11045284 ns/op</td>
    <td>4850678 ns/op</td>
    <td>56% faster</td>
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
