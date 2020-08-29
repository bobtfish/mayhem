package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/character"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/movable"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/rand"
)

type RangedCombat struct {
	*WithBoard
	PlayerIdx       int
	Character       movable.Movable
	MovedCharacters map[movable.Movable]bool
	OutOfRange      bool
	DisplayRange    bool
}

func (screen *RangedCombat) Enter(ss pixel.Picture, win *pixelgl.Window) {
	fmt.Printf("In ranged combat state\n")
	screen.DisplayRange = true
	screen.WithBoard.CursorSprite = CursorRangedAttack
}

func (screen *RangedCombat) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	attacker := screen.Character.(movable.Attackerable)
	attackRange := attacker.GetAttackRange()
	if attackRange == 0 || win.JustPressed(pixelgl.Key0) || win.JustPressed(pixelgl.KeyK) { // No ranged combat
		return &MoveFindCharacterScreen{
			WithBoard:       screen.WithBoard,
			PlayerIdx:       screen.PlayerIdx,
			MovedCharacters: screen.MovedCharacters,
		}
	}

	batch := screen.WithBoard.DrawBoard(ss, win)

	// FIXME - this code is stolen from flying movement, can we consolidate?
	if screen.DisplayRange {
		textBottom(fmt.Sprintf("Ranged attack (range=%d)", attackRange), ss, batch)
	}
	cursorMoved := screen.WithBoard.MoveCursor(win)
	if cursorMoved || (!screen.OutOfRange && !screen.DisplayRange) {
		screen.OutOfRange = false
		screen.DisplayRange = false
		screen.WithBoard.DrawCursor(ss, batch)
	}

	if win.JustPressed(pixelgl.KeyS) {
		characterLocation := screen.Character.GetBoardPosition()
		attackDistance := screen.WithBoard.CursorPosition.Distance(characterLocation)
		if attackDistance > 0 { // You can't ranged attack yourself
			if attackDistance > attackRange {
				fmt.Printf("Out of range\n")
				textBottom("Out of range", ss, batch)
				screen.OutOfRange = true
			} else {
				if !HaveLineOfSight(characterLocation, screen.WithBoard.CursorPosition, screen.WithBoard.Grid) {
					textBottom("No line of sight", ss, batch)
					screen.OutOfRange = true
				} else {
					// Do ranged attack
					fx := attacker.GetAttackFx()
					screen.WithBoard.Grid.PlaceGameObject(screen.WithBoard.CursorPosition, fx)

					return &WaitForFx{
						Fx:   fx,
						Grid: screen.WithBoard.Grid,
						NextScreen: &DoRangedAttack{
							WithBoard:       screen.WithBoard,
							PlayerIdx:       screen.PlayerIdx,
							MovedCharacters: screen.MovedCharacters,
							Attacker:        attacker,
						},
					}
				}
			}
		}
	}
	batch.Draw(win)

	return screen
}

type DoRangedAttack struct {
	*WithBoard
	PlayerIdx       int
	MovedCharacters map[movable.Movable]bool
	Attacker        movable.Attackerable
}

func (screen *DoRangedAttack) Enter(ss pixel.Picture, win *pixelgl.Window) {
	fmt.Printf("Do ranged attack\n")
}

func (screen *DoRangedAttack) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	target := screen.WithBoard.Grid.GetGameObject(screen.WithBoard.CursorPosition)
	needPause := false
	if !target.IsEmpty() {
		fmt.Printf("Target square is not empty\n")
		ob, attackable := target.(movable.Attackable)
		if attackable {
			fmt.Printf("Target square is attackable\n")
			if !ob.IsUndead() || screen.Attacker.CanAttackUndead() {
				// FIXME this is duplicate logic
				defenceRating := ob.GetDefence() + rand.Intn(9)
				attackRating := screen.Attacker.GetRangedCombat() + rand.Intn(9)

				fmt.Printf("Attack rating %d defence rating %d\n", attackRating, defenceRating)
				if attackRating > defenceRating {
					fmt.Printf("Attack kills defender\n")
					_, newScreen := PostSuccessfulAttack(target, screen.WithBoard, false)
					if newScreen != nil {
						return newScreen
					}
				}
			} else {
				textBottom("Undead - Cannot be attacked", ss, win)
				needPause = true
			}
		}
	}

	return &Pause{
		Skip: !needPause,
		Grid: screen.WithBoard.Grid,
		NextScreen: &MoveFindCharacterScreen{
			WithBoard:       screen.WithBoard,
			PlayerIdx:       screen.PlayerIdx,
			MovedCharacters: screen.MovedCharacters,
		},
	}
}

type EngagedAttack struct {
	*WithBoard
	PlayerIdx       int
	Character       movable.Movable
	MovedCharacters map[movable.Movable]bool
	ClearMsg        bool
}

func (screen *EngagedAttack) Enter(ss pixel.Picture, win *pixelgl.Window) {
	fmt.Printf("In engaged attack state\n")
	textBottom("Engaged to enemy", ss, win)
}

func (screen *EngagedAttack) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	if screen.ClearMsg {
		textBottom("", ss, batch)
	}

	direction := captureDirectionKey(win)
	if direction != logical.ZeroVec() {
		screen.ClearMsg = true
		fmt.Printf("Try to move in v(%d, %d)\n", direction.X, direction.Y)
		currentLocation := screen.Character.GetBoardPosition()
		newLocation := currentLocation.Add(direction)

		as := DoAttackMaybe(currentLocation, newLocation, screen.PlayerIdx, screen.WithBoard, screen.MovedCharacters, false)
		if as.NextScreen != nil {
			fmt.Printf("Can attack in that direction")
			return as.NextScreen
		}
		if as.IllegalUndeadAttack {
			screen.ClearMsg = false
			textBottom("Undead - Cannot be attacked", ss, batch)
		}
	}

	if win.JustPressed(pixelgl.Key0) || win.JustPressed(pixelgl.KeyK) {
		return &RangedCombat{
			WithBoard:       screen.WithBoard,
			PlayerIdx:       screen.PlayerIdx,
			Character:       screen.Character,
			MovedCharacters: screen.MovedCharacters,
		}
	}
	batch.Draw(win)
	return screen
}

type DoAttack struct {
	*WithBoard
	AttackerV       logical.Vec
	DefenderV       logical.Vec
	PlayerIdx       int
	MovedCharacters map[movable.Movable]bool
	IsDismount      bool
}

func (screen *DoAttack) Enter(ss pixel.Picture, win *pixelgl.Window) {}

func PostSuccessfulAttack(target grid.GameObject, withBoard *WithBoard, canMakeCorpse bool) (bool, GameScreen) {
	canMoveOnto := true
	// If the defender can be killed, kill them. Otherwise remove them
	ob, corpsable := target.(movable.Corpseable)
	fmt.Printf("Defender is %T corpsable %v ob %T(%v)\n", target, corpsable, ob, ob)

	// Store character and player for later
	var p *player.Player
	char, isCharacter := target.(*character.Character)
	if isCharacter {
		p = char.BelongsTo
	}

	makesCorpse := canMakeCorpse && corpsable && ob.CanMakeCorpse()
	if makesCorpse {
		fmt.Printf("make corpse\n")
		ob.MakeCorpse()
	} else {
		fmt.Printf("remove defender as no corpse\n")
		died := withBoard.Grid.GetGameObjectStack(target.(movable.Movable).GetBoardPosition()).RemoveTopObject()
		if KillIfPlayer(died, withBoard.Grid) {
			if WeHaveAWinner(withBoard.Players) {
				return canMoveOnto, &WinnerScreen{
					WithBoard: withBoard,
				}
			}
		}
	}

	// If the thing that was just killed was carrying the player, put the player back on the board
	if isCharacter && char.CarryingPlayer {
		canMoveOnto = false
		fmt.Printf("Was carrying player, put back: %T %v\n", p, p)
		withBoard.Grid.PlaceGameObject(char.BoardPosition, p)
	}

	return canMoveOnto, nil
}

func (screen *DoAttack) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	// Work out what happened. This is overly simple, but equivalent to what the original game does :)
	defender := screen.WithBoard.Grid.GetGameObject(screen.DefenderV)
	defenceRating := defender.(movable.Attackable).GetDefence() + rand.Intn(9)
	attacker := getAttacker(screen.WithBoard.Grid.GetGameObject(screen.AttackerV), screen.IsDismount)
	attackRating := attacker.(movable.Attackerable).GetCombat() + rand.Intn(9)

	fmt.Printf("Attack rating %d defence rating %d\n", attackRating, defenceRating)
	if attackRating > defenceRating {
		canMoveOnto, newScreen := PostSuccessfulAttack(defender, screen.WithBoard, true)
		if newScreen != nil {
			return newScreen
		}

		if canMoveOnto {
			doCharacterMove(screen.AttackerV, screen.DefenderV, screen.WithBoard.Grid, screen.IsDismount)
			screen.AttackerV = screen.DefenderV
		}
	}
	return &RangedCombat{
		WithBoard:       screen.WithBoard,
		PlayerIdx:       screen.PlayerIdx,
		Character:       screen.WithBoard.Grid.GetGameObject(screen.AttackerV).(movable.Movable),
		MovedCharacters: screen.MovedCharacters,
	}
}
