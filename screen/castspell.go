package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type CastSpellScreen struct {
	*WithBoard
	PlayerIdx int
}

func (screen *CastSpellScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
}

func (screen *CastSpellScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	thisPlayer := screen.Players[screen.PlayerIdx]
	spell := thisPlayer.Spells[thisPlayer.ChosenSpell]

	screen.WithBoard.DrawBoard(ss, win)

	// Text ${playername} ${spellname} ${range}
	// For spells with range = 0
	//    Any aqwedcxzs key casts

	// For spells with range > 0
	//   Cursor not visible until first button, then starts on player
	//   Text dissapears when cursor appears

	if win.JustPressed(pixelgl.KeyS) {
		target := screen.WithBoard.CursorPosition
		fmt.Printf("Cast spell %s (%d) on V(%d, %d)\n", spell.Name, spell.Range, target.X, target.Y)
	}

	return screen
}
