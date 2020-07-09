package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

type CastSpellScreen struct {
	*WithBoard
	PlayerIdx   int
	ReadyToCast bool
}

func (screen *CastSpellScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	screen.WithBoard.CursorPosition = screen.Players[screen.PlayerIdx].BoardPosition
}

func (screen *CastSpellScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	thisPlayer := screen.Players[screen.PlayerIdx]
	spell := thisPlayer.Spells[thisPlayer.ChosenSpell]

	batch := screen.WithBoard.DrawBoard(ss, win)

	/* For spells with range = 0 any key aqwedcxzs casts
	   For spells with range > 0
	     cursor not visible till aqwedcxzs pressed, then shown
	     then standard move cursor till s pressed for cast */
	if !screen.ReadyToCast {
		render.NewTextDrawer(ss).DrawText(fmt.Sprintf("%s %s %d", thisPlayer.Name, spell.GetName(), spell.GetRange()), logical.V(0, 0), win)
		if win.JustPressed(pixelgl.KeyS) || !captureDirectionKey(win).Equals(logical.ZeroVec()) {
			render.NewTextDrawer(ss).DrawText("                                  ", logical.V(0, 0), win) // clear bottom bar
			if spell.GetRange() == 0 {
				target := screen.WithBoard.CursorPosition
				fmt.Printf("Cast spell %s (%d) on V(%d, %d)\n", spell.GetName(), spell.GetRange(), target.X, target.Y)
				return screen.Cast()
			} else {
				screen.ReadyToCast = true
			}
		}
	} else {
		screen.WithBoard.MoveCursor(ss, win, batch)
		// FIXME does bottom bar text update when you move over something?
		if win.JustPressed(pixelgl.KeyS) {
			target := screen.WithBoard.CursorPosition
			fmt.Printf("Cast spell %s (%d) on V(%d, %d)\n", spell.GetName(), spell.GetRange(), target.X, target.Y)
			return screen.Cast()
		}
	}
	batch.Draw(win)

	return screen
}

func (screen *CastSpellScreen) Cast() GameScreen {
	target := screen.WithBoard.CursorPosition
	anim := fx.FxSpellCast()
	screen.WithBoard.Grid.PlaceGameObject(target, anim)
	return &DoSpellCast{
		WithBoard: screen.WithBoard,
		PlayerIdx: screen.PlayerIdx,
		Fx:        anim,
	}
}

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
	if screen.Fx.RemoveMe() {
		// Fx for spell cast finished
		// Work out what happened :)
		return screen.NextSpellOrMove()
	}
	return screen
}

func (screen *DoSpellCast) NextSpellOrMove() GameScreen {
	screen.PlayerIdx++
	if screen.PlayerIdx == len(screen.WithBoard.Players) {
		panic("Not written yet")
	}
	return &Pause{NextScreen: &CastSpellScreen{
		WithBoard: screen.WithBoard,
		PlayerIdx: screen.PlayerIdx,
	}}
}
