package screen

import (
	"fmt"

	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/character"
	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/game"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/movable"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/render"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

type MoveAnnounceScreen struct {
	PlayerIdx int
}

func (screen *MoveAnnounceScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
}

func (screen *MoveAnnounceScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	players := ctx.(*game.Window).GetPlayers()
	batch := DrawBoard(ctx)
	textBottomColor(fmt.Sprintf("%s's turn", players[screen.PlayerIdx].Name), render.GetColor(241, 241, 0), ss, batch)
	batch.Draw(win)

	// 0 skips movement turn
	if win.JustPressed(pixelgl.Key0) {
		return NextPlayerMove(screen.PlayerIdx, players, ctx)
	}

	// any other key displays the cursor
	if win.JustPressed(pixelgl.KeyS) || captureDirectionKey(win) != logical.ZeroVec() {
		return &MoveFindCharacterScreen{
			WithCursor: WithCursor{CursorPosition: players[screen.PlayerIdx].BoardPosition},
			PlayerIdx:  screen.PlayerIdx,
		}
	}
	return screen
}

type MoveFindCharacterScreen struct {
	WithCursor
	PlayerIdx       int
	MovedCharacters map[movable.Movable]bool
}

func (screen *MoveFindCharacterScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	screen.WithCursor.CursorSprite = CursorBox
	if screen.MovedCharacters == nil {
		screen.MovedCharacters = make(map[movable.Movable]bool)
	}
	fmt.Printf("Enter move find character screen for player %d\n", screen.PlayerIdx+1)
}

func (screen *MoveFindCharacterScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	grid := ctx.GetGrid()
	players := ctx.(*game.Window).GetPlayers()
	batch := DrawBoard(ctx)
	screen.WithCursor.DrawCursor(ctx, batch)
	screen.WithCursor.MoveCursor(ctx)
	batch.Draw(win)

	if win.JustPressed(pixelgl.Key0) {
		return NextPlayerMove(screen.PlayerIdx, players, ctx)
	}
	if win.JustPressed(pixelgl.KeyS) {
		// work out what's in this square, start moving it if movable and it belongs to the current player
		ob, isMovable := grid.GetGameObject(screen.WithCursor.CursorPosition).(movable.Movable)
		if isMovable {
			fmt.Printf("Is movable\n")
			if ob.CheckBelongsTo(players[screen.PlayerIdx]) {
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
						PlayerIdx:       screen.PlayerIdx,
						Character:       ob,
						MovedCharacters: screen.MovedCharacters,
					}
				}

				// Is it engaged?
				if IsNextToEngageable(screen.WithCursor.CursorPosition, screen.PlayerIdx, ctx) {
					fmt.Printf("Is next to engageable character\n")
					if !ob.BreakEngagement() {
						fmt.Printf("Did not break engagement, must do engaged attack\n")
						return &EngagedAttack{
							PlayerIdx:       screen.PlayerIdx,
							Character:       ob,
							MovedCharacters: screen.MovedCharacters,
						}
					}
					fmt.Printf("Broke engagement, can move normally\n")
				}

				// This check is after dismount for magic castle
				// and after engaged attack checking for shadow wood
				if ob.GetMovement() == 0 {
					fmt.Printf("Not movable (0 movement range)\n")
					return screen
				}

				// Not engaged, so move
				if ob.IsFlying() {
					return &MoveFlyingCharacterScreen{
						PlayerIdx:       screen.PlayerIdx,
						Character:       ob,
						MovedCharacters: screen.MovedCharacters,
					}
				}
				return &MoveGroundCharacterScreen{
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
	PlayerIdx       int
	MovedCharacters map[movable.Movable]bool
	Character       movable.Movable
}

func (screen *MaybeDismount) Enter(ctx screeniface.GameCtx) {}

func (screen *MaybeDismount) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	batch := DrawBoard(ctx)
	isStaticCharacter := false
	if screen.Character.GetMovement() == 0 { // Magic castle / dark citadel
		isStaticCharacter = true
	} else {
		textBottomColor("Dismount wizard (Y or N)", render.GetColor(244, 244, 0), ss, batch)
	}
	batch.Draw(win)

	if isStaticCharacter || win.JustPressed(pixelgl.KeyY) { // Do dismount
		// Dismount sets the character to not moved, the wizard moves
		delete(screen.MovedCharacters, screen.Character)
		wizard := screen.Character.(*character.Character).BelongsTo
		screen.MovedCharacters[wizard] = true
		if wizard.IsFlying() {
			return &MoveFlyingCharacterScreen{
				PlayerIdx:       screen.PlayerIdx,
				Character:       wizard,
				MovedCharacters: screen.MovedCharacters,
				IsDismount:      true,
			}
		}
		return &MoveGroundCharacterScreen{
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
				PlayerIdx:       screen.PlayerIdx,
				Character:       screen.Character,
				MovedCharacters: screen.MovedCharacters,
			}
		}
		return &MoveGroundCharacterScreen{
			PlayerIdx:       screen.PlayerIdx,
			Character:       screen.Character,
			MovementLeft:    screen.Character.GetMovement(),
			MovedCharacters: screen.MovedCharacters,
		}
	}
	return screen
}

func NextPlayerMove(playerIdx int, players []*player.Player, ctx screeniface.GameCtx) screeniface.GameScreen {
	nextIdx := NextPlayerIdx(playerIdx, players)
	if nextIdx == len(players) {
		return &GrowScreen{}
	}
	return &MoveAnnounceScreen{
		PlayerIdx: nextIdx,
	}
}

type MoveGroundCharacterScreen struct {
	PlayerIdx       int
	Character       movable.Movable
	MovementLeft    int
	NumDiagonals    int
	MovedCharacters map[movable.Movable]bool
	IsDismount      bool
}

func (screen *MoveGroundCharacterScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	fmt.Printf("Enter move ground character screen for player %d\n", screen.PlayerIdx+1)
}

func (screen *MoveGroundCharacterScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	grid := ctx.GetGrid()
	batch := DrawBoard(ctx)
	textBottomMulti([]TextWithColor{
		{Text: "Movement range=", Color: render.GetColor(0, 247, 0)},
		{Text: fmt.Sprintf("%d", screen.MovementLeft), Color: render.GetColor(237, 237, 0)},
	}, ss, batch)

	currentLocation := screen.Character.GetBoardPosition()

	direction := captureDirectionKey(win)
	if direction != logical.ZeroVec() {
		fmt.Printf("Try to move in v(%d, %d)\n", direction.X, direction.Y)
		// work out if we can move to this square at all - if not ignore, if nothing move to it, if something attack it if owned by other player
		newLocation := currentLocation.Add(direction)
		if !grid.AsRect().Contains(newLocation) {
			fmt.Printf("Cannot move out of screen\n")
			return screen
		}

		ms := MoveDoAttackMaybe(currentLocation, newLocation, screen.PlayerIdx, ctx, screen.MovedCharacters, screen.IsDismount)
		if ms.NextScreen != nil {
			return ms.NextScreen
		}
		if ms.IllegalUndeadAttack {
			textBottomColor("Undead - Cannot be attacked", render.GetColor(0, 244, 244), ss, win)
		}
		if ms.DidMove {
			if ms.MountMove {
				screen.Character = grid.GetGameObject(newLocation).(movable.Movable)
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
func MoveDoAttackMaybe(from, to logical.Vec, playerIdx int, ctx screeniface.GameCtx, movedCharacters map[movable.Movable]bool, isDismount bool) MoveStatus {
	as := DoAttackMaybe(from, to, playerIdx, ctx, movedCharacters, isDismount)
	grid := ctx.GetGrid()
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
			doMount(from, to, grid, isDismount)
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
	doCharacterMove(from, to, grid, isDismount)

	// If you move next to an engageable character, you always become engaged in combat
	if IsNextToEngageable(to, playerIdx, ctx) {
		fmt.Printf("Has moved next to engageable character, should be engaged\n")
		return MoveStatus{
			NextScreen: &EngagedAttack{
				PlayerIdx:       playerIdx,
				Character:       grid.GetGameObject(to).(movable.Movable),
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
	NextScreen          screeniface.GameScreen
	IsMount             bool
}

func getAttacker(ob grid.GameObject, isDismount bool) movable.Attackerable {
	if !isDismount {
		return ob.(movable.Attackerable)
	}
	return ob.(*character.Character).BelongsTo
}

func DoAttackMaybe(from, to logical.Vec, playerIdx int, ctx screeniface.GameCtx, movedCharacters map[movable.Movable]bool, isDismount bool) AttackStatus {
	grid := ctx.GetGrid()
	players := ctx.(*game.Window).GetPlayers()
	target := grid.GetGameObject(to)
	if !target.IsEmpty() {
		fmt.Printf("Target square is not empty\n")
		ob, attackable := target.(movable.Attackable)
		if attackable && ob.GetCombat() > 0 {
			fmt.Printf("Target square is attackable\n")
			if !ob.CheckBelongsTo(players[playerIdx]) {
				attacker := getAttacker(grid.GetGameObject(from), isDismount)
				if ob.IsUndead() && !attacker.CanAttackUndead() {
					fmt.Printf("Cannot attack undead\n")
					return AttackStatus{
						NotEmpty:            true,
						IllegalUndeadAttack: true,
					}
				}
				fmt.Printf("Target square belongs to a different player do attack\n")
				fx := fx.Attack()
				grid.PlaceGameObject(to, fx)
				return AttackStatus{
					NotEmpty: true,
					NextScreen: &WaitForFx{
						NextScreen: &DoAttack{
							AttackerV:       from,
							DefenderV:       to,
							PlayerIdx:       playerIdx,
							MovedCharacters: movedCharacters,
							IsDismount:      isDismount,
						},
						Fx: fx,
					},
				}
			}
		}
		// Does belong to this player, see if it's mountable, if so we can move to it
		_, isPlayer := grid.GetGameObject(from).(*player.Player)
		if isPlayer || isDismount { // We're moving the player, lets see if target is something we can mount
			fmt.Printf("Moving player, check for mount\n")
			if ob.IsMount() {
				fmt.Printf("  Is mount\n")
				return AttackStatus{NotEmpty: true, IsMount: true}
			}
		}
		return AttackStatus{NotEmpty: true}
	}
	return AttackStatus{}
}

func IsNextToEngageable(location logical.Vec, playerIdx int, ctx screeniface.GameCtx) bool {
	grid := ctx.GetGrid()
	players := ctx.(*game.Window).GetPlayers()
	for _, adjVec := range grid.AsRect().Adjacents(location) {
		adj := grid.GetGameObject(adjVec)
		if !adj.IsEmpty() {
			ob, attackable := adj.(movable.Attackable)
			if attackable {
				if !ob.CheckBelongsTo(players[playerIdx]) {
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

func doMount(from, to logical.Vec, grid *grid.GameGrid, isDismount bool) {
	fmt.Printf("doMount\n")
	stack := grid.GetGameObjectStack(from)
	if !isDismount {
		// Take the player off the board (as we just track them as mounted on the character)
		// but set their position (in case they dismount without moving first)
		stack.RemoveTopObject().SetBoardPosition(to)
	} else {
		// Was already mounted, need to dismount old place
		mount := stack.TopObject().(*character.Character)
		mount.CarryingPlayer = false
	}
	grid.GetGameObject(to).(*character.Character).Mount()
}

func (screen *MoveGroundCharacterScreen) MoveGroundCharacterScreenFinished() screeniface.GameScreen {
	return &RangedCombat{
		PlayerIdx:       screen.PlayerIdx,
		Character:       screen.Character,
		MovedCharacters: screen.MovedCharacters,
	}
}

type MoveFlyingCharacterScreen struct {
	WithCursor
	PlayerIdx       int
	Character       movable.Movable
	OutOfRange      bool
	DisplayRange    bool
	MovedCharacters map[movable.Movable]bool
	IsDismount      bool
}

func (screen *MoveFlyingCharacterScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	screen.WithCursor.CursorSprite = CursorFly
	screen.WithCursor.CursorPosition = screen.Character.GetBoardPosition()
	fmt.Printf("Enter move flying character screen for player %d\n", screen.PlayerIdx+1)
	screen.DisplayRange = true // Set this to start to suppress cursor till we move it
}

type MoveStatus struct {
	DidMove             bool
	IllegalUndeadAttack bool
	NextScreen          screeniface.GameScreen
	MountMove           bool
}

func (screen *MoveFlyingCharacterScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	grid := ctx.GetGrid()
	batch := DrawBoard(ctx)
	if screen.DisplayRange {
		textBottomMulti([]TextWithColor{
			{Text: "Movement range=", Color: render.GetColor(0, 247, 0)},
			{Text: fmt.Sprintf("%d", screen.Character.GetMovement()), Color: render.GetColor(237, 237, 0)},
			{Text: " (flying)", Color: render.GetColor(0, 245, 245)},
		}, ss, batch)
	}
	cursorMoved := screen.WithCursor.MoveCursor(ctx)
	if cursorMoved || (!screen.OutOfRange && !screen.DisplayRange) {
		screen.OutOfRange = false
		screen.DisplayRange = false
		screen.WithCursor.DrawCursor(ctx, batch)
	}

	if win.JustPressed(pixelgl.KeyS) {
		fmt.Printf("Try flying move\n")
		currentLocation := screen.Character.GetBoardPosition()
		target := screen.WithCursor.CursorPosition
		if target.Distance(currentLocation) > screen.Character.GetMovement() {
			fmt.Printf("Out of range\n")
			textBottomColor("Out of range", render.GetColor(0, 249, 249), ss, batch)
			screen.OutOfRange = true
		} else {
			// work out what's in this square, if nothing move to it, if something attack it
			ms := MoveDoAttackMaybe(currentLocation, target, screen.PlayerIdx, ctx, screen.MovedCharacters, screen.IsDismount)
			if ms.NextScreen != nil {
				return ms.NextScreen
			}
			if ms.IllegalUndeadAttack {
				screen.OutOfRange = true
				screen.DisplayRange = false
				textBottomColor("Undead - Cannot be attacked", render.GetColor(0, 244, 244), ss, win)
			}
			if ms.DidMove {
				fmt.Printf("Did do flying move, finish screen\n")

				if ms.MountMove {
					screen.Character = grid.GetGameObject(target).(movable.Movable)
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

func (screen *MoveFlyingCharacterScreen) MoveFlyingCharacterScreenFinished() screeniface.GameScreen {
	return &RangedCombat{
		PlayerIdx:       screen.PlayerIdx,
		Character:       screen.Character,
		MovedCharacters: screen.MovedCharacters,
	}
}
