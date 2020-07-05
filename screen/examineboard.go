package screen

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/render"
)

type ExamineBoardScreen struct {
	MainMenu        GameScreen
	Grid            *grid.GameGrid
	Players         []*player.Player
	CursorPosition  logical.Vec
	CursorFlash     bool
	CursorFlashTime time.Time
}

func (screen *ExamineBoardScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
}

func (screen *ExamineBoardScreen) ShouldIDrawCursor() bool {
	now := time.Now()
	if screen.CursorFlashTime.Before(now) {
		newFlash := true
		if screen.CursorFlash {
			newFlash = false
		}
		screen.CursorFlash = newFlash
		screen.CursorFlashTime = now.Add(time.Second / 8)
	}
	return screen.CursorFlash
}

func (screen *ExamineBoardScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	sd := render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())
	batch := screen.Grid.DrawBatch(sd)

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

	objectAtCursor := screen.Grid.GetGameObject(screen.CursorPosition)
	cursorColor := render.GetColor(0, 255, 255)
	if !objectAtCursor.IsEmpty() {
		cursorColor = objectAtCursor.GetColor()
	}
	render.NewTextDrawer(ss).DrawText(objectAtCursor.Describe(), logical.V(0, 0), win)
	if screen.ShouldIDrawCursor() || objectAtCursor.IsEmpty() {
		sd.DrawSpriteColor(cursorSprite(CURSOR_SPELL), screen.CursorPosition, cursorColor, batch)
	}

	batch.Draw(win)

	return screen
}
