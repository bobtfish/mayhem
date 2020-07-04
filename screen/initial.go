package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
)

type InitialScreen struct{}
type ComputerDifficultyScreen struct {
	WizardCount int
}

func (screen *InitialScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	td := TextDrawer(ss)
	td.DrawText("  MAYHEM - Remake of Chaos", logical.V(0, 9), win)
	td.DrawText("         By bobtfish", logical.V(0, 8), win)
	td.DrawText("How many wizards?", logical.V(0, 6), win)
	td.DrawText("(Press 2 to 8)", logical.V(0, 5), win)
}

func (screen *InitialScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	c := captureNumKey(win)
	if c >= 2 && c <= 8 {
		td := TextDrawer(ss)
		td.DrawText(fmt.Sprintf("%d", c), logical.V(18, 6), win)
		td.DrawText("Level of computer wizards?", logical.V(0, 3), win)
		td.DrawText("(Press 1 to 8)", logical.V(0, 2), win)
		return &ComputerDifficultyScreen{WizardCount: c}
	}
	return screen
}

func (screen *ComputerDifficultyScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	td := TextDrawer(ss)
	td.DrawText("Level of computer wizards?", logical.V(0, 3), win)
	td.DrawText("(Press 1 to 8)", logical.V(0, 2), win)
}

func (screen *ComputerDifficultyScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	c := captureNumKey(win)
	if c >= 1 && c <= 8 {
		TextDrawer(ss).DrawText(fmt.Sprintf("%d", c), logical.V(27, 3), win)
		return &PlayerNameScreen{PlayersScreen: PlayersScreen{
			WizardCount:        screen.WizardCount,
			ComputerDifficulty: c,
		},
		}
	}
	return screen
}
