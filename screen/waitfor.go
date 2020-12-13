package screen

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/render"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

type WaitFor struct {
	NextScreen screeniface.GameScreen
	Grid       *grid.GameGrid
	FinishedF  func() bool
}

func (screen *WaitFor) Enter(ss pixel.Picture, win *pixelgl.Window) {}

func (screen *WaitFor) Step(ss pixel.Picture, win *pixelgl.Window) screeniface.GameScreen {
	if screen.FinishedF() {
		//fmt.Printf("Waitfor Skip to next screen\n")
		return screen.NextScreen
	}
	if screen.Grid != nil {
		screen.Grid.DrawBatch(render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())).Draw(win)
	}
	return screen
}

type Pause struct {
	NextScreen screeniface.GameScreen
	Grid       *grid.GameGrid
	Skip       bool
	For        time.Duration
}

func (screen *Pause) Enter(ss pixel.Picture, win *pixelgl.Window) {
	if screen.For == time.Duration(0) {
		screen.For = time.Second
	}
}

func (screen *Pause) Step(ss pixel.Picture, win *pixelgl.Window) screeniface.GameScreen {
	started := time.Now()
	return &WaitFor{
		NextScreen: screen.NextScreen,
		Grid:       screen.Grid,
		FinishedF: func() bool {
			return screen.Skip || started.Add(screen.For).Before(time.Now())
		},
	}
}

type WaitForFx struct {
	Fx         *fx.Fx
	NextScreen screeniface.GameScreen
	Grid       *grid.GameGrid
}

func (screen *WaitForFx) Enter(ss pixel.Picture, win *pixelgl.Window) {}

func (screen *WaitForFx) Step(ss pixel.Picture, win *pixelgl.Window) screeniface.GameScreen {
	return &WaitFor{
		NextScreen: screen.NextScreen,
		Grid:       screen.Grid,
		FinishedF: func() bool {
			return screen.Fx == nil || screen.Fx.RemoveMe()
		},
	}
}
