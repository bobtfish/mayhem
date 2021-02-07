package screen

import (
	"fmt"

	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/character"
	"github.com/bobtfish/mayhem/game"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/movable"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/rand"
	"github.com/bobtfish/mayhem/render"

	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

type RangedCombat struct {
	WithCursor
	PlayerIdx       int
	Character       movable.Movable
	MovedCharacters map[movable.Movable]bool
	OutOfRange      bool
	DisplayRange    bool
}

func (screen *RangedCombat) Enter(ctx screeniface.GameCtx) {
	screen.DisplayRange = true
	screen.WithCursor.CursorSprite = CursorRangedAttack
	screen.WithCursor.CursorPosition = screen.Character.GetBoardPosition()
}

func (screen *RangedCombat) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	grid := ctx.GetGrid()
	attacker := screen.Character.(movable.Attackerable)
	attackRange := attacker.GetAttackRange()
	if attackRange == 0 || win.JustPressed(pixelgl.Key0) || win.JustPressed(pixelgl.KeyK) { // No ranged combat
		return &MoveFindCharacterScreen{
			WithCursor:      WithCursor{CursorPosition: screen.WithCursor.CursorPosition},
			PlayerIdx:       screen.PlayerIdx,
			MovedCharacters: screen.MovedCharacters,
		}
	}

	batch := DrawBoard(ctx)

	// FIXME - this code is stolen from flying movement, can we consolidate?
	if screen.DisplayRange {
		textBottomMulti([]TextWithColor{
			{Text: "Ranged attack, range=", Color: render.GetColor(0, 246, 0)},
			{Text: fmt.Sprintf("%d", attackRange), Color: render.GetColor(241, 241, 0)},
		}, ss, batch)
	}
	cursorMoved := screen.WithCursor.MoveCursor(ctx)
	if cursorMoved || (!screen.OutOfRange && !screen.DisplayRange) {
		screen.OutOfRange = false
		screen.DisplayRange = false
		screen.WithCursor.DrawCursor(ctx, batch)
	}

	if win.JustPressed(pixelgl.KeyS) {
		attackPosition := screen.WithCursor.CursorPosition
		characterLocation := screen.Character.GetBoardPosition()
		attackDistance := attackPosition.Distance(characterLocation)
		if attackDistance > 0 { // You can't ranged attack yourself
			if attackDistance > attackRange {
				textBottomColor("Out of range", render.GetColor(247, 247, 0), ss, batch)
				screen.OutOfRange = true
			} else {
				if !grid.HaveLineOfSight(characterLocation, attackPosition) {
					textBottomColor("No line of sight", render.GetColor(0, 233, 233), ss, batch)
					screen.OutOfRange = true
				} else {
					return &AnimateRangedAttack{
						AttackPosition:  attackPosition,
						PlayerIdx:       screen.PlayerIdx,
						MovedCharacters: screen.MovedCharacters,
						Attacker:        attacker,
					}
				}
			}
		}
	}
	batch.Draw(win)

	return screen
}

type AnimateRangedAttack struct {
	AttackPosition  logical.Vec
	PlayerIdx       int
	MovedCharacters map[movable.Movable]bool
	Attacker        movable.Attackerable
	StepIdx         int
	AnimationSteps  []logical.Vec
}

func (screen *AnimateRangedAttack) Enter(ctx screeniface.GameCtx) {
	from := screen.Attacker.GetBoardPosition()
	to := screen.AttackPosition
	fmt.Printf("Attack from %d, %d TO %d, %d\n", from.X, from.Y, to.X, to.Y)
	screen.AnimationSteps = to.Subtract(from).Path()
	for i, s := range screen.AnimationSteps {
		screen.AnimationSteps[i] = from.Add(s)
	}
	screen.AnimationSteps = append(screen.AnimationSteps, to)
}

func (screen *AnimateRangedAttack) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	if screen.StepIdx == len(screen.AnimationSteps) {
		// Do ranged attack
		fx := screen.Attacker.GetAttackFx()
		grid := ctx.GetGrid()
		grid.PlaceGameObject(screen.AttackPosition, fx)

		return &WaitForFx{
			Fx: fx,
			NextScreen: &DoRangedAttack{
				AttackPosition:  screen.AttackPosition,
				PlayerIdx:       screen.PlayerIdx,
				MovedCharacters: screen.MovedCharacters,
				Attacker:        screen.Attacker,
			},
		}
	}
	step := screen.AnimationSteps[screen.StepIdx]
	fmt.Printf("Animation step %d, %d\n", step.X, step.Y)
	screen.StepIdx++
	return screen
}

type DoRangedAttack struct {
	AttackPosition  logical.Vec
	PlayerIdx       int
	MovedCharacters map[movable.Movable]bool
	Attacker        movable.Attackerable
}

func (screen *DoRangedAttack) Enter(ctx screeniface.GameCtx) {
	fmt.Printf("Do ranged attack\n")
}

func (screen *DoRangedAttack) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	target := ctx.GetGrid().GetGameObject(screen.AttackPosition)
	needPause := false
	if !target.IsEmpty() {
		fmt.Printf("Target square is not empty\n")
		ob, attackable := target.(movable.Attackable)
		if attackable && ob.GetCombat() > 0 {
			fmt.Printf("Target square is attackable\n")
			if !ob.IsUndead() || screen.Attacker.CanAttackUndead() {
				// FIXME this is duplicate logic
				defenceRating := ob.GetDefence() + rand.Intn(9)
				attackRating := screen.Attacker.GetRangedCombat() + rand.Intn(9)

				fmt.Printf("Attack rating %d defence rating %d\n", attackRating, defenceRating)
				if attackRating > defenceRating {
					fmt.Printf("Attack kills defender\n")
					_, newScreen := PostSuccessfulAttack(target, ctx, false)
					if newScreen != nil {
						return newScreen
					}
				}
			} else {
				textBottomColor("Undead - Cannot be attacked", render.GetColor(0, 244, 244), ss, win)
				needPause = true
			}
		}
	}

	return &Pause{
		Skip: !needPause,
		NextScreen: &MoveFindCharacterScreen{
			WithCursor:      WithCursor{CursorPosition: screen.AttackPosition},
			PlayerIdx:       screen.PlayerIdx,
			MovedCharacters: screen.MovedCharacters,
		},
	}
}

type EngagedAttack struct {
	PlayerIdx       int
	Character       movable.Movable
	MovedCharacters map[movable.Movable]bool
	ClearMsg        bool
}

func (screen *EngagedAttack) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	fmt.Printf("In engaged attack state\n")
	textBottomColor("Engaged to enemy", render.GetColor(247, 247, 0), ss, win)
}

func (screen *EngagedAttack) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	batch := DrawBoard(ctx)
	if screen.ClearMsg {
		textBottom("", ss, batch)
	}

	direction := captureDirectionKey(win)
	if direction != logical.ZeroVec() {
		screen.ClearMsg = true
		fmt.Printf("Try to move in v(%d, %d)\n", direction.X, direction.Y)
		currentLocation := screen.Character.GetBoardPosition()
		newLocation := currentLocation.Add(direction)

		as := DoAttackMaybe(currentLocation, newLocation, screen.PlayerIdx, ctx, screen.MovedCharacters, false)
		if as.NextScreen != nil {
			fmt.Printf("Can attack in that direction")
			return as.NextScreen
		}
		if as.IllegalUndeadAttack {
			screen.ClearMsg = false
			textBottomColor("Undead - Cannot be attacked", render.GetColor(0, 244, 244), ss, win)
		}
	}

	if win.JustPressed(pixelgl.Key0) || win.JustPressed(pixelgl.KeyK) {
		return &RangedCombat{
			PlayerIdx:       screen.PlayerIdx,
			Character:       screen.Character,
			MovedCharacters: screen.MovedCharacters,
		}
	}
	batch.Draw(win)
	return screen
}

type DoAttack struct {
	AttackerV       logical.Vec
	DefenderV       logical.Vec
	PlayerIdx       int
	MovedCharacters map[movable.Movable]bool
	IsDismount      bool
}

func (screen *DoAttack) Enter(ctx screeniface.GameCtx) {}

func PostSuccessfulAttack(target grid.GameObject, ctx screeniface.GameCtx, canMakeCorpse bool) (bool, screeniface.GameScreen) {
	grid := ctx.GetGrid()
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
		died := grid.GetGameObjectStack(target.(movable.Movable).GetBoardPosition()).RemoveTopObject()
		if KillIfPlayer(died, grid) {
			players := ctx.(*game.Window).GetPlayers()
			if WeHaveAWinner(players) {
				return canMoveOnto, &WinnerScreen{
					Players: players,
				}
			}
		}
	}

	// If the thing that was just killed was carrying the player, put the player back on the board
	if isCharacter && char.CarryingPlayer {
		canMoveOnto = false
		fmt.Printf("Was carrying player, put back: %T %v\n", p, p)
		grid.PlaceGameObject(char.BoardPosition, p)
	}

	return canMoveOnto, nil
}

func (screen *DoAttack) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	grid := ctx.GetGrid()
	// Work out what happened. This is overly simple, but equivalent to what the original game does :)
	defender := grid.GetGameObject(screen.DefenderV)
	defenceRating := defender.(movable.Attackable).GetDefence() + rand.Intn(9)
	attacker := getAttacker(grid.GetGameObject(screen.AttackerV), screen.IsDismount)
	attackRating := attacker.(movable.Attackerable).GetCombat() + rand.Intn(9)

	fmt.Printf("Attack rating %d defence rating %d\n", attackRating, defenceRating)
	if attackRating > defenceRating {
		canMoveOnto, newScreen := PostSuccessfulAttack(defender, ctx, true)
		if newScreen != nil {
			return newScreen
		}

		if canMoveOnto {
			doCharacterMove(screen.AttackerV, screen.DefenderV, grid, screen.IsDismount)
			screen.AttackerV = screen.DefenderV
		}
	}
	return &RangedCombat{
		PlayerIdx:       screen.PlayerIdx,
		Character:       grid.GetGameObject(screen.AttackerV).(movable.Movable),
		MovedCharacters: screen.MovedCharacters,
	}
}
