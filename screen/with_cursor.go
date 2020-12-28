package screen

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"

	"github.com/bobtfish/mayhem/game"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

type WithCursor struct {
	CursorPosition  logical.Vec
	CursorShow      bool
	CursorFlashTime time.Time
	CursorSprite    int // Defaults to CursorSpell
}

func (screen *WithCursor) ShouldIDrawCursor() bool {
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

func (screen *WithCursor) DrawBoard(ctx screeniface.GameCtx) *pixel.Batch {
	ss := ctx.GetSpriteSheet()
	sd := render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())
	return ctx.GetGrid().DrawBatch(sd)
}

func (screen *WithCursor) MoveCursor(ctx screeniface.GameCtx) bool {
	win := ctx.GetWindow()
	grid := ctx.GetGrid()
	players := ctx.(*game.Window).GetPlayers()
	c := captureNumKey(win)
	if c > 0 && c <= len(players) {
		fmt.Printf("Flash player %d characters\n", c)
	}
	v := captureDirectionKey(win)
	if !v.Equals(logical.ZeroVec()) {
		fmt.Printf("Move cursor V(%d, %d)\n", v.X, v.Y)
		screen.CursorPosition = grid.AsRect().Clamp(screen.CursorPosition.Add(v))
		return true
	}
	return false
}

func (screen *WithCursor) DrawCursor(ctx screeniface.GameCtx, batch pixel.Target) {
	ss := ctx.GetSpriteSheet()
	grid := ctx.GetGrid()
	sd := render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())
	objectAtCursor := grid.GetGameObject(screen.CursorPosition)
	cursorColor := render.GetColor(0, 255, 255)
	if !objectAtCursor.IsEmpty() {
		cursorColor = objectAtCursor.GetColor()
	}
	description := objectAtCursor.Describe()
	textBottom(description, ss, batch)
	if screen.ShouldIDrawCursor() || objectAtCursor.IsEmpty() {
		sd.DrawSpriteColor(cursorSprite(screen.CursorSprite), screen.CursorPosition, cursorColor, batch)
	}
}

func (screen *WithCursor) MoveAndDrawCursor(ctx screeniface.GameCtx, batch pixel.Target) {
	screen.MoveCursor(ctx)
	screen.DrawCursor(ctx, batch)
}