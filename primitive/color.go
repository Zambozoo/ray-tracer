package primitive

import (
	"fmt"
	"sync"
)

type ColorI struct {
	R uint8
	G uint8
	B uint8
}

func (c ColorI) String() string {
	return fmt.Sprintf("Color{R:%d,G:%d,B:%d}", c.R, c.G, c.B)
}

func (c *ColorI) ToColor() Color {
	return Color{
		R: float64(c.R) / 255,
		G: float64(c.G) / 255,
		B: float64(c.B) / 255,
	}
}

type Color struct {
	R float64
	G float64
	B float64
}

var (
	Black = Color{}
	White = Color{R: 1, G: 1, B: 1}
)

func (c Color) String() string {
	return fmt.Sprintf("Color{R:%f,G:%f,B:%f}", c.R, c.G, c.B)
}

func (c Color) Truncate() ColorI {
	return ColorI{
		R: uint8(c.R * 255),
		G: uint8(c.G * 255),
		B: uint8(c.B * 255),
	}
}

func clamp(f, low, high float64) float64 {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}

func (c Color) Clamp() Color {
	return Color{
		R: clamp(c.R, 0, 1),
		G: clamp(c.G, 0, 1),
		B: clamp(c.B, 0, 1),
	}
}

func (c Color) Scale(f float64) Color {
	return Color{
		R: c.R * f,
		G: c.G * f,
		B: c.B * f,
	}
}

func (c Color) Add(d Color) Color {
	return Color{
		R: c.R + d.R,
		G: c.G + d.G,
		B: c.B + d.B,
	}
}

func (c Color) Hadamard(d Color) Color {
	return Color{
		R: c.R * d.R,
		G: c.G * d.G,
		B: c.B * d.B,
	}
}

type AtomicPixel struct {
	Mu    sync.Mutex
	Pixel Color
}

func (ap *AtomicPixel) Add(c Color) {
	ap.Mu.Lock()
	defer ap.Mu.Unlock()
	ap.Pixel = ap.Pixel.Add(c)
}
