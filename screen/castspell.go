package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/logical"
)

// Move state onto next player spell cast (if there are players left)
// or onto the movement phase if all spells have been cast
func NextSpellCastOrMove(playerIdx int, wb *WithBoard, skipPause bool) GameScreen {
	var nextScreen GameScreen
	nextIdx := NextPlayerIdx(playerIdx, wb.Players)
	nextScreen = &DisplaySpellCastScreen{
		WithBoard: wb,
		PlayerIdx: nextIdx,
	}

	if nextIdx == len(wb.Players) {
		// All players have cast their spells, movement comes next
		nextScreen = &MoveAnnounceScreen{
			WithBoard: wb,
		}
	}
	return &Pause{
		Skip:       skipPause,
		Grid:       wb.Grid,
		NextScreen: nextScreen,
	}
}

// Display the spell name in the bottom bar until player presses a direction key or s

type DisplaySpellCastScreen struct {
	*WithBoard
	PlayerIdx int
}

func (screen *DisplaySpellCastScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	thisPlayer := screen.Players[screen.PlayerIdx]
	screen.WithBoard.CursorPosition = thisPlayer.BoardPosition
	if thisPlayer.ChosenSpell >= 0 {
		spell := thisPlayer.Spells[thisPlayer.ChosenSpell]
		batch := screen.WithBoard.DrawBoard(ss, win)
		textBottom(fmt.Sprintf("%s %s %d", thisPlayer.Name, spell.GetName(), spell.GetCastRange()), ss, batch)
		batch.Draw(win)
	}
}

func (screen *DisplaySpellCastScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	thisPlayer := screen.Players[screen.PlayerIdx]
	if (thisPlayer.ChosenSpell < 0) || win.JustPressed(pixelgl.Key0) {
		return NextSpellCastOrMove(screen.PlayerIdx, screen.WithBoard, true)
	}
	if win.JustPressed(pixelgl.KeyS) || !captureDirectionKey(win).Equals(logical.ZeroVec()) {
		return &TargetSpellScreen{
			WithBoard:      screen.WithBoard,
			PlayerIdx:      screen.PlayerIdx,
			CastsRemaining: thisPlayer.Spells[thisPlayer.ChosenSpell].CastQuantity(),
		}
	}
	return screen
}

// If range 0 then cast the spell straight away (on the wizard
// If range > 0 then move cursor around to find a target until S is pressed
type TargetSpellScreen struct {
	*WithBoard
	PlayerIdx      int
	CastsRemaining int
	MessageShown   bool
	CastBefore     bool
}

func (screen *TargetSpellScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	textBottom("", ss, win) // clear bottom bar
}

func (screen *TargetSpellScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	thisPlayer := screen.Players[screen.PlayerIdx]
	spell := thisPlayer.Spells[thisPlayer.ChosenSpell]
	batch := screen.WithBoard.DrawBoard(ss, win)

	if spell.GetCastRange() == 0 {
		target := screen.WithBoard.CursorPosition
		fmt.Printf("Cast spell %s (%d) on V(%d, %d)\n", spell.GetName(), spell.GetCastRange(), target.X, target.Y)
		return screen.AnimateAndCast()
	}
	if screen.WithBoard.MoveCursor(win) || !screen.MessageShown {
		screen.MessageShown = false
		screen.WithBoard.DrawCursor(ss, batch)
	}
	if win.JustPressed(pixelgl.KeyS) {
		target := screen.WithBoard.CursorPosition
		if spell.GetCastRange() < target.Distance(screen.Players[screen.PlayerIdx].BoardPosition) {
			textBottom("Out of range", ss, batch)
			fmt.Printf("Out of range! Spell cast range %d but distance to target is %d\n", spell.GetCastRange(), target.Distance(screen.Players[screen.PlayerIdx].BoardPosition))
			screen.MessageShown = true
		} else {
			if spell.NeedsLineOfSight() && !HaveLineOfSight(screen.Players[screen.PlayerIdx].BoardPosition, screen.WithBoard.CursorPosition, screen.WithBoard.Grid) {
				textBottom("No line of sight", ss, batch)
				screen.MessageShown = true
			} else {
				if spell.CanCast(screen.WithBoard.Grid.GetGameObject(target)) {
					fmt.Printf("Cast spell %s (%d) on V(%d, %d)\n", spell.GetName(), spell.GetCastRange(), target.X, target.Y)
					return screen.AnimateAndCast()
				}
				fmt.Printf("Cannot cast '%s' on non-empty square\n", spell.GetName())
			}
		}
	}
	batch.Draw(win)
	if win.JustPressed(pixelgl.Key0) {
		return NextSpellCastOrMove(screen.PlayerIdx, screen.WithBoard, true)
	}
	return screen
}

func (screen *TargetSpellScreen) AnimateAndCast() GameScreen {
	thisPlayer := screen.Players[screen.PlayerIdx]
	spell := thisPlayer.Spells[thisPlayer.ChosenSpell]
	target := screen.WithBoard.CursorPosition
	anim := spell.CastFx()
	if anim != nil {
		screen.WithBoard.Grid.PlaceGameObject(target, anim)
	}
	return &WaitForFx{
		NextScreen: &DoSpellCast{
			WithBoard:      screen.WithBoard,
			PlayerIdx:      screen.PlayerIdx,
			CastsRemaining: screen.CastsRemaining,
			CastBefore:     screen.CastBefore,
		},
		Grid: screen.WithBoard.Grid,
		Fx:   anim,
	}
}

// This screen does the actual mechanics of the animation
// and then casting the spell once animation is finished

type DoSpellCast struct {
	*WithBoard
	PlayerIdx      int
	CastsRemaining int
	CastBefore     bool
}

func (screen *DoSpellCast) Enter(ss pixel.Picture, win *pixelgl.Window) {
}

func (screen *DoSpellCast) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	castsRemaining := screen.CastsRemaining - 1
	batch := screen.WithBoard.DrawBoard(ss, win)
	// Fx for spell cast finished
	// Work out what happened :)
	targetVec := screen.WithBoard.CursorPosition
	fmt.Printf("About to call player CastSpell method\n")
	p := screen.Players[screen.PlayerIdx]

	fmt.Printf("IN Player spell cast\n")
	spell := p.Spells[p.ChosenSpell]
	fmt.Printf("Player spell %T cast on %T\n", spell, targetVec)
	var success bool
	var canCastMore bool
	var anim *fx.Fx
	if screen.CastBefore || spell.CastSucceeds(p.CastIllusion, screen.LawRating) {
		canCastMore = true
		success, anim = spell.DoCast(p.CastIllusion, targetVec, screen.WithBoard.Grid, p)
	}

	fmt.Printf("Finished player CastSpell method\n")
	if !screen.CastBefore {
		if (canCastMore && castsRemaining > 0) || success {
			fmt.Printf("Spell Succeeds\n")
			textBottom("Spell Succeeds", ss, batch)
			screen.WithBoard.LawRating += spell.GetLawRating()
		} else {
			fmt.Printf("Spell failed\n")
			textBottom("Spell Failed", ss, batch)
		}
	}
	batch.Draw(win)
	nextScreen := NextSpellCastOrMove(screen.PlayerIdx, screen.WithBoard, false)
	if castsRemaining > 0 && canCastMore {
		nextScreen = &TargetSpellScreen{
			WithBoard:      screen.WithBoard,
			PlayerIdx:      screen.PlayerIdx,
			CastsRemaining: castsRemaining,
			CastBefore:     true,
		}
	} else {
		if !spell.IsReuseable() {
			p.Spells = append(p.Spells[:p.ChosenSpell], p.Spells[p.ChosenSpell+1:]...)
		}
		p.ChosenSpell = -1
	}
	return &WaitForFx{
		NextScreen: nextScreen,
		Grid:       screen.Grid,
		Fx:         anim,
	}
}
