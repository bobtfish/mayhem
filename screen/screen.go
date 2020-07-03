package screen

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/bobtfish/mayhem/render"
)

type GameScreen interface {
	Setup(pixel.Picture, *pixelgl.Window)
	Draw(*pixelgl.Window)
	Finished() bool
	NextScreen() GameScreen
}

type ScreenBasics struct {
	SpriteSheet  pixel.Picture
	SpriteDrawer render.SpriteDrawer
	TextDrawer   render.SpriteDrawer
}

func (screen *ScreenBasics) Setup(ss pixel.Picture, win *pixelgl.Window) {
	win.Clear(colornames.Black)
	screen.SpriteSheet = ss
	drawMainBorder(win, render.NewSpriteDrawer(ss).WithOffset(render.MainScreenV()))
	screen.SpriteDrawer = render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())
	screen.TextDrawer = render.NewTextDrawer(ss).WithOffset(render.GameBoardV())
}
