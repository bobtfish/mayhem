package screen

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type GrowScreen struct {
	*WithBoard
}

func (screen *GrowScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
}

func (screen *GrowScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	batch.Draw(win)
	return &Pause{
		Grid: screen.WithBoard.Grid,
		NextScreen: &TurnMenuScreen{
			Players: screen.WithBoard.Players,
			Grid:    screen.WithBoard.Grid,
		},
	}
}
