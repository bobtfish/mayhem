package screen

import (
	"time"

	"github.com/bobtfish/mayhem/fx"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

type WaitFor struct {
	NextScreen screeniface.GameScreen
	FinishedF  func() bool
}

func (screen *WaitFor) Enter(ctx screeniface.GameCtx) {}

func (screen *WaitFor) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	if screen.FinishedF() {
		//fmt.Printf("Waitfor Skip to next screen\n")
		return screen.NextScreen
	}
	DrawBoard(ctx).Draw(ctx.GetWindow())
	return screen
}

type Pause struct {
	NextScreen screeniface.GameScreen
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
		FinishedF: func() bool {
			return screen.Skip || started.Add(screen.For).Before(time.Now())
		},
	}
}

type WaitForFx struct {
	Fx         *fx.Fx
	NextScreen screeniface.GameScreen
}

func (screen *WaitForFx) Enter(ctx screeniface.GameCtx) {}

func (screen *WaitForFx) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	return &WaitFor{
		NextScreen: screen.NextScreen,
		FinishedF: func() bool {
			return screen.Fx == nil || screen.Fx.RemoveMe()
		},
	}
}
