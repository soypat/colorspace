package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/soypat/colorspace"
)

var splitcolor = color.RGBA{R: 255, A: 255}

func main() {
	const width, height = 700, 50
	const dx = 1. / width
	imgHeight := (height + 1) * len(lerps) // 1 pixel for differentiating lerps.
	for _, crange := range ranges {
		img := image.NewRGBA64(image.Rect(0, 0, width, imgHeight))
		for ilerp, lerp := range lerps {
			flerp := lerp.F
			c1, c2 := crange.c1, crange.c2
			for ix := 0; ix < width; ix++ {
				x := float32(ix) * dx
				c := flerp(c1, c2, x)
				yoff := ilerp * (height + 1)
				for iy := yoff; iy < yoff+height; iy++ {
					img.Set(ix, iy, c)
				}
				img.Set(ix, yoff+height, splitcolor)
			}
		}
		fp, _ := os.Create(crange.name + ".png")
		png.Encode(fp, img)
		fp.Close()
	}
}

type colorrange struct {
	name   string
	c1, c2 color.RGBA
}

var ranges = []colorrange{
	{name: "white-black", c1: color.RGBA{R: 255, G: 255, B: 255}, c2: color.RGBA{}},
	{name: "white-blue", c1: color.RGBA{R: 255, G: 255, B: 255}, c2: color.RGBA{B: 255}},
	{name: "red-blue", c1: color.RGBA{R: 255}, c2: color.RGBA{B: 255}},
	{name: "greyred-blue", c1: color.RGBA{R: 160, G: 127, B: 127}, c2: color.RGBA{B: 255}},
}

type Lerp struct {
	Name string
	F    func(c1, c2 color.Color, v float32) color.Color
}

var lerps = []Lerp{
	{
		Name: "sRGB",
		F:    colorspace.LerpSRGB,
	},
	{
		Name: "lin-sRGB",
		F:    colorspace.LerpLSRGB,
	},
	{
		Name: "CIE-XYZ",
		F:    colorspace.LerpCIEXYZ,
	},
	{
		Name: "OKLAB",
		F:    colorspace.LerpOKLAB,
	},
	{
		Name: "OKLCH",
		F:    colorspace.LerpOKLCH,
	},
}
