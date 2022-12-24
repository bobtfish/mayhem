package render

import (
	"image/color"

	"github.com/faiface/pixel"
)

func GetColor(r, g, b int) color.Color {
	if r == 0 && g == 0 && b == 0 {
		return pixel.RGB(1, 1, 1)
	}
	return pixel.RGB(float64(r)/255, float64(g)/255, float64(b)/255)
}

func ColorWhite() color.Color {
	return GetColor(255, 255, 255)
}

func ColorRed() color.Color {
	return GetColor(255, 0, 0)
}

func ColorYellow() color.Color {
	return GetColor(255, 255, 0)
}

func ColorGreen() color.Color {
	return GetColor(0, 255, 0)
}

func ColorBlue() color.Color {
	return GetColor(0, 0, 255)
}

func ColorPurple() color.Color {
	return GetColor(255, 0, 255)
}

func ColorCyan() color.Color {
	return GetColor(1, 255, 255)
}
