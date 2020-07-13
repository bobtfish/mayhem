package screen

import (
	"fmt"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/render"
)

type WithBoard struct {
	CursorPosition  logical.Vec
	CursorShow      bool
	CursorFlashTime time.Time
	Grid            *grid.GameGrid
	Players         []*player.Player
}

func (screen *WithBoard) ShouldIDrawCursor() bool {
	now := time.Now()
	if screen.CursorFlashTime.Before(now) {
		newFlash := true
		if screen.CursorShow {
			newFlash = false
		}
		screen.CursorShow = newFlash
		screen.CursorFlashTime = now.Add(time.Second / 8)
	}
	return screen.CursorShow
}

func (screen *WithBoard) DrawBoard(ss pixel.Picture, win *pixelgl.Window) *pixel.Batch {
	sd := render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())
	return screen.Grid.DrawBatch(sd)
}

func (screen *WithBoard) MoveCursor(win *pixelgl.Window) bool {
	c := captureNumKey(win)
	if c > 0 && c <= len(screen.Players) {
		fmt.Printf("Flash player %d characters\n", c)
	}
	v := captureDirectionKey(win)
	if !v.Equals(logical.ZeroVec()) {
		fmt.Printf("Move cursor V(%d, %d)\n", v.X, v.Y)
		screen.CursorPosition = screen.Grid.AsRect().Clamp(screen.CursorPosition.Add(v))
		return true
	}
	return false
}

func (screen *WithBoard) DrawCursor(ss pixel.Picture, batch *pixel.Batch) {
	sd := render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())
	objectAtCursor := screen.Grid.GetGameObject(screen.CursorPosition)
	cursorColor := render.GetColor(0, 255, 255)
	if !objectAtCursor.IsEmpty() {
		cursorColor = objectAtCursor.GetColor()
	}
	description := objectAtCursor.Describe()
	td := render.NewTextDrawer(ss)
	td.DrawText(description, logical.V(0, 0), batch)
	td.DrawText(strings.Repeat(" ", 32-len(description)), logical.V(len(description), 0), batch)
	if screen.ShouldIDrawCursor() || objectAtCursor.IsEmpty() {
		sd.DrawSpriteColor(cursorSprite(CURSOR_SPELL), screen.CursorPosition, cursorColor, batch)
	}
}

func (screen *WithBoard) MoveAndDrawCursor(ss pixel.Picture, win *pixelgl.Window, batch *pixel.Batch) {
	screen.MoveCursor(win)
	screen.DrawCursor(ss, batch)
}
