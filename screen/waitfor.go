package screen

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/render"
)

type WaitFor struct {
	NextScreen GameScreen
	Grid       *grid.GameGrid
	Skip       bool
}

func (screen *WaitFor) Enter(ss pixel.Picture, win *pixelgl.Window) {}

func (screen *WaitFor) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	if screen.Skip {
		return screen.NextScreen
	}
	if screen.Grid != nil {
		screen.Grid.DrawBatch(render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())).Draw(win)
	}
	return screen
}

type Pause struct {
	*WaitFor
	Started time.Time
}

func (screen *Pause) Enter(ss pixel.Picture, win *pixelgl.Window) {
	screen.Started = time.Now()
}

func (screen *Pause) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	if screen.Started.Add(time.Second).Before(time.Now()) {
		screen.WaitFor.Skip = true
	}
	return screen.WaitFor.Step(ss, win)
}

type WaitForFx struct {
	Fx *fx.Fx
	*WaitFor
}

func (screen *WaitForFx) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	if screen.Fx.RemoveMe() {
		screen.WaitFor.Skip = true
	}
	return screen.WaitFor.Step(ss, win)
}
