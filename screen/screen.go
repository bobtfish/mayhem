package screen

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/bobtfish/mayhem/render"
)

func SpriteDrawer(ss pixel.Picture) render.SpriteDrawer {
	return render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())
}

func TextDrawer(ss pixel.Picture) render.SpriteDrawer {
	return render.NewTextDrawer(ss).WithOffset(render.GameBoardV())
}

func ClearScreen(ss pixel.Picture, win *pixelgl.Window) {
	win.Clear(colornames.Black)
	drawMainBorder(win, render.NewSpriteDrawer(ss).WithOffset(render.MainScreenV()))
}
