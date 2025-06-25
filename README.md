# colorspace
[![go.dev reference](https://pkg.go.dev/badge/github.com/soypat/colorspace)](https://pkg.go.dev/github.com/soypat/colorspace)
[![Go Report Card](https://goreportcard.com/badge/github.com/soypat/colorspace)](https://goreportcard.com/report/github.com/soypat/colorspace)
[![codecov](https://codecov.io/gh/soypat/colorspace/branch/main/graph/badge.svg)](https://codecov.io/gh/soypat/colorspace)
[![Go](https://github.com/soypat/colorspace/actions/workflows/go.yml/badge.svg)](https://github.com/soypat/colorspace/actions/workflows/go.yml)
[![sourcegraph](https://sourcegraph.com/github.com/soypat/colorspace/-/badge.svg)](https://sourcegraph.com/github.com/soypat/colorspace?badge)

colorspace implements different color space logic to allow for conversion from colorspace to colorspace and interpolation within each colorspace.

How to install package with newer versions of Go (+1.16):
```sh
go mod download github.com/soypat/colorspace@latest
```

## Linear interpolation example
Shown in each image are 5 different color gradients generated with linear interpolation in each available colorspace. See [`examples/lerp`](./examples/lerp/lerp.go):
1. Topmost: **sRGB**. This is the naive linear interpolation
2. **Linear sRGB**. 
3. **CIE XYZ**
4. **OKLAB**
5. **OKLCH**. Designed to yield the most perceptively uniform gradient.


![greyred to blue lerp](./greyred-blue.png)
![red to blue lerp](./red-blue.png)
![white to blue lerp](./white-blue.png)
![white to black lerp](./white-black.png)
