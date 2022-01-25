package hq2x

import "image/color"

func _interp1(a, b uint8) uint8 {
	return uint8((uint(a)*3 + uint(b)) / 4)
}

func interp1(a, b color.RGBA) color.RGBA {
	R := _interp1(a.R, b.R)
	G := _interp1(a.G, b.G)
	B := _interp1(a.B, b.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func _interp2(a, b, c uint8) uint8 {
	return uint8((uint(a)*2 + uint(b) + uint(c)) / 4)
}

func interp2(a, b, c color.RGBA) color.RGBA {
	R := _interp2(a.R, b.R, c.R)
	G := _interp2(a.G, b.G, c.G)
	B := _interp2(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func _interp5(a, b uint8) uint8 {
	return uint8((uint(a) + uint(b)) / 2)
}

func interp5(a, b color.RGBA) color.RGBA {
	R := _interp5(a.R, b.R)
	G := _interp5(a.G, b.G)
	B := _interp5(a.B, b.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func _interp6(a, b, c uint8) uint8 {
	return uint8((uint(a)*5 + uint(b)*2 + uint(c)) / 8)
}

func interp6(a, b, c color.RGBA) color.RGBA {
	R := _interp6(a.R, b.R, c.R)
	G := _interp6(a.G, b.G, c.G)
	B := _interp6(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func _interp7(a, b, c uint8) uint8 {
	return uint8((uint(a)*6 + uint(b) + uint(c)) / 8)
}

func interp7(a, b, c color.RGBA) color.RGBA {
	R := _interp7(a.R, b.R, c.R)
	G := _interp7(a.G, b.G, c.G)
	B := _interp7(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func _interp9(a, b, c uint8) uint8 {
	return uint8((uint(a)*2 + uint(b)*3 + uint(c)*3) / 8)
}

func interp9(a, b, c color.RGBA) color.RGBA {
	R := _interp9(a.R, b.R, c.R)
	G := _interp9(a.G, b.G, c.G)
	B := _interp9(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}

func _interp10(a, b, c uint8) uint8 {
	return uint8((uint(a)*14 + uint(b) + uint(c)) / 16)
}

func interp10(a, b, c color.RGBA) color.RGBA {
	R := _interp10(a.R, b.R, c.R)
	G := _interp10(a.G, b.G, c.G)
	B := _interp10(a.B, b.B, c.B)
	return color.RGBA{
		R: R,
		G: G,
		B: B,
	}
}
