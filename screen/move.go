package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/movable"
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
		// FIXME work out what's in this square, start moving it if movable and it belongs to the current player
		ob, isMovable := screen.WithBoard.Grid.GetGameObject(screen.WithBoard.CursorPosition).(movable.Movable)
		if isMovable {
			fmt.Printf("Is movable\n")
			if ob.CheckBelongsTo(screen.Players[screen.PlayerIdx]) {
				fmt.Printf("Belongs to this player\n")
				if ob.IsFlying() {
					return &MoveFlyingCharacterScreen{
						WithBoard: screen.WithBoard,
						PlayerIdx: screen.PlayerIdx,
						Character: ob,
					}
				}
				return &MoveGroundCharacterScreen{
					WithBoard:    screen.WithBoard,
					PlayerIdx:    screen.PlayerIdx,
					Character:    ob,
					MovementLeft: ob.GetMovement(),
				}
			}
		}
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

type MoveGroundCharacterScreen struct {
	*WithBoard
	PlayerIdx    int
	Character    movable.Movable
	MovementLeft int
	NumDiagonals int
}

func (screen *MoveGroundCharacterScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	fmt.Printf("Enter move ground character screen for player %d\n", screen.PlayerIdx+1)
}

func (screen *MoveGroundCharacterScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	render.NewTextDrawer(ss).DrawText(fmt.Sprintf("Movement range=%d", screen.MovementLeft), logical.V(0, 0), batch)
	batch.Draw(win)

	direction := captureDirectionKey(win)
	if direction != logical.ZeroVec() {
		fmt.Printf("Try to move in v(%d, %d)\n", direction.X, direction.Y)
		// work out if we can move to this square at all - if not ignore, if nothing move to it, if something attack it if owned by other player
		currentLocation := screen.Character.GetBoardPosition()
		newLocation := currentLocation.Add(direction)
		if !screen.WithBoard.Grid.AsRect().Contains(newLocation) {
			fmt.Printf("Cannot move out of screen\n")
			return screen
		}
		target := screen.WithBoard.Grid.GetGameObject(newLocation)
		// FIXME we can move into non-empty squares to attack
		if !target.IsEmpty() {
			fmt.Printf("Target square is not empty\n")
			return screen
		}
		screen.WithBoard.Grid.GetGameObjectStack(currentLocation).RemoveTopObject()
		// FIXME type cast here, puke
		screen.WithBoard.Grid.PlaceGameObject(newLocation, screen.Character.(grid.GameObjectStackable))
		screen.Character.SetBoardPosition(newLocation)
		screen.WithBoard.CursorPosition = newLocation

		// Do the D&D diagonal move thing
		if direction.IsDiagonal() {
			screen.NumDiagonals++
			if screen.NumDiagonals%2 == 0 {
				screen.MovementLeft--
			}
		}
		screen.MovementLeft--

		if screen.MovementLeft <= 0 {
			return &MoveFindCharacterScreen{
				WithBoard: screen.WithBoard,
				PlayerIdx: screen.PlayerIdx,
			}
		}
	}
	return screen
}

type MoveFlyingCharacterScreen struct {
	*WithBoard
	PlayerIdx  int
	Character  movable.Movable
	OutOfRange bool
}

func (screen *MoveFlyingCharacterScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	screen.WithBoard.CursorSprite = CURSOR_FLY
	fmt.Printf("Enter move flying character screen for player %d\n", screen.PlayerIdx+1)
}

func (screen *MoveFlyingCharacterScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	if screen.WithBoard.MoveCursor(win) || !screen.OutOfRange {
		screen.OutOfRange = false
		screen.WithBoard.DrawCursor(ss, batch)
	}

	if win.JustPressed(pixelgl.KeyS) {
		currentLocation := screen.Character.GetBoardPosition()
		if screen.WithBoard.CursorPosition.Distance(currentLocation) > screen.Character.GetMovement() {
			render.NewTextDrawer(ss).DrawText("Out of range                   ", logical.ZeroVec(), batch)
			screen.OutOfRange = true
		} else {
			// FIXME work out what's in this square, if nothing move to it, if something attack it
			target := screen.WithBoard.Grid.GetGameObject(screen.WithBoard.CursorPosition)
			// FIXME we can move into non-empty squares to attack
			if !target.IsEmpty() {
				fmt.Printf("Target square is not empty\n")
				return screen
			}
			screen.WithBoard.Grid.GetGameObjectStack(currentLocation).RemoveTopObject()
			// FIXME type cast here, puke
			screen.WithBoard.Grid.PlaceGameObject(screen.WithBoard.CursorPosition, screen.Character.(grid.GameObjectStackable))
			screen.Character.SetBoardPosition(screen.WithBoard.CursorPosition)

			return &MoveFindCharacterScreen{
				WithBoard: screen.WithBoard,
				PlayerIdx: screen.PlayerIdx,
			}
		}

	}
	batch.Draw(win)
	return screen
}
