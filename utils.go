package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"math/rand"

)

func pickColour() pixel.RGBA {
	return pixel.RGB(rand.Float64(), rand.Float64(), rand.Float64())
}

func drawHydra(ss pixel.Picture, win *pixelgl.Window) {
	rect := pixel.R(0, 16, 16, 32)
	sprite := pixel.NewSprite(ss, rect)
	mat := pixel.IM
	mat = mat.Moved(win.Bounds().Center())
	mat = mat.ScaledXY(win.Bounds().Center(), pixel.V(4, 4))
	sprite.Draw(win, mat)
}
/*
func placeCharactersTest(grid *grid.GameGrid, ct CharacterTypes) {
	x := 0
	y := 0
	for k := range ct {
		grid.PlaceGameObject(logical.V(x, y), ct.NewCharacter(k))
		x++
		if x == 15 {
			x = 0
			y++
		}
	}
}*/

/*
func blowSomethingUp(grid *grid.GameGrid) {
	x := rand.Intn(14)
	y := rand.Intn(9)
	fxA := []*Fx{FxWarp(), FxBlam(), FxFire(), FxBoom(), FxPop()}
	fxn := rand.Intn(len(fxA))
	grid.PlaceGameObject(logical.V(x, y), fxA[fxn])
}*/
