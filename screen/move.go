package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

type MoveAnnounceScreen struct {
	*WithBoard
	PlayerIdx int
}

func (screen *MoveAnnounceScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	screen.WithBoard.CursorPosition = screen.Players[screen.PlayerIdx].BoardPosition
}

func (screen *MoveAnnounceScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	render.NewTextDrawer(ss).DrawText(fmt.Sprintf("%s's turn", screen.Players[screen.PlayerIdx].Name), logical.V(0, 0), batch)
	batch.Draw(win)

	if win.JustPressed(pixelgl.KeyS) || captureDirectionKey(win) != logical.ZeroVec() {
		return &MoveFindCharacterScreen{
			WithBoard: screen.WithBoard,
			PlayerIdx: screen.PlayerIdx,
		}
	}
	return screen
}

type MoveFindCharacterScreen struct {
	*WithBoard
	PlayerIdx int
}

func (screen *MoveFindCharacterScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	screen.WithBoard.CursorSprite = CURSOR_BOX
	fmt.Printf("Enter move find character screen for player %d\n", screen.PlayerIdx+1)
}

func (screen *MoveFindCharacterScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	screen.WithBoard.DrawCursor(ss, batch)
	screen.WithBoard.MoveCursor(win)
	batch.Draw(win)

	if win.JustPressed(pixelgl.Key0) {
		return screen.NextMove()
	}
	if win.JustPressed(pixelgl.KeyS) {
		// FIXME work out what's in this square, start moving it if movable
	}

	return screen
}

func (screen *MoveFindCharacterScreen) NextMove() GameScreen {
	if screen.PlayerIdx+1 == len(screen.WithBoard.Players) {
		return &GrowScreen{
			WithBoard: screen.WithBoard,
		}
	}
	return &MoveAnnounceScreen{
		WithBoard: screen.WithBoard,
		PlayerIdx: screen.PlayerIdx + 1,
	}
}
