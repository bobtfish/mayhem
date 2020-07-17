package screen

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/render"
)

type Pause struct {
	NextScreen GameScreen
	Started    time.Time
	Grid       *grid.GameGrid
	Skip       bool
}

func (screen *Pause) Enter(ss pixel.Picture, win *pixelgl.Window) {
	screen.Started = time.Now()
}

func (screen *Pause) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	if screen.Skip {
		return screen.NextScreen
	}
	if screen.Grid != nil {
		screen.Grid.DrawBatch(render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())).Draw(win)
	}
	if screen.Started.Add(time.Second).Before(time.Now()) {
		return screen.NextScreen
	}
	return screen
}
