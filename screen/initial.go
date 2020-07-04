package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
)

type InitialScreen struct {
	DrawnFirst         bool
	WizardCount        int
	DrawnSecond        bool
	ComputerDifficulty int
}

func (screen *InitialScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
}

func (screen *InitialScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	td := TextDrawer(ss)
	if !screen.DrawnFirst {
		td.DrawText("  MAYHEM - Remake of Chaos", logical.V(0, 9), win)
		td.DrawText("         By bobtfish", logical.V(0, 8), win)
		td.DrawText("How many wizards?", logical.V(0, 6), win)
		td.DrawText("(Press 2 to 8)", logical.V(0, 5), win)
		screen.DrawnFirst = true
	} else {
		if screen.WizardCount == 0 {
			c := captureNumKey(win)
			if c >= 2 && c <= 8 {
				screen.WizardCount = c
				td.DrawText(fmt.Sprintf("%d", c), logical.V(18, 6), win)
			}
		} else {
			if !screen.DrawnSecond {
				td.DrawText("Level of computer wizards?", logical.V(0, 3), win)
				td.DrawText("(Press 1 to 8)", logical.V(0, 2), win)
				screen.DrawnSecond = true
			} else {
				c := captureNumKey(win)
				if c >= 1 && c <= 8 {
					screen.ComputerDifficulty = c
					td.DrawText(fmt.Sprintf("%d", c), logical.V(27, 3), win)
				}
			}
		}
	}
	if screen.ComputerDifficulty > 0 && screen.WizardCount > 0 {
		return &PlayersScreen{
			WizardCount:        screen.WizardCount,
			ComputerDifficulty: screen.ComputerDifficulty,
		}
	}
	return screen
}
