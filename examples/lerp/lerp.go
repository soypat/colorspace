package main

import (
	"image/color"

	colorspace "github.com/YOURUSER/YOURREPONAME"
)

func main() {

}

var lerps = map[string]func(c1, c2 color.Color, v float32) color.Color{
	"lin-sRGB": colorspace.LerpLSRGB,
	"CIE-XYZ":  colorspace.LerpCIEXYZ,
	"OKLAB":    colorspace.LerpOKLAB,
	"OKLCH":    colorspace.LerpOKLCH,
}
