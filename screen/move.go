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

	// 0 skips movement turn
	if win.JustPressed(pixelgl.Key0) {
		return NextPlayerMove(screen.PlayerIdx, screen.WithBoard)
	}

	// any other key displays the cursor
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
	PlayerIdx       int
	MovedCharacters map[movable.Movable]bool
}

func (screen *MoveFindCharacterScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	screen.WithBoard.CursorSprite = CURSOR_BOX
	if screen.MovedCharacters == nil {
		screen.MovedCharacters = make(map[movable.Movable]bool, 0)
	}
	fmt.Printf("Enter move find character screen for player %d\n", screen.PlayerIdx+1)
}

func (screen *MoveFindCharacterScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	screen.WithBoard.DrawCursor(ss, batch)
	screen.WithBoard.MoveCursor(win)
	batch.Draw(win)

	if win.JustPressed(pixelgl.Key0) {
		return NextPlayerMove(screen.PlayerIdx, screen.WithBoard)
	}
	if win.JustPressed(pixelgl.KeyS) {
		// work out what's in this square, start moving it if movable and it belongs to the current player
		ob, isMovable := screen.WithBoard.Grid.GetGameObject(screen.WithBoard.CursorPosition).(movable.Movable)
		if isMovable {
			fmt.Printf("Is movable\n")
			if ob.CheckBelongsTo(screen.Players[screen.PlayerIdx]) {
				fmt.Printf("Belongs to this player\n")
				if ob.GetMovement() == 0 {
					fmt.Printf("Not movable (0 movement range)\n")
					return screen
				}

				if _, movedAlready := screen.MovedCharacters[ob]; movedAlready {
					fmt.Printf("Not movable (has moved already)\n")
					return screen
				} else {
					screen.MovedCharacters[ob] = true
				}

				// Definitely something we can move
				if ob.IsFlying() {
					return &MoveFlyingCharacterScreen{
						WithBoard:       screen.WithBoard,
						PlayerIdx:       screen.PlayerIdx,
						Character:       ob,
						MovedCharacters: screen.MovedCharacters,
					}
				}
				return &MoveGroundCharacterScreen{
					WithBoard:       screen.WithBoard,
					PlayerIdx:       screen.PlayerIdx,
					Character:       ob,
					MovementLeft:    ob.GetMovement(),
					MovedCharacters: screen.MovedCharacters,
				}
			}
		}
	}

	return screen
}

func NextPlayerMove(playerIdx int, withBoard *WithBoard) GameScreen {
	if playerIdx+1 == len(withBoard.Players) {
		return &GrowScreen{
			WithBoard: withBoard,
		}
	}
	return &MoveAnnounceScreen{
		WithBoard: withBoard,
		PlayerIdx: playerIdx + 1,
	}
}

type MoveGroundCharacterScreen struct {
	*WithBoard
	PlayerIdx       int
	Character       movable.Movable
	MovementLeft    int
	NumDiagonals    int
	MovedCharacters map[movable.Movable]bool
}

func (screen *MoveGroundCharacterScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	fmt.Printf("Enter move ground character screen for player %d\n", screen.PlayerIdx+1)
}

func (screen *MoveGroundCharacterScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	render.NewTextDrawer(ss).DrawText(fmt.Sprintf("Movement range=%d", screen.MovementLeft), logical.V(0, 0), batch)
	batch.Draw(win)

	currentLocation := screen.Character.GetBoardPosition()

	direction := captureDirectionKey(win)
	if direction != logical.ZeroVec() {
		fmt.Printf("Try to move in v(%d, %d)\n", direction.X, direction.Y)
		// work out if we can move to this square at all - if not ignore, if nothing move to it, if something attack it if owned by other player
		newLocation := currentLocation.Add(direction)
		if !screen.WithBoard.Grid.AsRect().Contains(newLocation) {
			fmt.Printf("Cannot move out of screen\n")
			return screen
		}

		newScreen, didMove := MoveDoAttackMaybe(currentLocation, newLocation, screen.PlayerIdx, screen.WithBoard, screen.MovedCharacters)
		if newScreen != nil {
			return newScreen
		}
		if !didMove {
			return screen
		}

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
			return screen.MoveGroundCharacterScreenFinished()
		}
	}
	if win.JustPressed(pixelgl.Key0) || win.JustPressed(pixelgl.KeyK) {
		return screen.MoveGroundCharacterScreenFinished()
	}
	return screen
}

// If moving to an empty square does the move then returns true allowing progress to continue (and ground/air specific logic to follow)
// If moving to a square which can be attacked, returns the attack screen to make the attack happen
// If moving to a square with something that cannot be moved into or attacked, return false
func MoveDoAttackMaybe(from, to logical.Vec, playerIdx int, withBoard *WithBoard, movedCharacters map[movable.Movable]bool) (GameScreen, bool) {
	newScreen, notEmpty := DoAttackMaybe(from, to, playerIdx, withBoard, movedCharacters)
	if notEmpty {
		if newScreen != nil {
			return newScreen, true
		}
		return nil, false
	}

	// Is an empty square, move to it
	doCharacterMove(from, to, withBoard.Grid)

	// If you move next to an engageable character, you always become engaged in combat
	if IsNextToEngageable(to, playerIdx, withBoard) {
		fmt.Printf("Has moved next to engageable character, should be engaged\n")
		return &EngagedAttack{
			WithBoard:       withBoard,
			PlayerIdx:       playerIdx,
			Character:       withBoard.Grid.GetGameObject(to).(movable.Movable),
			MovedCharacters: movedCharacters,
		}, true
	}

	return nil, true
}

func DoAttackMaybe(from, to logical.Vec, playerIdx int, withBoard *WithBoard, movedCharacters map[movable.Movable]bool) (GameScreen, bool) {
	target := withBoard.Grid.GetGameObject(to)
	if !target.IsEmpty() {
		fmt.Printf("Target square is not empty\n")
		ob, attackable := target.(movable.Attackable)
		if attackable {
			fmt.Printf("Target square is attackable\n")
			if !ob.CheckBelongsTo(withBoard.Players[playerIdx]) {
				fmt.Printf("Target square belongs to a different player do attack\n")
				return &DoAttack{
					AttackerV:       from,
					DefenderV:       to,
					WithBoard:       withBoard,
					PlayerIdx:       playerIdx,
					MovedCharacters: movedCharacters,
				}, true
			}
		}
		return nil, true
	}
	return nil, false
}

func IsNextToEngageable(location logical.Vec, playerIdx int, withBoard *WithBoard) bool {
	for _, adjVec := range withBoard.Grid.AsRect().Adjacents(location) {
		adj := withBoard.Grid.GetGameObject(adjVec)
		ob, attackable := adj.(movable.Attackable)
		if attackable {
			if !ob.CheckBelongsTo(withBoard.Players[playerIdx]) {
				if ob.Engageable() {
					return true
				}
			}
		}
	}
	return false
}

// this gets reused from screen/attack.go
func doCharacterMove(from, to logical.Vec, grid *grid.GameGrid) {
	character := grid.GetGameObjectStack(from).RemoveTopObject()
	character.SetBoardPosition(to)
	grid.PlaceGameObject(to, character)
}

func (screen *MoveGroundCharacterScreen) MoveGroundCharacterScreenFinished() GameScreen {
	return &MoveFindCharacterScreen{
		WithBoard:       screen.WithBoard,
		PlayerIdx:       screen.PlayerIdx,
		MovedCharacters: screen.MovedCharacters,
	}
}

type MoveFlyingCharacterScreen struct {
	*WithBoard
	PlayerIdx       int
	Character       movable.Movable
	OutOfRange      bool
	DisplayRange    bool
	MovedCharacters map[movable.Movable]bool
}

func (screen *MoveFlyingCharacterScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	screen.WithBoard.CursorSprite = CURSOR_FLY
	fmt.Printf("Enter move flying character screen for player %d\n", screen.PlayerIdx+1)
	screen.DisplayRange = true // Set this to start to supress cursor till we move it
}

func (screen *MoveFlyingCharacterScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	if screen.DisplayRange {
		render.NewTextDrawer(ss).DrawText(fmt.Sprintf("Movement range=%d (flying)", screen.Character.GetMovement()), logical.ZeroVec(), batch)
	}
	cursorMoved := screen.WithBoard.MoveCursor(win)
	if cursorMoved || (!screen.OutOfRange && !screen.DisplayRange) {
		screen.OutOfRange = false
		screen.DisplayRange = false
		screen.WithBoard.DrawCursor(ss, batch)
	}

	if win.JustPressed(pixelgl.KeyS) {
		fmt.Printf("Try flying move\n")
		currentLocation := screen.Character.GetBoardPosition()
		if screen.WithBoard.CursorPosition.Distance(currentLocation) > screen.Character.GetMovement() {
			fmt.Printf("Out of range\n")
			render.NewTextDrawer(ss).DrawText("Out of range                   ", logical.ZeroVec(), batch)
			screen.OutOfRange = true
		} else {
			// work out what's in this square, if nothing move to it, if something attack it
			newScreen, didMove := MoveDoAttackMaybe(currentLocation, screen.WithBoard.CursorPosition, screen.PlayerIdx, screen.WithBoard, screen.MovedCharacters)
			if newScreen != nil {
				return newScreen
			}
			if !didMove {
				return screen
			}

			return screen.MoveFlyingCharacterScreenFinished()
		}

	}

	// Allow you to cancel movement, but that character then counts as moved
	if win.JustPressed(pixelgl.Key0) || win.JustPressed(pixelgl.KeyK) {
		return screen.MoveFlyingCharacterScreenFinished()
	}

	batch.Draw(win)
	return screen
}

func (screen *MoveFlyingCharacterScreen) MoveFlyingCharacterScreenFinished() GameScreen {
	return &MoveFindCharacterScreen{
		WithBoard:       screen.WithBoard,
		PlayerIdx:       screen.PlayerIdx,
		MovedCharacters: screen.MovedCharacters,
	}
}
