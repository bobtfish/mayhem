package screen

import (
	"fmt"

	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

type ExamineBoardScreen struct {
	MainMenu screeniface.GameScreen
	WithCursor
}

func (screen *ExamineBoardScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
}

func (screen *ExamineBoardScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	batch := DrawBoard(ctx)
	screen.WithCursor.MoveAndDrawCursor(ctx, batch)
	batch.Draw(win)

	c := captureNumKey(win)
	if c == 0 {
		fmt.Println("Return to main menu")
		return screen.MainMenu
	}

	return screen
}
