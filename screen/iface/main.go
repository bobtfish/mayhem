package iface

import (
	"github.com/bobtfish/mayhem/grid"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type GameScreen interface {
	Enter(GameCtx)
	Step(GameCtx) GameScreen
}

type GameCtx interface {
	GetWindow() *pixelgl.Window
	GetSpriteSheet() pixel.Picture
	GetGrid() *grid.GameGrid
}
