package screen

import (
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

func captureNumKey(win *pixelgl.Window) int {
	if win.JustPressed(pixelgl.Key0) {
		return 0
	}
	if win.JustPressed(pixelgl.Key1) {
		return 1
	}
	if win.JustPressed(pixelgl.Key2) {
		return 2
	}
	if win.JustPressed(pixelgl.Key3) {
		return 3
	}
	if win.JustPressed(pixelgl.Key4) {
		return 4
	}
	if win.JustPressed(pixelgl.Key5) {
		return 5
	}
	if win.JustPressed(pixelgl.Key6) {
		return 6
	}
	if win.JustPressed(pixelgl.Key7) {
		return 7
	}
	if win.JustPressed(pixelgl.Key8) {
		return 8
	}
	if win.JustPressed(pixelgl.Key9) {
		return 9
	}
	return -1
}

func drawMainBorder(win *pixelgl.Window, sd *render.SpriteDrawer) {
	batch := sd.GetNewBatch()
	// Bottom left
	sd.DrawSprite(logical.V(6, 20), logical.V(0, 0), batch)
	// Bottom row
	for i := 1; i < 15; i++ {
		sd.DrawSprite(logical.V(7, 20), logical.V(i, 0), batch)
	}
	// Bottom right
	sd.DrawSprite(logical.V(8, 20), logical.V(15, 0), batch)
	// LHS and RHS
	for i := 1; i < 10; i++ {
		sd.DrawSprite(logical.V(5, 20), logical.V(0, i), batch)
		sd.DrawSprite(logical.V(9, 20), logical.V(15, i), batch)
	}
	// Top Left
	sd.DrawSprite(logical.V(2, 20), logical.V(0, 10), batch)
	// Top row
	for i := 1; i < 15; i++ {
		sd.DrawSprite(logical.V(3, 20), logical.V(i, 10), batch)
	}
	// Top right
	sd.DrawSprite(logical.V(4, 20), logical.V(15, 10), batch)
	batch.Draw(win)
}
