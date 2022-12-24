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
	td.DrawTextColor("  MAYHEM - Remake of Chaos", logical.V(0, 9), render.ColorWhite(), win)
	td.DrawTextColor("         By bobtfish", logical.V(0, 8), render.ColorWhite(), win)
	td.DrawTextColor("How many wizards?", logical.V(0, 6), render.ColorWhite(), win)
	td.DrawTextColor("(Press 2 to 8)", logical.V(0, 5), render.ColorWhite(), win)
	textBottomColor("       Press H for help", render.ColorWhite(), ss, win)
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
		td.DrawTextColor(fmt.Sprintf("%d", c), logical.V(18, 6), render.ColorWhite(), win)
		td.DrawTextColor("Level of computer wizards?", logical.V(0, 3), render.ColorWhite(), win)
		td.DrawTextColor("(Press 1 to 8)", logical.V(0, 2), render.ColorWhite(), win)
		textBottomColor("", render.ColorWhite(), ss, win)
		return &ComputerDifficultyScreen{WizardCount: c}
	}
	return screen
}

func (screen *ComputerDifficultyScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	td := TextDrawer(ss)
	td.DrawTextColor("Level of computer wizards?", logical.V(0, 3), render.ColorWhite(), win)
	td.DrawTextColor("(Press 1 to 8)", logical.V(0, 2), render.ColorWhite(), win)
}

func (screen *ComputerDifficultyScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	ss := ctx.GetSpriteSheet()
	win := ctx.GetWindow()
	c := captureNumKey(win)
	if c >= 1 && c <= 8 {
		TextDrawer(ss).DrawTextColor(fmt.Sprintf("%d", c), logical.V(27, 3), render.ColorWhite(), win)
		return &PlayerNameScreen{
			PlayersScreen: PlayersScreen{
				WizardCount:        screen.WizardCount,
				ComputerDifficulty: c,
			},
		}
	}
	return screen
}
