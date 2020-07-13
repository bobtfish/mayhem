package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type ExamineBoardScreen struct {
	MainMenu GameScreen
	*WithBoard
}

func (screen *ExamineBoardScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
}

func (screen *ExamineBoardScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	screen.WithBoard.MoveAndDrawCursor(ss, win, batch)
	batch.Draw(win)

	c := captureNumKey(win)
	if c == 0 {
		fmt.Println("Return to main menu")
		return screen.MainMenu
	}

	return screen
}
