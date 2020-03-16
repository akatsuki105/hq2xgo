package hq2x

import "image/color"

func interp1(a, b color.RGBA) color.RGBA {
	f := func(a, b uint8) uint8 {
		return uint8((uint(a)*3 + uint(b)) / 4)
	}

	R := f(a.R, b.R)
	G := f(a.G, b.G)
	B := f(a.B, b.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func interp2(a, b, c color.RGBA) color.RGBA {
	f := func(a, b, c uint8) uint8 {
		return uint8((uint(a)*2 + uint(b) + uint(c)) / 4)
	}

	R := f(a.R, b.R, c.R)
	G := f(a.G, b.G, c.G)
	B := f(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func interp5(a, b color.RGBA) color.RGBA {
	f := func(a, b uint8) uint8 {
		return uint8((uint(a) + uint(b)) / 2)
	}

	R := f(a.R, b.R)
	G := f(a.G, b.G)
	B := f(a.B, b.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func interp6(a, b, c color.RGBA) color.RGBA {
	f := func(a, b, c uint8) uint8 {
		return uint8((uint(a)*5 + uint(b)*2 + uint(c)) / 8)
	}

	R := f(a.R, b.R, c.R)
	G := f(a.G, b.G, c.G)
	B := f(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func interp7(a, b, c color.RGBA) color.RGBA {
	f := func(a, b, c uint8) uint8 {
		return uint8((uint(a)*6 + uint(b) + uint(c)) / 8)
	}

	R := f(a.R, b.R, c.R)
	G := f(a.G, b.G, c.G)
	B := f(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func interp9(a, b, c color.RGBA) color.RGBA {
	f := func(a, b, c uint8) uint8 {
		return uint8((uint(a)*2 + uint(b)*3 + uint(c)*3) / 8)
	}

	R := f(a.R, b.R, c.R)
	G := f(a.G, b.G, c.G)
	B := f(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func interp10(a, b, c color.RGBA) color.RGBA {
	f := func(a, b, c uint8) uint8 {
		return uint8((uint(a)*14 + uint(b) + uint(c)) / 16)
	}

	R := f(a.R, b.R, c.R)
	G := f(a.G, b.G, c.G)
	B := f(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}
