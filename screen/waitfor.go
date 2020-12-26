package screen

import (
	"time"

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

func (screen *WaitFor) Enter(ctx screeniface.GameCtx) {}

func (screen *WaitFor) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
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

func (screen *Pause) Enter(ctx screeniface.GameCtx) {
	if screen.For == time.Duration(0) {
		screen.For = time.Second
	}
}

func (screen *Pause) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
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

func (screen *WaitForFx) Enter(ctx screeniface.GameCtx) {}

func (screen *WaitForFx) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	return &WaitFor{
		NextScreen: screen.NextScreen,
		Grid:       screen.Grid,
		FinishedF: func() bool {
			return screen.Fx == nil || screen.Fx.RemoveMe()
		},
	}
}
