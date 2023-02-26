package vsv

import "image/color"

// === BGR ===

type BGR struct {
	B, G, R uint8
}

func (c BGR) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(0xFFFF)
	return
}

func bgrModel(c color.Color) color.Color {
	if _, ok := c.(BGR); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	return BGR{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}

var BGRColorModel color.Model = color.ModelFunc(bgrModel)


// === BGRA ===

type BGRA struct {
	B, G, R, A uint8
}

func (c BGRA) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(c.A)
	a |= a << 8
	return
}

func bgraModel(c color.Color) color.Color {
	if _, ok := c.(BGRA); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return BGRA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

var BGRAColorModel color.Model = color.ModelFunc(bgraModel)
