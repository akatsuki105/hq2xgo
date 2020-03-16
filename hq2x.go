package hq2x

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

const (
	TOP_LEFT = iota
	TOP
	TOP_RIGHT
	LEFT
	CENTER
	RIGHT
	BOTTOM_LEFT
	BOTTOM
	BOTTOM_RIGHT
)

// HQ2x - Enlarge image by 2x with hq2x algorithm
func HQ2x(src *image.RGBA) (*image.RGBA, error) {
	srcX, srcY := src.Bounds().Dx(), src.Bounds().Dy()

	dest := image.NewRGBA(image.Rect(0, 0, srcX*2, srcY*2))

	for x := 0; x < srcX; x++ {
		for y := 0; y < srcY; y++ {
			context := [9]color.RGBA{
				getPixel(src, x-1, y-1), getPixel(src, x, y-1), getPixel(src, x+1, y-1),
				getPixel(src, x-1, y), getPixel(src, x, y), getPixel(src, x+1, y),
				getPixel(src, x-1, y+1), getPixel(src, x, y+1), getPixel(src, x+1, y+1),
			}

			tl, tr, bl, br := hq2xPixel(context)
			dest.Set(x*2, y*2, tl)
			dest.Set(x*2+1, y*2, tr)
			dest.Set(x*2, y*2+1, bl)
			dest.Set(x*2+1, y*2+1, br)
		}
	}

	return dest, nil
}

func getPixel(src *image.RGBA, x, y int) color.RGBA {
	width, height := src.Bounds().Dx(), src.Bounds().Dy()

	if x < 0 {
		x = 0
	} else if x >= width {
		x = width - 1
	}

	if y < 0 {
		y = 0
	} else if y >= height {
		y = height - 1
	}

	return src.RGBAAt(x, y)
}

func hq2xPixel(context [9]color.RGBA) (tl, tr, bl, br color.RGBA) {
	yuvContext := [9]color.YCbCr{}
	yuvPixel := RGBAToYCbCr(context[CENTER])
	for i := 0; i <= 9; i++ {
		yuvContext[i] = RGBAToYCbCr(context[i])
	}

	contextFlag := newContextFlag()
	var pattern uint8
	for bit := 0; bit <= 9; bit++ {
		if bit != CENTER && !equalYuv(yuvContext[bit], yuvPixel) {
			pattern |= contextFlag[bit]
		}
	}

	switch pattern {
	case 0, 1, 4, 32, 128, 5, 132, 160, 33, 129, 36, 133, 164, 161, 37, 165:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 2, 34, 130, 162:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 16, 17, 48, 49:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 64, 65, 68, 69:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 8, 12, 136, 140:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 3, 35, 131, 163:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 6, 38, 134, 166:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 20, 21, 52, 53:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 144, 145, 176, 177:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp1(context[CENTER], context[BOTTOM])

	case 192, 193, 196, 197:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 96, 97, 100, 101:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 40, 44, 168, 172:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 9, 13, 137, 141:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 18, 50:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 80, 81:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 72, 76:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 10, 138:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 66:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 24:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 7, 39, 135:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 148, 149, 180:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp1(context[CENTER], context[BOTTOM])

	case 224, 228, 225:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 41, 169, 45:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 22, 54:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 208, 209:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 104, 108:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 11, 139:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 19, 51:
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tl = interp1(context[CENTER], context[LEFT])
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tl = interp6(context[CENTER], context[TOP], context[LEFT])
			tr = interp9(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 146, 178:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
			br = interp1(context[CENTER], context[BOTTOM])
		} else {
			tr = interp9(context[CENTER], context[TOP], context[RIGHT])
			br = interp6(context[CENTER], context[RIGHT], context[BOTTOM])
		}
		bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])

	case 84, 85:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP])
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			tr = interp6(context[CENTER], context[RIGHT], context[TOP])
			br = interp9(context[CENTER], context[RIGHT], context[BOTTOM])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])

	case 112, 113:
		tl = interp2(context[CENTER], context[LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			bl = interp1(context[CENTER], context[LEFT])
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			bl = interp6(context[CENTER], context[BOTTOM], context[LEFT])
			br = interp9(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 200, 204:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
			br = interp1(context[CENTER], context[RIGHT])
		} else {
			bl = interp9(context[CENTER], context[BOTTOM], context[LEFT])
			br = interp6(context[CENTER], context[BOTTOM], context[RIGHT])
		}

	case 73, 77:
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			tl = interp1(context[CENTER], context[TOP])
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			tl = interp6(context[CENTER], context[LEFT], context[TOP])
			bl = interp9(context[CENTER], context[BOTTOM], context[LEFT])
		}
		tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 42, 170:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
			bl = interp1(context[CENTER], context[BOTTOM])
		} else {
			tl = interp9(context[CENTER], context[LEFT], context[TOP])
			bl = interp6(context[CENTER], context[LEFT], context[BOTTOM])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 14, 142:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
			tr = interp1(context[CENTER], context[RIGHT])
		} else {
			tl = interp9(context[CENTER], context[LEFT], context[TOP])
			tr = interp6(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])

	case 67:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 70:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 28:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 152:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 194:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 98:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 56:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 25:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 26, 31:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 82, 214:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 88, 248:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 74, 107:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 27:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp1(context[CENTER], context[TOP_RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 86:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp1(context[CENTER], context[BOTTOM_RIGHT])

	case 216:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 106:
		tl = interp1(context[CENTER], context[TOP_LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 30:
		tl = interp1(context[CENTER], context[TOP_LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = context[CENTER]
		} else {
			tr = interp2(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 210:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp1(context[CENTER], context[TOP_RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = context[CENTER]
		} else {
			br = interp2(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 120:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = context[CENTER]
		} else {
			bl = interp2(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp1(context[CENTER], context[BOTTOM_RIGHT])
	case 75:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = context[CENTER]
		} else {
			tl = interp2(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 29:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 198:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 184:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 99:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 57:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 71:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 156:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 226:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 60:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 195:
		tl = interp1(context[CENTER], context[LEFT])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		br = interp1(context[CENTER], context[RIGHT])

	case 102:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[LEFT])
		tr = interp1(context[CENTER], context[RIGHT])
		bl = interp1(context[CENTER], context[LEFT])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[RIGHT])

	case 153:
		tl = interp1(context[CENTER], context[TOP])
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[TOP])
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[BOTTOM])
		br = interp1(context[CENTER], context[BOTTOM])

	case 58:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp1(context[CENTER], context[BOTTOM])
		br = interp2(context[CENTER], context[BOTTOM_RIGHT], context[BOTTOM])

	case 83:
		tl = interp1(context[CENTER], context[LEFT])
		if !equalYuv(yuvContext[TOP], yuvContext[RIGHT]) {
			tr = interp1(context[CENTER], context[TOP_RIGHT])
		} else {
			tr = interp7(context[CENTER], context[TOP], context[RIGHT])
		}
		bl = interp2(context[CENTER], context[BOTTOM_LEFT], context[LEFT])
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 92:
		tl = interp2(context[CENTER], context[TOP_LEFT], context[TOP])
		tr = interp1(context[CENTER], context[TOP])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		if !equalYuv(yuvContext[RIGHT], yuvContext[BOTTOM]) {
			br = interp1(context[CENTER], context[BOTTOM_RIGHT])
		} else {
			br = interp7(context[CENTER], context[RIGHT], context[BOTTOM])
		}

	case 202:
		if !equalYuv(yuvContext[LEFT], yuvContext[TOP]) {
			tl = interp1(context[CENTER], context[TOP_LEFT])
		} else {
			tl = interp7(context[CENTER], context[LEFT], context[TOP])
		}
		tr = interp2(context[CENTER], context[TOP_RIGHT], context[RIGHT])
		if !equalYuv(yuvContext[BOTTOM], yuvContext[LEFT]) {
			bl = interp1(context[CENTER], context[BOTTOM_LEFT])
		} else {
			bl = interp7(context[CENTER], context[BOTTOM], context[LEFT])
		}
		br = interp1(context[CENTER], context[RIGHT])

	case 78:
	case 154:
	case 114:
	case 89:
	case 90:
	case 55, 23:
	case 182, 150:
	case 213, 212:
	case 241, 240:
	case 236, 232:
	case 109, 105:
	case 171, 43:
	case 143, 15:
	case 124:
	case 203:
	case 62:
	case 211:
	case 118:
	case 217:
	case 110:
	case 155:
	case 188:
	case 185:
	case 61:
	case 157:
	case 103:
	case 227:
	case 230:
	case 199:
	case 220:
	case 158:
	case 234:
	case 242:
	case 59:
	case 121:
	case 87:
	case 79:
	case 122:
	case 94:
	case 218:
	case 91:
	case 229:
	case 167:
	case 173:
	case 181:
	case 186:
	case 115:
	case 93:
	case 206:
	case 205, 201:
	case 174, 46:
	case 179, 147:
	case 117, 116:
	case 189:
	case 231:
	case 126:
	case 219:
	case 125:
	case 221:
	case 207:
	case 238:
	case 190:
	case 187:
	case 243:
	case 119:
	case 237, 233:
	case 175, 47:
	case 183, 151:
	case 245, 244:
	case 250:
	case 123:
	case 95:
	case 222:
	case 252:
	case 249:
	case 235:
	case 111:
	case 63:
	case 159:
	case 215:
	case 246:
	case 254:
	case 253:
	case 251:
	case 239:
	case 127:
	case 191:
	case 223:
	case 247:
	case 255:
	default:
		panic(fmt.Errorf("invalid pattern: %d", pattern))
	}

	return tl, tr, bl, br
}

func equalYuv(a color.YCbCr, b color.YCbCr) bool {
	const (
		yThreshhold = 48.
		uThreshhold = 7.
		vThreshhold = 6.
	)

	aY, aU, aV := a.Y, a.Cb, a.Cr
	bY, bU, bV := b.Y, b.Cb, b.Cr

	if math.Abs(float64(aY)-float64(bY)) > yThreshhold {
		return false
	}
	if math.Abs(float64(aU)-float64(bU)) > uThreshhold {
		return false
	}
	if math.Abs(float64(aV)-float64(bV)) > vThreshhold {
		return false
	}

	return true
}

func newContextFlag() [9]uint8 {
	contextFlag := [9]uint8{}
	curFlag := uint8(1)

	for i := 0; i < 9; i++ {
		if i == CENTER {
			continue
		}

		contextFlag[i] = curFlag
		curFlag = curFlag << 1
	}

	return contextFlag
}

func RGBAToYCbCr(c color.RGBA) color.YCbCr {
	r, g, b := c.R, c.G, c.B
	y, u, v := color.RGBToYCbCr(r, g, b)
	return color.YCbCr{
		Y:  y,
		Cb: u,
		Cr: v,
	}
}
