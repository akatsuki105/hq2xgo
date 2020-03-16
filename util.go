package hq2x

import "image/color"

func colorToRGB(c color.Color) (uint8, uint8, uint8) {
	r, g, b, _ := c.RGBA()
	return uint8(r >> 24), uint8(g >> 24), uint8(b >> 24)
}
