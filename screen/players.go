package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
)

type PlayersScreen struct {
	ScreenBasics
	WizardCount        int
	ComputerDifficulty int

	Players       []Player
	CurrentPlayer Player
}

type Player struct {
	Name         string
	NameFinished bool
	HumanPlayer  bool
	AIFinished   bool
}

func (screen *PlayersScreen) Setup(ss pixel.Picture, win *pixelgl.Window) {
	screen.ScreenBasics.Setup(ss, win)
	if screen.Players == nil {
		screen.Players = make([]Player, 0)
	}
}

func (screen *PlayersScreen) Draw(win *pixelgl.Window) {
	if !screen.CurrentPlayer.NameFinished {
		playerCount := len(screen.Players) + 1

		if win.JustPressed(pixelgl.KeyEnter) && len(screen.CurrentPlayer.Name) >= 0 {
			screen.CurrentPlayer.NameFinished = true
		} else {
			if win.JustPressed(pixelgl.KeyBackspace) && len(screen.CurrentPlayer.Name) >= 0 {
				length := len(screen.CurrentPlayer.Name) - 1
				screen.CurrentPlayer.Name = screen.CurrentPlayer.Name[:length]
			} else {
				screen.CurrentPlayer.Name = fmt.Sprintf("%s%s", screen.CurrentPlayer.Name, win.Typed())
			}
		}

		if len(screen.CurrentPlayer.Name) > 12 {
			screen.CurrentPlayer.Name = screen.CurrentPlayer.Name[:12]
		}

		screen.TextDrawer.DrawText(fmt.Sprintf("PLAYER %d", playerCount), logical.V(0, 9), win)
		screen.TextDrawer.DrawText("Enter name (12 letters max.)", logical.V(0, 8), win)
		screen.TextDrawer.DrawText("            ", logical.V(0, 7), win)
		screen.TextDrawer.DrawText(screen.CurrentPlayer.Name, logical.V(0, 7), win)

		if win.JustPressed(pixelgl.KeyEnter) && len(screen.CurrentPlayer.Name) >= 0 {
			screen.CurrentPlayer.NameFinished = true
		}
	} else {
		if !screen.CurrentPlayer.AIFinished {
			screen.TextDrawer.DrawText("Computer controlled?", logical.V(0, 5), win)
			if win.JustPressed(pixelgl.KeyY) {
				screen.CurrentPlayer.AIFinished = true
				screen.TextDrawer.DrawText("Y", logical.V(21, 5), win)
			}
			if win.JustPressed(pixelgl.KeyN) {
				screen.CurrentPlayer.AIFinished = true
				screen.CurrentPlayer.HumanPlayer = true
				screen.TextDrawer.DrawText("N", logical.V(21, 5), win)
			}
		}
	}
}

func (screen *PlayersScreen) NextScreen() GameScreen {
	screen.Players = append(screen.Players, screen.CurrentPlayer)
	screen.CurrentPlayer = Player{}
	if len(screen.Players) == screen.WizardCount {
		return &InitialScreen{}
	}
	return screen
}

func (screen *PlayersScreen) Finished() bool {
	if len(screen.Players) == screen.WizardCount {
		return true
	}
	return false
}
