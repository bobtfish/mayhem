package screen

import (
	"fmt"
	"strings"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

type WinnerScreen struct {
	Players []*player.Player
}

func (screen *WinnerScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	td := TextDrawer(ss)
	td.DrawText("  WE HAVE A WINNER", logical.V(0, 9), win)

	var winner *player.Player
	for i := 0; i < len(screen.Players); i++ {
		if screen.Players[i].Alive {
			winner = screen.Players[i]
			break
		}
	}
	spaceLen := 16 - len(winner.Name)/2
	td.DrawText(fmt.Sprintf("%s%s", strings.Repeat(" ", spaceLen), winner.Name), logical.V(0, 8), win)
}

func (screen *WinnerScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	if win.Typed() != "" {
		return &InitialScreen{}
	}
	return screen
}
