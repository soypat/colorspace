package colorspace

import (
	"image/color"
	"math/rand"
	"testing"
)

func TestBasic(t *testing.T) {
	red := SRGB{R: 1, G: 0, B: 0}
	redlsrgb := red.LsRGB()
	wantlsrgb := LSRGB{R: 1, G: 0, B: 0}
	if redlsrgb != wantlsrgb {
		t.Errorf("lsrgb for red mismatch, want %v, got %v", wantlsrgb, redlsrgb)
	}
	redxyz := redlsrgb.CIEXYZ()
	wantredxyz := CIEXYZ{X: 0.41239079926595934, Y: 0.21263900587151027, Z: 0.01933081871559182}
	if redxyz != wantredxyz {
		t.Errorf("xyz for red not match: want %v, got %v", wantredxyz, redxyz)
	}
	redoklab := redxyz.OKLAB()
	expectoklab := OKLAB{
		L: 0.6279553639214311,
		A: 0.2248630684262744,
		B: 0.125846277330585,
	}
	if expectoklab.DeltaE(redoklab) > 0.0001 {
		t.Errorf("mismatch oklab for red: got %v, want %v", redoklab, expectoklab)
	}
}

func TestColor(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	palette := jet
	for i := 0; i < 10; i++ {
		idx1 := rng.Intn(len(palette))
		idx2 := rng.Intn(len(palette))
		v := rng.Float32()
		c := LerpOKLAB(palette[idx1], palette[idx2], v)
		_ = c
	}
}

var jet = color.Palette{
	color.RGBA64{R: 0x8180, G: 0xe0e, B: 0x707, A: 0xffff},
	color.RGBA64{R: 0xa3a3, G: 0x1413, B: 0x808, A: 0xffff},
	color.RGBA64{R: 0xc1c1, G: 0x1f1f, B: 0xd0c, A: 0xffff},
	color.RGBA64{R: 0xc6c6, G: 0x2221, B: 0xe0e, A: 0xffff},
	color.RGBA64{R: 0xe1e1, G: 0x3636, B: 0x1616, A: 0xffff},
	color.RGBA64{R: 0xebeb, G: 0x403f, B: 0x1918, A: 0xffff},
	color.RGBA64{R: 0xffff, G: 0x6464, B: 0x2726, A: 0xffff},
	color.RGBA64{R: 0xffff, G: 0x8282, B: 0x3434, A: 0xffff},
	color.RGBA64{R: 0xffff, G: 0x9f9e, B: 0x4141, A: 0xffff},
	color.RGBA64{R: 0xffff, G: 0xa2a1, B: 0x4242, A: 0xffff},
	color.RGBA64{R: 0xffff, G: 0xc1c1, B: 0x4e4d, A: 0xffff},
	color.RGBA64{R: 0xdedd, G: 0xe0df, B: 0x5150, A: 0xffff},
	color.RGBA64{R: 0xbfbf, G: 0xf1f0, B: 0x5352, A: 0xffff},
	color.RGBA64{R: 0x9d9c, G: 0xf9f9, B: 0x5857, A: 0xffff},
	color.RGBA64{R: 0x7d7c, G: 0xf9f9, B: 0x6160, A: 0xffff},
	color.RGBA64{R: 0x5e5d, G: 0xf9f9, B: 0x6d6d, A: 0xffff},
	color.RGBA64{R: 0x3b3b, G: 0xf9f9, B: 0x7979, A: 0xffff},
	color.RGBA64{R: 0x1a1a, G: 0xf9f9, B: 0x8180, A: 0xffff},
	color.RGBA64{R: 0x0, G: 0xdcdb, B: 0xc2c1, A: 0xffff},
	color.RGBA64{R: 0x0, G: 0xbfbf, B: 0xe5e4, A: 0xffff},
	color.RGBA64{R: 0x0, G: 0x9f9e, B: 0xf9f9, A: 0xffff},
	color.RGBA64{R: 0x201f, G: 0x8787, B: 0xf5f4, A: 0xffff},
	color.RGBA64{R: 0x2726, G: 0x7f7e, B: 0xefef, A: 0xffff},
	color.RGBA64{R: 0x3636, G: 0x5e5d, B: 0xcaca, A: 0xffff},
	color.RGBA64{R: 0x3838, G: 0x5554, B: 0xbebd, A: 0xffff},
	color.RGBA64{R: 0x3a3a, G: 0x3d3d, B: 0x9393, A: 0xffff},
	color.RGBA64{R: 0x3838, G: 0x3131, B: 0x7d7c, A: 0xffff},
	color.RGBA64{R: 0x3636, G: 0x1f1f, B: 0x5656, A: 0xffff},
	color.RGBA64{R: 0x3333, G: 0x1313, B: 0x3938, A: 0xffff},
}
