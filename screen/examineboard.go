package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/render"
)

type ExamineBoardScreen struct {
	MainMenu       GameScreen
	Grid           *grid.GameGrid
	Players        []*player.Player
	CursorPosition logical.Vec
}

func (screen *ExamineBoardScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
}

func (screen *ExamineBoardScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	sd := render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())
	batch := screen.Grid.DrawBatch(sd)

	fmt.Printf("Draw Cursor at V(%d, %d)\n", screen.CursorPosition.X, screen.CursorPosition.Y)
	sd.DrawSprite(cursorSprite(), screen.CursorPosition, batch)

	batch.Draw(win)

	c := captureNumKey(win)
	if c == 0 {
		fmt.Println("Return to main menu")
		return screen.MainMenu
	}
	if c > 0 && c <= len(screen.Players) {
		fmt.Printf("Flash player %d characters\n", c)
	}
	v := captureDirectionKey(win)
	if !v.Equals(logical.ZeroVec()) {
		fmt.Printf("Move cursor V(%d, %d)\n", v.X, v.Y)
		screen.CursorPosition = screen.Grid.AsRect().Clamp(screen.CursorPosition.Add(v))
	}
	return screen
}
