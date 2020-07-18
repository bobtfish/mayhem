package screen

import (
	"fmt"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
)

type WinnerScreen struct {
	*WithBoard
}

func (screen *WinnerScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	td := TextDrawer(ss)
	td.DrawText("  WE HAVE A WINNER", logical.V(0, 9), win)

	var winner *player.Player
	for i := 0; i < len(screen.WithBoard.Players); i++ {
		if screen.WithBoard.Players[i].Alive {
			winner = screen.WithBoard.Players[i]
			break
		}
	}
	spaceLen := 16 - len(winner.Name)/2
	td.DrawText(fmt.Sprintf("%s%s", strings.Repeat(" ", spaceLen), winner.Name), logical.V(0, 8), win)
}

func (screen *WinnerScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	if win.Typed() != "" {
		return &InitialScreen{}
	}
	return screen
}
