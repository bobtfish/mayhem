package screen

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/bobtfish/mayhem/render"
)

type GameScreen interface {
	Enter(pixel.Picture, *pixelgl.Window)
	Step(*pixelgl.Window) GameScreen
}

type ScreenBasics struct {
	SpriteSheet pixel.Picture
}

func SpriteDrawer(ss pixel.Picture) render.SpriteDrawer {
	return render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())
}

func TextDrawer(ss pixel.Picture) render.SpriteDrawer {
	return render.NewTextDrawer(ss).WithOffset(render.GameBoardV())
}

func (screen *ScreenBasics) Enter(ss pixel.Picture, win *pixelgl.Window) {
	win.Clear(colornames.Black)
	screen.SpriteSheet = ss
	drawMainBorder(win, render.NewSpriteDrawer(ss).WithOffset(render.MainScreenV()))
}
