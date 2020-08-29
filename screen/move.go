package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/character"
	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/movable"
	"github.com/bobtfish/mayhem/player"
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
	textBottom(fmt.Sprintf("%s's turn", screen.Players[screen.PlayerIdx].Name), ss, batch)
	batch.Draw(win)

	// 0 skips movement turn
	if win.JustPressed(pixelgl.Key0) {
		return NextPlayerMove(screen.PlayerIdx, screen.Players, screen.WithBoard)
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
	screen.WithBoard.CursorSprite = CursorBox
	if screen.MovedCharacters == nil {
		screen.MovedCharacters = make(map[movable.Movable]bool)
	}
	fmt.Printf("Enter move find character screen for player %d\n", screen.PlayerIdx+1)
}

func (screen *MoveFindCharacterScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	screen.WithBoard.DrawCursor(ss, batch)
	screen.WithBoard.MoveCursor(win)
	batch.Draw(win)

	if win.JustPressed(pixelgl.Key0) {
		return NextPlayerMove(screen.PlayerIdx, screen.Players, screen.WithBoard)
	}
	if win.JustPressed(pixelgl.KeyS) {
		// work out what's in this square, start moving it if movable and it belongs to the current player
		ob, isMovable := screen.WithBoard.Grid.GetGameObject(screen.WithBoard.CursorPosition).(movable.Movable)
		if isMovable {
			fmt.Printf("Is movable\n")
			if ob.CheckBelongsTo(screen.Players[screen.PlayerIdx]) {
				fmt.Printf("Belongs to this player\n")
				// Definitely something we can move

				// Skip if already moved
				if _, movedAlready := screen.MovedCharacters[ob]; movedAlready {
					fmt.Printf("Not movable (has moved already)\n")
					return screen
				}
				screen.MovedCharacters[ob] = true

				// Are we a mount with a wizard mounted
				char, isChar := ob.(*character.Character)
				if isChar && char.CarryingPlayer {
					return &MaybeDismount{
						WithBoard:       screen.WithBoard,
						PlayerIdx:       screen.PlayerIdx,
						Character:       ob,
						MovedCharacters: screen.MovedCharacters,
					}
				}

				// This check is after dismount for magic castle
				if ob.GetMovement() == 0 {
					fmt.Printf("Not movable (0 movement range)\n")
					return screen
				}

				// Is it engaged?
				if IsNextToEngageable(screen.WithBoard.CursorPosition, screen.PlayerIdx, screen.WithBoard) {
					fmt.Printf("Is next to engageable character\n")
					if !ob.BreakEngagement() {
						fmt.Printf("Did not break engagement, must do engaged attack\n")
						return &EngagedAttack{
							WithBoard:       screen.WithBoard,
							PlayerIdx:       screen.PlayerIdx,
							Character:       ob,
							MovedCharacters: screen.MovedCharacters,
						}
					}
					fmt.Printf("Broke engagement, can move normally\n")
				}

				// Not engaged, so move
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

type MaybeDismount struct {
	*WithBoard
	PlayerIdx       int
	MovedCharacters map[movable.Movable]bool
	Character       movable.Movable
}

func (screen *MaybeDismount) Enter(ss pixel.Picture, win *pixelgl.Window) {}

func (screen *MaybeDismount) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	isStaticCharacter := false
	if screen.Character.GetMovement() == 0 { // Magic castle / dark citadel
		isStaticCharacter = true
	} else {
		textBottom("Dismount wizard (Y or N)", ss, batch)
	}
	batch.Draw(win)

	if isStaticCharacter || win.JustPressed(pixelgl.KeyY) { // Do dismount
		// Dismount sets the character to not moved, the wizard moves
		delete(screen.MovedCharacters, screen.Character)
		wizard := screen.Character.(*character.Character).BelongsTo
		screen.MovedCharacters[wizard] = true
		if wizard.IsFlying() {
			return &MoveFlyingCharacterScreen{
				WithBoard:       screen.WithBoard,
				PlayerIdx:       screen.PlayerIdx,
				Character:       wizard,
				MovedCharacters: screen.MovedCharacters,
				IsDismount:      true,
			}
		}
		return &MoveGroundCharacterScreen{
			WithBoard:       screen.WithBoard,
			PlayerIdx:       screen.PlayerIdx,
			Character:       wizard,
			MovementLeft:    wizard.GetMovement(),
			MovedCharacters: screen.MovedCharacters,
			IsDismount:      true,
		}
	}

	if win.JustPressed(pixelgl.KeyN) { // Move character as normal
		// Note in the original game, characters with a player never seem to get engaged, so I deliberately skip that here
		if screen.Character.IsFlying() {
			return &MoveFlyingCharacterScreen{
				WithBoard:       screen.WithBoard,
				PlayerIdx:       screen.PlayerIdx,
				Character:       screen.Character,
				MovedCharacters: screen.MovedCharacters,
			}
		}
		return &MoveGroundCharacterScreen{
			WithBoard:       screen.WithBoard,
			PlayerIdx:       screen.PlayerIdx,
			Character:       screen.Character,
			MovementLeft:    screen.Character.GetMovement(),
			MovedCharacters: screen.MovedCharacters,
		}
	}
	return screen
}

func NextPlayerMove(playerIdx int, players []*player.Player, withBoard *WithBoard) GameScreen {
	nextIdx := NextPlayerIdx(playerIdx, players)
	if nextIdx == len(withBoard.Players) {
		return &GrowScreen{
			WithBoard: withBoard,
		}
	}
	return &MoveAnnounceScreen{
		WithBoard: withBoard,
		PlayerIdx: nextIdx,
	}
}

type MoveGroundCharacterScreen struct {
	*WithBoard
	PlayerIdx       int
	Character       movable.Movable
	MovementLeft    int
	NumDiagonals    int
	MovedCharacters map[movable.Movable]bool
	IsDismount      bool
}

func (screen *MoveGroundCharacterScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	fmt.Printf("Enter move ground character screen for player %d\n", screen.PlayerIdx+1)
}

func (screen *MoveGroundCharacterScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	textBottom(fmt.Sprintf("Movement range=%d", screen.MovementLeft), ss, batch)

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

		ms := MoveDoAttackMaybe(currentLocation, newLocation, screen.PlayerIdx, screen.WithBoard, screen.MovedCharacters, screen.IsDismount)
		if ms.NextScreen != nil {
			return ms.NextScreen
		}
		if ms.IllegalUndeadAttack {
			textBottom("Undead - Cannot be attacked", ss, batch)
		}
		if ms.DidMove {
			screen.WithBoard.CursorPosition = newLocation

			if ms.MountMove {
				screen.Character = screen.WithBoard.Grid.GetGameObject(newLocation).(movable.Movable)
				screen.MovedCharacters[screen.Character] = true
				return screen.MoveGroundCharacterScreenFinished()
			}

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
	}
	batch.Draw(win)
	if win.JustPressed(pixelgl.Key0) || win.JustPressed(pixelgl.KeyK) {
		return screen.MoveGroundCharacterScreenFinished()
	}
	return screen
}

// If moving to an empty square does the move then returns true allowing progress to continue (and ground/air specific logic to follow)
// If moving to a square which can be attacked, returns the attack screen to make the attack happen
// If moving to a square with something that cannot be moved into or attacked, return false
func MoveDoAttackMaybe(from, to logical.Vec, playerIdx int, withBoard *WithBoard, movedCharacters map[movable.Movable]bool, isDismount bool) MoveStatus {
	as := DoAttackMaybe(from, to, playerIdx, withBoard, movedCharacters, isDismount)
	if as.NotEmpty {
		fmt.Printf("Do attack maybe, was not empty\n")
		// FIXME AttackStatus.ToMoveStatus method?
		if as.NextScreen != nil {
			return MoveStatus{
				DidMove:    true,
				NextScreen: as.NextScreen,
			}
		}

		if as.IsMount {
			doMount(from, to, withBoard.Grid)
			return MoveStatus{
				DidMove:   true,
				MountMove: true,
			}
		}

		return MoveStatus{
			IllegalUndeadAttack: as.IllegalUndeadAttack,
			DidMove:             false,
			NextScreen:          nil,
		}
	}

	// Is an empty square, move to it
	doCharacterMove(from, to, withBoard.Grid, isDismount)

	// If you move next to an engageable character, you always become engaged in combat
	if IsNextToEngageable(to, playerIdx, withBoard) {
		fmt.Printf("Has moved next to engageable character, should be engaged\n")
		return MoveStatus{
			NextScreen: &EngagedAttack{
				WithBoard:       withBoard,
				PlayerIdx:       playerIdx,
				Character:       withBoard.Grid.GetGameObject(to).(movable.Movable),
				MovedCharacters: movedCharacters,
			},
			DidMove: true,
		}
	}
	return MoveStatus{
		DidMove:    true,
		NextScreen: nil,
	}
}

type AttackStatus struct {
	NotEmpty            bool
	IllegalUndeadAttack bool
	NextScreen          GameScreen
	IsMount             bool
}

func getAttacker(ob grid.GameObject, isDismount bool) movable.Attackerable {
	if !isDismount {
		return ob.(movable.Attackerable)
	}
	return ob.(*character.Character).BelongsTo
}

func DoAttackMaybe(from, to logical.Vec, playerIdx int, withBoard *WithBoard, movedCharacters map[movable.Movable]bool, isDismount bool) AttackStatus {
	target := withBoard.Grid.GetGameObject(to)
	if !target.IsEmpty() {
		fmt.Printf("Target square is not empty\n")
		ob, attackable := target.(movable.Attackable)
		if attackable {
			fmt.Printf("Target square is attackable\n")
			if !ob.CheckBelongsTo(withBoard.Players[playerIdx]) {
				attacker := getAttacker(withBoard.Grid.GetGameObject(from), isDismount)
				if ob.IsUndead() && !attacker.CanAttackUndead() {
					fmt.Printf("Cannot attack undead\n")
					return AttackStatus{
						NotEmpty:            true,
						IllegalUndeadAttack: true,
					}
				}
				fmt.Printf("Target square belongs to a different player do attack\n")
				fx := fx.Attack()
				withBoard.Grid.PlaceGameObject(to, fx)
				return AttackStatus{
					NotEmpty: true,
					NextScreen: &WaitForFx{
						NextScreen: &DoAttack{
							AttackerV:       from,
							DefenderV:       to,
							WithBoard:       withBoard,
							PlayerIdx:       playerIdx,
							MovedCharacters: movedCharacters,
							IsDismount:      isDismount,
						},
						Grid: withBoard.Grid,
						Fx:   fx,
					},
				}
			}
			// Does belong to this player, see if it's mountable, if so we can move to it
			_, isPlayer := withBoard.Grid.GetGameObject(from).(*player.Player)
			if isPlayer { // We're moving the player, lets see if target is something we can mount
				fmt.Printf("Moving player, check for mount\n")
				if ob.IsMount() {
					fmt.Printf("  Is mount\n")
					return AttackStatus{NotEmpty: true, IsMount: true}
				}
			}
		}
		return AttackStatus{NotEmpty: true}
	}
	return AttackStatus{}
}

func IsNextToEngageable(location logical.Vec, playerIdx int, withBoard *WithBoard) bool {
	for _, adjVec := range withBoard.Grid.AsRect().Adjacents(location) {
		adj := withBoard.Grid.GetGameObject(adjVec)
		if !adj.IsEmpty() {
			ob, attackable := adj.(movable.Attackable)
			if attackable {
				if !ob.CheckBelongsTo(withBoard.Players[playerIdx]) {
					if ob.Engageable() {
						return true
					}
				}
			}
		}
	}
	return false
}

// this gets reused from screen/attack.go
func doCharacterMove(from, to logical.Vec, grid *grid.GameGrid, isDismount bool) {
	stack := grid.GetGameObjectStack(from)
	if isDismount {
		mount := stack.TopObject().(*character.Character)
		mount.CarryingPlayer = false
		mount.BelongsTo.SetBoardPosition(to)
		grid.PlaceGameObject(to, mount.BelongsTo)
		return
	}
	character := stack.RemoveTopObject()
	character.SetBoardPosition(to)
	grid.PlaceGameObject(to, character)
}

func doMount(from, to logical.Vec, grid *grid.GameGrid) {
	fmt.Printf("doMount\n")
	// Take the player off the board (as we just track them as mounted on the character)
	// but set their position (in case they dismount without moving first)
	grid.GetGameObjectStack(from).RemoveTopObject().SetBoardPosition(to)
	grid.GetGameObject(to).(*character.Character).Mount()
}

func (screen *MoveGroundCharacterScreen) MoveGroundCharacterScreenFinished() GameScreen {
	return &RangedCombat{
		WithBoard:       screen.WithBoard,
		PlayerIdx:       screen.PlayerIdx,
		Character:       screen.Character,
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
	IsDismount      bool
}

func (screen *MoveFlyingCharacterScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	screen.WithBoard.CursorSprite = CursorFly
	fmt.Printf("Enter move flying character screen for player %d\n", screen.PlayerIdx+1)
	screen.DisplayRange = true // Set this to start to suppress cursor till we move it
}

type MoveStatus struct {
	DidMove             bool
	IllegalUndeadAttack bool
	NextScreen          GameScreen
	MountMove           bool
}

func (screen *MoveFlyingCharacterScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	if screen.DisplayRange {
		textBottom(fmt.Sprintf("Movement range=%d (flying)", screen.Character.GetMovement()), ss, batch)
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
			textBottom("Out of range", ss, batch)
			screen.OutOfRange = true
		} else {
			// work out what's in this square, if nothing move to it, if something attack it
			ms := MoveDoAttackMaybe(currentLocation, screen.WithBoard.CursorPosition, screen.PlayerIdx, screen.WithBoard, screen.MovedCharacters, screen.IsDismount)
			if ms.NextScreen != nil {
				return ms.NextScreen
			}
			if ms.IllegalUndeadAttack {
				screen.OutOfRange = true
				screen.DisplayRange = false
				textBottom("Undead - Cannot be attacked", ss, batch)
			}
			if ms.DidMove {
				fmt.Printf("Did do flying move, finish screen\n")

				if ms.MountMove {
					screen.Character = screen.WithBoard.Grid.GetGameObject(screen.WithBoard.CursorPosition).(movable.Movable)
					screen.MovedCharacters[screen.Character] = true
				}

				return screen.MoveFlyingCharacterScreenFinished()
			}
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
	return &RangedCombat{
		WithBoard:       screen.WithBoard,
		PlayerIdx:       screen.PlayerIdx,
		Character:       screen.Character,
		MovedCharacters: screen.MovedCharacters,
	}
}
