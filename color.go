package colorspace

import (
	"image/color"

	"github.com/chewxy/math32"
	"github.com/soypat/geometry/ms1"
	"github.com/soypat/geometry/ms3"
)

const undefinedHue = 0.0

var (
	// standard white points, defined by 4-figure CIE x,y chromaticities
	d50 = ms3.Vec{X: 0.3457 / 0.3585, Y: 1, Z: (1.0 - 0.3457 - 0.3585) / 0.3585}
	d65 = ms3.Vec{X: 0.3127 / 0.3290, Y: 1, Z: (1.0 - 0.3127 - 0.3290) / 0.3290}

	// Transposed due to being defined in column major format.
	linSRGBToXYZ = ms3.NewMat3([]float32{
		506752. / 1228815, 87881. / 245763, 12673. / 70218,
		87098. / 409605, 175762. / 245763, 12673. / 175545,
		7918. / 409605, 87881. / 737289, 1001167. / 1053270,
	})
	xyzToLinSRGB = ms3.NewMat3([]float32{12831. / 3959, -329. / 214, -1974. / 3959,
		-851781. / 878810, 1648619. / 878810, 36519. / 878810,
		705. / 12673, -2585. / 12673, 705. / 667,
	})
	d65Tod50 = ms3.NewMat3([]float32{1.0479297925449969, 0.022946870601609652, -0.05019226628920524,
		0.02962780877005599, 0.9904344267538799, -0.017073799063418826,
		-0.009243040646204504, 0.015055191490298152, 0.7518742814281371})
	d50Tod65 = ms3.NewMat3([]float32{0.955473421488075, -0.02309845494876471, 0.06325924320057072,
		-0.0283697093338637, 1.0099953980813041, 0.021041441191917323,
		0.012314014864481998, -0.020507649298898964, 1.330365926242124})
	xyzToLMS = ms3.NewMat3([]float32{0.8190224379967030, 0.3619062600528904, -0.1288737815209879,
		0.0329836539323885, 0.9292868615863434, 0.0361446663506424,
		0.0481771893596242, 0.2642395317527308, 0.6335478284694309})
	lmsToOKLAB = ms3.NewMat3([]float32{0.2104542683093140, 0.7936177747023054, -0.0040720430116193,
		1.9779985324311684, -2.4285922420485799, 0.4505937096174110,
		0.0259040424655478, 0.7827717124575296, -0.8086757549230774})
	lmsToXYZ = ms3.NewMat3([]float32{1.2268798758459243, -0.5578149944602171, 0.2813910456659647,
		-0.0405757452148008, 1.1122868032803170, -0.0717110580655164,
		-0.0763729366746601, -0.4214933324022432, 1.5869240198367816})
	oklabToLMS = ms3.NewMat3([]float32{1.0000000000000000, 0.3963377773761749, 0.2158037573099136,
		1.0000000000000000, -0.1055613458156586, -0.0638541728258133,
		1.0000000000000000, -0.0894841775298119, -1.2914855480194092})
)

// OKLAB is a uniform color space for device independent coloring designed to improve preceptual uniformity,
// hue and lightness prediction, color blending and usability regarding numerical stability.
type OKLAB struct {
	// Preceptual lightness. See [LAB]
	L float32
	// A and B for opposite channels of the four unique hues. unbounded but in practice ranging from -0.5 to 0.5.
	// CSS assigns ±100% to ±0.4 for both.
	A float32
	B float32
}

// OKLCH is cylindrical representation of [OKLAB] color space.
type OKLCH struct {
	L float32 // Perceptual luminosity. Same as for [OKLAB].
	C float32 // Chroma. Defines intensity of hue.
	H float32 // Hue in degrees.
}

// CIELCH is the cylindiracl hue color space representation of CIELAB.
type CIELCH struct {
	L float32 // Perceptual luminosity. Same as for [CIELAB].
	C float32 // chroma.
	H float32 // hue.
}

// CIELAB or also known as LAB, is a color model defined by the international commission on illumination (CIE) in 1976.
// It is designed so that a given numerical change always corresponds to a similar preceived change in color.
// Since a* and b* axes are unbounded a correct CIELAB color may not be representable in sRGB gamut.
type CIELAB struct {
	// L* (L-star) Perceptual Lightness calcuilated using the cube root of relative luminance with an offset near black.
	// Defines black at 0 and white at 1.
	L float32
	// a* axis (unbounded) Varies greenish appearance.
	A float32
	// b* axis (unbounded) varies red/green/yellow to blue.
	B float32
}

// SRGB usually denoted as sRGB (standard Red-Green-Blue) is a color space for use on printers, monitors and the world wide web.
// This is to say most monitors/printers receive color data as a triplets for each pixel representing redness, greeness and blueness.
type SRGB struct {
	R float32 // Red.
	G float32 // Green.
	B float32 // Blue.
}

// LSRGB is linear-light (un-companded) color space.
type LSRGB struct {
	R float32 // Red.
	G float32 // Green.
	B float32 // Blue.
}

// CIEXYZ refers to the 1931 CIE color space defined such that a mixture between two colors
// in some proportion lies on the line between those two colors in this space. One disadvantage
// of this model is that it is not perceptually uniform. The disadvantage is remedied in subsequent color models
// such as CIELUV and [CIELAB].
type CIEXYZ struct {
	X, Y, Z float32
}

func LerpSRGB(c1, c2 color.Color, v float32) color.Color {
	o1 := ColorToSRGB(c1)
	o2 := ColorToSRGB(c2)
	return o1.Lerp(o2, v)
}

func LerpLSRGB(c1, c2 color.Color, v float32) color.Color {
	o1 := ColorToSRGB(c1).LsRGB()
	o2 := ColorToSRGB(c2).LsRGB()
	return o1.Lerp(o2, v).ClipToGamut().SRGB()
}

func LerpCIEXYZ(c1, c2 color.Color, v float32) color.Color {
	o1 := ColorToSRGB(c1).LsRGB().CIEXYZ()
	o2 := ColorToSRGB(c2).LsRGB().CIEXYZ()
	return o1.Lerp(o2, v).LsRGB().ClipToGamut().SRGB()
}

func LerpOKLAB(c1, c2 color.Color, v float32) color.Color {
	o1 := ColorToSRGB(c1).LsRGB().CIEXYZ().OKLAB()
	o2 := ColorToSRGB(c2).LsRGB().CIEXYZ().OKLAB()
	lch := o1.Lerp(o2, v).OKLCH()
	mapped := lch.AutoGamutMapping()
	return mapped.OKLAB().CIEXYZ().LsRGB().ClipToGamut().SRGB()
}

func LerpOKLCH(c1, c2 color.Color, v float32) color.Color {
	o1 := ColorToSRGB(c1).LsRGB().CIEXYZ().OKLAB().OKLCH()
	o2 := ColorToSRGB(c2).LsRGB().CIEXYZ().OKLAB().OKLCH()
	mapped := o1.Lerp(o2, v).AutoGamutMapping()
	result := mapped.OKLAB().CIEXYZ().LsRGB().ClipToGamut().SRGB()
	return result
}

// ColorToSRGB converts the color to [SRGB] discarding the opacity/alpha (A) field.
func ColorToSRGB(c color.Color) SRGB {
	r, g, b, _ := c.RGBA()
	return SRGB{
		R: float32(r) / 0xffff,
		G: float32(g) / 0xffff,
		B: float32(b) / 0xffff,
	}
}

// transferFunc is the gamma function.
func transferFunc(v float32) float32 {
	sign := math32.Copysign(1, v)
	abs := math32.Abs(v)
	if abs <= 0.04045 {
		return v / 12.92
	}
	return sign * math32.Pow((abs+0.055)/1.055, 2.4)
}

// invTransferFunc is the inverse gamma function as defined by IEC2003.
func invTransferFunc(v float32) float32 {
	sign := math32.Copysign(1, v)
	abs := math32.Abs(v)
	if abs <= 0.0031308 {
		return 12.92 * v
	}
	return sign * (1.055*math32.Pow(abs, 1/2.4) - 0.055)
}

func (c SRGB) LsRGB() LSRGB {
	return LSRGB{
		R: transferFunc(c.R),
		G: transferFunc(c.G),
		B: transferFunc(c.B),
	}
}

func (c LSRGB) SRGB() SRGB {
	return SRGB{
		R: invTransferFunc(c.R),
		G: invTransferFunc(c.G),
		B: invTransferFunc(c.B),
	}
}

func (c LSRGB) Vec() ms3.Vec { return ms3.Vec{X: c.R, Y: c.G, Z: c.B} }
func (c CIEXYZ) Vec() ms3.Vec {
	return ms3.Vec{X: c.X, Y: c.Y, Z: c.Z}
}
func (c CIELAB) Vec() ms3.Vec { return ms3.Vec{X: c.L, Y: c.A, Z: c.B} }
func (c CIELCH) Vec() ms3.Vec { return ms3.Vec{X: c.L, Y: c.C, Z: c.H} }
func (c OKLCH) Vec() ms3.Vec  { return ms3.Vec{X: c.L, Y: c.C, Z: c.H} }
func (c OKLAB) Vec() ms3.Vec  { return ms3.Vec{X: c.L, Y: c.A, Z: c.B} }

func (c LSRGB) CIEXYZ() CIEXYZ {
	v := ms3.MulMatVec(linSRGBToXYZ, c.Vec())
	return CIEXYZ{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
	}
}

func (c CIEXYZ) LsRGB() LSRGB {
	v := ms3.MulMatVec(xyzToLinSRGB, c.Vec())
	return LSRGB{R: v.X, G: v.Y, B: v.Z}
}

func (c SRGB) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R * 0xffff)
	g = uint32(c.G * 0xffff)
	b = uint32(c.B * 0xffff)
	return r, g, b, 0xffff
}

func (c CIEXYZ) OKLAB() OKLAB {
	lms := ms3.MulMatVec(xyzToLMS, c.Vec())

	v := ms3.MulMatVec(lmsToOKLAB, ms3.Vec{
		X: math32.Cbrt(lms.X),
		Y: math32.Cbrt(lms.Y),
		Z: math32.Cbrt(lms.Z),
	})
	return OKLAB{
		L: v.X,
		A: v.Y,
		B: v.Z,
	}
}

func (c OKLAB) CIEXYZ() CIEXYZ {
	LMSnl := ms3.MulMatVec(oklabToLMS, c.Vec())
	v := ms3.MulMatVec(lmsToXYZ, ms3.Vec{
		X: LMSnl.X * LMSnl.X * LMSnl.X,
		Y: LMSnl.Y * LMSnl.Y * LMSnl.Y,
		Z: LMSnl.Z * LMSnl.Z * LMSnl.Z,
	})
	return CIEXYZ{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
	}
}

func (c OKLAB) OKLCH() OKLCH {
	const eps = 0.000004
	hue := math32.Atan2(c.B, c.A) * 180 / math32.Pi
	chroma := math32.Sqrt(c.A*c.A + c.B*c.B)
	if hue < 0 {
		hue += 360
	}
	if chroma <= eps {
		hue = undefinedHue
	}
	return OKLCH{
		L: c.L,
		C: chroma,
		H: hue,
	}
}

func (c OKLCH) AutoGamutMapping() OKLCH {
	// Early return for Lightness exceed range.
	origin := c
	if origin.L < 0 || origin.L > 100 {
		return OKLCH{
			L: math32.Min(math32.Max(origin.L, 0), 1),
			C: 0,
			H: 0,
		}
	}
	const (
		JND = 0.02
		eps = 0.0001
	)
	current := origin
	clipped := current.OKLAB().CIEXYZ().LsRGB().ClipToGamut()
	E := origin.OKLAB().DeltaE(clipped.CIEXYZ().OKLAB())
	if E < JND {
		return clipped.CIEXYZ().OKLAB().OKLCH()
	}
	var cmin, cmax float32 = 0, origin.C
	minInGamut := true
	for cmax-cmin > eps {
		chroma := 0.5 * (cmin + cmax)
		current.C = chroma
		currentRGB := current.OKLAB().CIEXYZ().LsRGB()
		if minInGamut && currentRGB.InGamut() {
			cmin = chroma
			minInGamut = OKLCH{L: current.L, C: chroma, H: current.H}.OKLAB().CIEXYZ().LsRGB().InGamut()
			continue
		}
		clipped = currentRGB.ClipToGamut()
		E = clipped.CIEXYZ().OKLAB().DeltaE(current.OKLAB())
		if E < JND {
			if JND-E < eps {
				return clipped.CIEXYZ().OKLAB().OKLCH()
			}
			minInGamut = false
			cmin = chroma
			minInGamut = OKLCH{L: current.L, C: chroma, H: current.H}.OKLAB().CIEXYZ().LsRGB().InGamut()
		} else {
			cmax = chroma
		}
	}
	return clipped.CIEXYZ().OKLAB().OKLCH()
}

func (c LSRGB) InGamut() bool {
	return c.R <= 1 && c.G <= 1 && c.B <= 1 && c.R >= 0 && c.G >= 0 && c.B >= 0
}

func (c SRGB) InGamut() bool {
	return c.R <= 1 && c.G <= 1 && c.B <= 1 && c.R >= 0 && c.G >= 0 && c.B >= 0
}

func (c LSRGB) ClipToGamut() LSRGB {
	return LSRGB{
		R: ms1.Clamp(c.R, 0, 1),
		G: ms1.Clamp(c.G, 0, 1),
		B: ms1.Clamp(c.B, 0, 1),
	}
}

func (c SRGB) ClipToGamut() SRGB {
	return SRGB{
		R: ms1.Clamp(c.R, 0, 1),
		G: ms1.Clamp(c.G, 0, 1),
		B: ms1.Clamp(c.B, 0, 1),
	}
}

func (c OKLCH) OKLAB() OKLAB {
	sin, cos := math32.Sincos(c.H * math32.Pi / 180)
	return OKLAB{
		L: c.L,
		A: c.C * cos,
		B: c.C * sin,
	}
}

func (c CIEXYZ) CIELAB() CIELAB {
	// Assuming XYZ is relative to D50, convert to CIE Lab
	// from CIE standard, which now defines these as a rational fraction
	const (
		ε = 216. / 24389 // 6^3/29^3
		κ = 24389. / 27  // 29^3/3^3
	)
	// compute xyz, which is XYZ scaled relative to reference white
	xyz := ms3.MulElem(c.Vec(), d50)
	f := func(x float32) float32 {
		if x > ε {
			return math32.Cbrt(x)
		}
		return (κ*x + 16) / 116
	}
	return CIELAB{
		L: 116*f(xyz.Y) - 16,
		A: 500 * (f(xyz.X) - f(xyz.Y)),
		B: 200 * (f(xyz.Y) - f(xyz.Z)),
	}
}

func (c CIELAB) CIELCH() CIELCH {
	const eps = 0.0015
	chroma := math32.Sqrt(c.A*c.A + c.B*c.B)
	hue := math32.Atan2(c.B, c.A) * 180 / math32.Pi
	if hue < 0 {
		hue += 360
	}
	if chroma <= eps {
		hue = undefinedHue
	}
	return CIELCH{
		L: c.L,
		C: chroma,
		H: hue,
	}
}

func (c CIELAB) CIEXYZ() CIEXYZ {
	const κ = 24389. / 27  // 29^3/3^3
	const ε = 216. / 24389 // 6^3/29^3
	const ecbrt = 6. / 29
	f1 := (c.L + 16) / 116
	f0 := c.A/500 + f1
	f2 := f1 - c.B/200

	var xyz CIEXYZ
	if f0 > ecbrt {
		xyz.X = f0 * f0 * f0
	} else {
		xyz.X = (116*f0 - 16) / κ
	}
	if c.L > κ*ε {
		ycbrt := (c.L + 16) / 116
		xyz.Y = ycbrt * ycbrt * ycbrt
	} else {
		xyz.Y = (116*f0 - 16) / κ
	}
	if f2 > ecbrt {
		xyz.Z = f2 * f2 * f2
	} else {
		xyz.Z = (116*f2 - 16) / κ
	}
	// Compute XYZ by scaling xyz by reference white
	v := ms3.MulElem(xyz.Vec(), d50)
	return CIEXYZ{X: v.X, Y: v.Y, Z: v.Z}
}

func (reference OKLAB) DeltaE(sample OKLAB) float32 {
	e := ms3.Sub(reference.Vec(), sample.Vec())
	return math32.Sqrt(ms3.Dot(e, e))
}

func (c CIELCH) CIELAB() CIELAB {
	sin, cos := math32.Sincos(c.H * math32.Pi / 180)
	return CIELAB{
		L: c.L,
		A: c.C * sin,
		B: c.C * cos,
	}
}

func (from OKLCH) Lerp(to OKLCH, v float32) OKLCH {
	return OKLCH(CIELCH(from).Lerp(CIELCH(to), v))
}

func (from CIELCH) Lerp(to CIELCH, v float32) CIELCH {
	// First handle achromatic or "powerless hue" colors.
	const eps = 0.000004
	fromPowerless := from.C < eps
	toPowerless := to.C < eps
	if fromPowerless || toPowerless {
		if fromPowerless && toPowerless {
			// Both colors are grayish, just ignore hue interpolation entirely.
			return CIELCH{
				L: ms1.Interp(from.L, to.L, v),
				C: 0,
				H: undefinedHue,
			}
		} else if !toPowerless {
			from.H = to.H
		} else {
			to.H = from.H
		}
	}
	// These branches are security, they are not taken regularly it seems.
	Hdelta := to.H - from.H // Hue is cylindrical in space [0..360], so the difference between 1 and 359 is 2, not 358.
	if Hdelta > 180 {
		Hdelta -= 360
	} else if Hdelta < -180 {
		Hdelta += 360
	}
	H := from.H + v*Hdelta
	if H > 360 {
		H -= 360
	} else if H < 0 {
		H += 360
	}
	return CIELCH{
		L: ms1.Interp(from.L, to.L, v),
		H: H,
		C: ms1.Interp(from.C, to.C, v),
	}
}

func (from CIELAB) Lerp(to CIELAB, v float32) CIELAB {
	return CIELAB{
		L: ms1.Interp(from.L, to.L, v),
		A: ms1.Interp(from.A, to.A, v),
		B: ms1.Interp(from.B, to.B, v),
	}
}

func (from OKLAB) Lerp(to OKLAB, v float32) OKLAB {
	return OKLAB(CIELAB(from).Lerp(CIELAB(to), v))
}

func (from CIEXYZ) Lerp(to CIEXYZ, v float32) CIEXYZ {
	return CIEXYZ{
		X: ms1.Interp(from.X, to.X, v),
		Y: ms1.Interp(from.Y, to.Y, v),
		Z: ms1.Interp(from.Z, to.Z, v),
	}
}

func (from LSRGB) Lerp(to LSRGB, v float32) LSRGB {
	return LSRGB{
		R: ms1.Interp(from.R, to.R, v),
		G: ms1.Interp(from.G, to.G, v),
		B: ms1.Interp(from.B, to.B, v),
	}
}
func (from SRGB) Lerp(to SRGB, v float32) SRGB {
	return SRGB{
		R: ms1.Interp(from.R, to.R, v),
		G: ms1.Interp(from.G, to.G, v),
		B: ms1.Interp(from.B, to.B, v),
	}
}
