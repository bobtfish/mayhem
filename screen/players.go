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
	Players            []Player
	CurrentPlayer      Player
}

type Player struct {
	Name string
}

func (screen *PlayersScreen) Setup(ss pixel.Picture, win *pixelgl.Window) {
	screen.ScreenBasics.Setup(ss, win)
	screen.Players = make([]Player, 0)
}

func (screen *PlayersScreen) Draw(win *pixelgl.Window) {
	if screen.CurrentPlayer.Name == "" {
		playerCount := len(screen.Players) + 1
		screen.TextDrawer.DrawText(fmt.Sprintf("PLAYER %d", playerCount), logical.V(0, 9), win)
		screen.TextDrawer.DrawText("Enter name (12 letters max.)", logical.V(0, 8), win)
	}
	/*
		if !screen.DrawnFirst {
			screen.TextDrawer.DrawText("  MAYHEM - Remake of Chaos", logical.V(0, 9), win)
			screen.TextDrawer.DrawText("         By bobtfish", logical.V(0, 8), win)
			screen.TextDrawer.DrawText("How many wizards?", logical.V(0, 6), win)
			screen.TextDrawer.DrawText("(Press 2 to 8)", logical.V(0, 5), win)
			screen.DrawnFirst = true
		} else {
			if screen.WizardCount == 0 {
				c := captureNumKey(win)
				if c >= 2 && c <= 8 {
					screen.WizardCount = c
					screen.TextDrawer.DrawText(fmt.Sprintf("%d", c), logical.V(18, 6), win)
				}
			} else {
				if !screen.DrawnSecond {
					screen.TextDrawer.DrawText("Level of computer wizards?", logical.V(0, 3), win)
					screen.TextDrawer.DrawText("(Press 1 to 8)", logical.V(0, 2), win)
					screen.DrawnSecond = true
				} else {
					c := captureNumKey(win)
					if c >= 1 && c <= 8 {
						screen.ComputerDifficulty = c
						screen.TextDrawer.DrawText(fmt.Sprintf("%d", c), logical.V(27, 3), win)
					}
				}
			}
		}*/
}

func (screen *PlayersScreen) NextScreen() GameScreen {
	return &InitialScreen{}
}

func (screen *PlayersScreen) Finished() bool {
	/*	if screen.ComputerDifficulty > 0 && screen.WizardCount > 0 {
		return true
	}*/
	return false
}
