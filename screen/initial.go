package screen

import (
	"fmt"

	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

type InitialScreen struct{}
type ComputerDifficultyScreen struct {
	WizardCount int
}

func (screen *InitialScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	td := TextDrawer(ss)
	td.DrawText("  MAYHEM - Remake of Chaos", logical.V(0, 9), render.ColorPurple(), win)
	td.DrawText("         By bobtfish", logical.V(0, 8), render.ColorRed(), win)
	td.DrawText("How many wizards?", logical.V(0, 6), render.ColorYellow(), win)
	td.DrawText("(Press 2 to 8)", logical.V(0, 5), render.ColorGreen(), win)
	textBottom("       Press H for help", render.ColorWhite(), ss, win)
}

func (screen *InitialScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	if win.JustPressed(pixelgl.KeyH) {
		return &HelpScreenMenu{}
	}

	c := captureNumKey(win)
	if c >= 2 && c <= 8 {
		td := TextDrawer(ss)
		td.DrawText(fmt.Sprintf("%d", c), logical.V(18, 6), render.ColorWhite(), win)
		td.DrawText("Level of computer wizards?", logical.V(0, 3), render.ColorYellow(), win)
		td.DrawText("(Press 1 to 8)", logical.V(0, 2), render.ColorPurple(), win)
		textBottom("", render.ColorWhite(), ss, win)
		return &ComputerDifficultyScreen{WizardCount: c}
	}
	return screen
}

func (screen *ComputerDifficultyScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	td := TextDrawer(ss)
	td.DrawText("Level of computer wizards?", logical.V(0, 3), render.ColorWhite(), win)
	td.DrawText("(Press 1 to 8)", logical.V(0, 2), render.ColorWhite(), win)
}

func (screen *ComputerDifficultyScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	ss := ctx.GetSpriteSheet()
	win := ctx.GetWindow()
	c := captureNumKey(win)
	if c >= 1 && c <= 8 {
		TextDrawer(ss).DrawText(fmt.Sprintf("%d", c), logical.V(27, 3), render.ColorWhite(), win)
		return &PlayerNameScreen{
			PlayersScreen: PlayersScreen{
				WizardCount:        screen.WizardCount,
				ComputerDifficulty: c,
			},
		}
	}
	return screen
}
