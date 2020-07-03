package screen

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/bobtfish/mayhem/logical"
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
	offset := logical.V(0, render.CHAR_PIXELS)
	sd := render.NewSpriteDrawer(ss, offset)
	drawMainBorder(win, sd)
	offset = offset.Add(logical.V(render.CHAR_PIXELS/2, render.CHAR_PIXELS/2))
	screen.SpriteDrawer = sd.WithOffset(offset)
	screen.TextDrawer = render.NewTextDrawer(ss, offset)
}
