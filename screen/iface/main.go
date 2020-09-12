package iface

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type GameScreen interface {
	Enter(pixel.Picture, *pixelgl.Window)
	Step(pixel.Picture, *pixelgl.Window) GameScreen
}
