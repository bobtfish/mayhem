package screen

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/render"
)

type CastSpellScreen struct {
	Grid *grid.GameGrid
}

func (screen *CastSpellScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
}

func (screen *CastSpellScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	sd := render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())
	batch := screen.Grid.DrawBatch(sd)
	batch.Draw(win)
	return screen
}
