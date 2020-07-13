package screen

import (
	"fmt"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/render"
)

// Move state onto next player spell cast (if there are players left)
// or onto the movement phase if all spells have been cast
func NextSpellCastOrMove(playerIdx int, players []*player.Player, grid *grid.GameGrid) GameScreen {
	playerIdx++
	if playerIdx == len(players) {
		// All players have cast their spells, movement comes next
		return &Pause{
			Grid: grid,
			NextScreen: &TurnMenuScreen{ // FIXME - skip straight to the next turn
				Players: players,
				Grid:    grid,
			},
		}
	}
	return &Pause{
		Grid: grid,
		NextScreen: &DisplaySpellCastScreen{
			WithBoard: &WithBoard{
				Grid:    grid,
				Players: players,
			},
			PlayerIdx: playerIdx,
		},
	}
}

// Display the spell name in the bottom bar until player presses a direction key or s

type DisplaySpellCastScreen struct {
	*WithBoard
	PlayerIdx int
}

func (screen *DisplaySpellCastScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
}

func (screen *DisplaySpellCastScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	thisPlayer := screen.Players[screen.PlayerIdx]
	spell := thisPlayer.Spells[thisPlayer.ChosenSpell]
	if thisPlayer.ChosenSpell < 0 {
		return NextSpellCastOrMove(screen.PlayerIdx, screen.Players, screen.Grid)
	}
	batch := screen.WithBoard.DrawBoard(ss, win)
	render.NewTextDrawer(ss).DrawText(fmt.Sprintf("%s %s %d", thisPlayer.Name, spell.GetName(), spell.GetRange()), logical.V(0, 0), batch)
	batch.Draw(win)
	if win.JustPressed(pixelgl.KeyS) || !captureDirectionKey(win).Equals(logical.ZeroVec()) {
		return &TargetSpellScreen{
			WithBoard: screen.WithBoard,
			PlayerIdx: screen.PlayerIdx,
		}
	}
	if win.JustPressed(pixelgl.Key0) {
		return NextSpellCastOrMove(screen.PlayerIdx, screen.Players, screen.Grid)
	}
	return screen
}

// If range 0 then cast the spell straight away (on the wizard
// If range > 0 then move cursor around to find a target until S is pressed

type TargetSpellScreen struct {
	*WithBoard
	PlayerIdx  int
	OutOfRange bool
}

func (screen *TargetSpellScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	render.NewTextDrawer(ss).DrawText(strings.Repeat(" ", 32), logical.ZeroVec(), win) // clear bottom bar
	screen.WithBoard.CursorPosition = screen.Players[screen.PlayerIdx].BoardPosition
}

func (screen *TargetSpellScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	thisPlayer := screen.Players[screen.PlayerIdx]
	spell := thisPlayer.Spells[thisPlayer.ChosenSpell]
	batch := screen.WithBoard.DrawBoard(ss, win)

	if spell.GetRange() == 0 {
		target := screen.WithBoard.CursorPosition
		fmt.Printf("Cast spell %s (%d) on V(%d, %d)\n", spell.GetName(), spell.GetRange(), target.X, target.Y)
		return screen.AnimateAndCast()
	} else {
		if screen.WithBoard.MoveCursor(win) || !screen.OutOfRange {
			screen.OutOfRange = false
			screen.WithBoard.DrawCursor(ss, batch)
		}
		// FIXME does bottom bar text update when you move over something?
		if win.JustPressed(pixelgl.KeyS) {
			target := screen.WithBoard.CursorPosition
			// FIXME can we cast the spell here?
			if spell.GetRange() < target.Distance(screen.Players[screen.PlayerIdx].BoardPosition) {
				render.NewTextDrawer(ss).DrawText("Out of range", logical.ZeroVec(), batch)
				screen.OutOfRange = true
			} else {
				if spell.CanCast(screen.WithBoard.Grid.GetGameObject(target)) {
					fmt.Printf("Cast spell %s (%d) on V(%d, %d)\n", spell.GetName(), spell.GetRange(), target.X, target.Y)
					return screen.AnimateAndCast()
				} else {
					fmt.Printf("Cannot cast on non-empty square\n")
				}
			}
		}
	}
	batch.Draw(win)
	if win.JustPressed(pixelgl.Key0) {
		return NextSpellCastOrMove(screen.PlayerIdx, screen.Players, screen.Grid)
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
	return &DoSpellCast{
		WithBoard: screen.WithBoard,
		PlayerIdx: screen.PlayerIdx,
		Fx:        anim,
	}
}

// This screen does the actual mechanics of the animation
// and then casting the spell once animation is finished

type DoSpellCast struct {
	*WithBoard
	Fx        *fx.Fx
	PlayerIdx int
}

func (screen *DoSpellCast) Enter(ss pixel.Picture, win *pixelgl.Window) {
}

func (screen *DoSpellCast) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	batch.Draw(win)
	// Wait until the spell animation is finished
	if screen.Fx == nil || screen.Fx.RemoveMe() {
		// Fx for spell cast finished
		// Work out what happened :)
		targetVec := screen.WithBoard.CursorPosition
		fmt.Printf("About to call player CastSpell method\n")
		success := screen.Players[screen.PlayerIdx].CastSpell(targetVec, screen.WithBoard.Grid)
		fmt.Printf("Finished player CastSpell method\n")
		if success {
			fmt.Printf("Spell Succeeds\n")
			render.NewTextDrawer(ss).DrawText("Spell Succeeds", logical.V(0, 0), win)
		} else {
			fmt.Printf("Spell failed\n")
			render.NewTextDrawer(ss).DrawText("Spell Failed", logical.V(0, 0), win)
		}
		return NextSpellCastOrMove(screen.PlayerIdx, screen.Players, screen.Grid)
	}
	return screen
}
