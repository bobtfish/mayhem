package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type MoveScreen struct {
	*WithBoard
	PlayerIdx int
}

func (screen *MoveScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	fmt.Printf("Enter move screen for player %d\n", screen.PlayerIdx+1)
}

func (screen *MoveScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	batch.Draw(win)

	if win.JustPressed(pixelgl.Key0) {
		return screen.NextMove()
	}

	return screen
}

func (screen *MoveScreen) NextMove() GameScreen {
	if screen.PlayerIdx == len(screen.WithBoard.Players) {
		return &GrowScreen{
			WithBoard: screen.WithBoard,
		}
	}
	return &MoveScreen{
		WithBoard: screen.WithBoard,
		PlayerIdx: screen.PlayerIdx + 1,
	}
}
