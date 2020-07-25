package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
)

// Move state onto next player spell cast (if there are players left)
// or onto the movement phase if all spells have been cast
func NextSpellCastOrMove(playerIdx int, players []*player.Player, grid *grid.GameGrid, skipPause bool) GameScreen {
	var nextScreen GameScreen
	nextIdx := NextPlayerIdx(playerIdx, players)
	nextScreen = &DisplaySpellCastScreen{
		WithBoard: &WithBoard{
			Grid:    grid,
			Players: players,
		},
		PlayerIdx: nextIdx,
	}

	if nextIdx == len(players) {
		// All players have cast their spells, movement comes next
		nextScreen = &MoveAnnounceScreen{
			WithBoard: &WithBoard{
				Players: players,
				Grid:    grid,
			},
		}
	}
	return &Pause{
		Skip:       skipPause,
		Grid:       grid,
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
}

func (screen *DisplaySpellCastScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	thisPlayer := screen.Players[screen.PlayerIdx]
	if thisPlayer.ChosenSpell < 0 {
		return NextSpellCastOrMove(screen.PlayerIdx, screen.Players, screen.Grid, true)
	}
	spell := thisPlayer.Spells[thisPlayer.ChosenSpell]
	batch := screen.WithBoard.DrawBoard(ss, win)
	textBottom(fmt.Sprintf("%s %s %d", thisPlayer.Name, spell.GetName(), spell.GetCastRange()), ss, batch)
	batch.Draw(win)
	if win.JustPressed(pixelgl.KeyS) || !captureDirectionKey(win).Equals(logical.ZeroVec()) {
		return &TargetSpellScreen{
			WithBoard:     screen.WithBoard,
			PlayerIdx:     screen.PlayerIdx,
			NumberOfCasts: spell.NumberOfCasts(),
		}
	}
	if win.JustPressed(pixelgl.Key0) {
		return NextSpellCastOrMove(screen.PlayerIdx, screen.Players, screen.Grid, false)
	}
	return screen
}

// If range 0 then cast the spell straight away (on the wizard
// If range > 0 then move cursor around to find a target until S is pressed

type TargetSpellScreen struct {
	*WithBoard
	PlayerIdx     int
	OutOfRange    bool
	NumberOfCasts int
}

func (screen *TargetSpellScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	textBottom("", ss, win) // clear bottom bar
	screen.WithBoard.CursorPosition = screen.Players[screen.PlayerIdx].BoardPosition
}

func (screen *TargetSpellScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	thisPlayer := screen.Players[screen.PlayerIdx]
	spell := thisPlayer.Spells[thisPlayer.ChosenSpell]
	batch := screen.WithBoard.DrawBoard(ss, win)

	if spell.GetCastRange() == 0 {
		target := screen.WithBoard.CursorPosition
		fmt.Printf("Cast spell %s (%d) on V(%d, %d)\n", spell.GetName(), spell.GetCastRange(), target.X, target.Y)
		return screen.AnimateAndCast()
	} else {
		if screen.WithBoard.MoveCursor(win) || !screen.OutOfRange {
			screen.OutOfRange = false
			screen.WithBoard.DrawCursor(ss, batch)
		}
		if win.JustPressed(pixelgl.KeyS) {
			target := screen.WithBoard.CursorPosition
			if spell.GetCastRange() < target.Distance(screen.Players[screen.PlayerIdx].BoardPosition) {
				textBottom("Out of range", ss, batch)
				fmt.Printf("Out of range! Spell cast range %d but distance to target is %d\n", spell.GetCastRange(), target.Distance(screen.Players[screen.PlayerIdx].BoardPosition))
				screen.OutOfRange = true
			} else {
				if !HaveLineOfSight(screen.Players[screen.PlayerIdx].BoardPosition, screen.WithBoard.CursorPosition, screen.WithBoard.Grid) {
					textBottom("No line of sight", ss, batch)
					screen.OutOfRange = true
				} else {
					if spell.CanCast(screen.WithBoard.Grid.GetGameObject(target)) {
						fmt.Printf("Cast spell %s (%d) on V(%d, %d)\n", spell.GetName(), spell.GetCastRange(), target.X, target.Y)
						return screen.AnimateAndCast()
					} else {
						fmt.Printf("Cannot cast on non-empty square\n")
					}
				}
			}
		}
	}
	batch.Draw(win)
	if win.JustPressed(pixelgl.Key0) {
		return NextSpellCastOrMove(screen.PlayerIdx, screen.Players, screen.Grid, true)
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
			WithBoard:         screen.WithBoard,
			PlayerIdx:         screen.PlayerIdx,
			NumberOfCastsLeft: screen.NumberOfCasts - 1,
		},
		Grid: screen.WithBoard.Grid,
		Fx:   anim,
	}
}

// This screen does the actual mechanics of the animation
// and then casting the spell once animation is finished

type DoSpellCast struct {
	*WithBoard
	PlayerIdx         int
	NumberOfCastsLeft int
}

func (screen *DoSpellCast) Enter(ss pixel.Picture, win *pixelgl.Window) {
}

func (screen *DoSpellCast) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	var nextScreen GameScreen
	batch := screen.WithBoard.DrawBoard(ss, win)
	// Fx for spell cast finished
	// Work out what happened :)
	targetVec := screen.WithBoard.CursorPosition
	fmt.Printf("About to call player CastSpell method\n")
	success, fx := screen.Players[screen.PlayerIdx].CastSpell(targetVec, screen.WithBoard.Grid)
	fmt.Printf("Finished player CastSpell method\n")
	if success {
		fmt.Printf("Spell Succeeds\n")
		textBottom("Spell Succeeds", ss, batch)
	} else {
		fmt.Printf("Spell failed\n")
		textBottom("Spell Failed", ss, batch)
	}
	batch.Draw(win)
	nextScreen = &TargetSpellScreen{
		WithBoard:     screen.WithBoard,
		PlayerIdx:     screen.PlayerIdx,
		NumberOfCasts: screen.NumberOfCastsLeft,
	}
	if screen.NumberOfCastsLeft == 0 {
		nextScreen = NextSpellCastOrMove(screen.PlayerIdx, screen.Players, screen.Grid, false)
	}
	return &WaitForFx{
		NextScreen: nextScreen,
		Grid:       screen.Grid,
		Fx:         fx,
	}
}
