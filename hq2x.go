package hq2x

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

const (
	topLeft = iota
	top
	topRight
	left
	center
	right
	bottomLeft
	bottom
	bottomRight
)

// HQ2x - Enlarge image by 2x with hq2x algorithm
func HQ2x(src image.Image) (*image.RGBA, error) {
	srcX, srcY := src.Bounds().Dx(), src.Bounds().Dy()

	dest := image.NewRGBA(image.Rect(0, 0, srcX*2, srcY*2))

	for x := 0; x < srcX; x++ {
		for y := 0; y < srcY; y++ {
			context := [9]color.Color{
				getPixel(src, x-1, y-1), getPixel(src, x, y-1), getPixel(src, x+1, y-1),
				getPixel(src, x-1, y), getPixel(src, x, y), getPixel(src, x+1, y),
				getPixel(src, x-1, y+1), getPixel(src, x, y+1), getPixel(src, x+1, y+1),
			}

			tmp := hq2xPixel(context)
			tl, tr, bl, br := tmp[0], tmp[1], tmp[2], tmp[3]
			dest.Set(x*2, y*2, tl)
			dest.Set(x*2+1, y*2, tr)
			dest.Set(x*2, y*2+1, bl)
			dest.Set(x*2+1, y*2+1, br)
		}
	}

	return dest, nil
}

func getPixel(src image.Image, x, y int) color.Color {
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

	return src.At(x, y)
}

func hq2xPixel(context [9]color.Color) [4]color.Color {
	result := [4]color.Color{}

	yuvContext := [9]color.YCbCr{}
	yuvPixel := colorToYCbCr(context[center])
	for i := 0; i <= 9; i++ {
		yuvContext[i] = colorToYCbCr(context[i])
	}

	contextFlag := newContextFlag()
	var pattern uint8
	for bit := 0; bit <= 9; bit++ {
		if bit != center && !equalYuv(yuvContext[bit], yuvPixel) {
			pattern |= contextFlag[bit]
		}
	}

	switch pattern {
	case 0, 1, 4, 32, 128, 5, 132, 160, 33, 129, 36, 133, 164, 161, 37, 165:
	case 2, 34, 130, 162:
	case 16, 17, 48, 49:
	case 64, 65, 68, 69:
	case 8, 12, 136, 140:
	case 3, 35, 131, 163:
	case 6, 38, 134, 166:
	case 20, 21, 52, 53:
	case 144, 145, 176, 177:
	case 192, 193, 196, 197:
	case 96, 97, 100, 101:
	case 40, 44, 168, 172:
	case 9, 13, 137, 141:
	case 18, 50:
	case 80, 81:
	case 72, 76:
	case 10, 138:
	case 66:
	case 24:
	case 7, 39, 135:
	case 148, 149, 180:
	case 224, 228, 225:
	case 41, 169, 45:
	case 22, 54:
	case 208, 209:
	case 104, 108:
	case 11, 139:
	case 19, 51:
	case 146, 178:
	case 84, 85:
	case 112, 113:
	case 200, 204:
	case 73, 77:
	case 42, 170:
	case 14, 142:
	case 67:
	case 70:
	case 28:
	case 152:
	case 194:
	case 98:
	case 56:
	case 25:
	case 26, 31:
	case 82, 214:
	case 88, 248:
	case 74, 107:
	case 27:
	case 86:
	case 216:
	case 106:
	case 30:
	case 210:
	case 120:
	case 75:
	case 29:
	case 198:
	case 184:
	case 99:
	case 57:
	case 71:
	case 156:
	case 226:
	case 60:
	case 195:
	case 102:
	case 153:
	case 58:
	case 83:
	case 92:
	case 202:
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

	return result
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
	result := [9]uint8{}
	curFlag := uint8(1)

	for i := 0; i < 9; i++ {
		if i == center {
			continue
		}

		result[i] = curFlag
		curFlag = curFlag << 1
	}

	return result
}

func colorToYCbCr(c color.Color) color.YCbCr {
	r, g, b := colorToRGB(c)
	y, u, v := color.RGBToYCbCr(r, g, b)
	return color.YCbCr{
		Y:  y,
		Cb: u,
		Cr: v,
	}
}
