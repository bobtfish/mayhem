package screen

import (
    "fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

type TurnMenuScreen struct {
	Players       []*Player
	SelectedIndex int
	ScreenBasics
	Continue bool
}

func (screen *TurnMenuScreen) Setup(ss pixel.Picture, win *pixelgl.Window) {
    fmt.Println(fmt.Sprintf("index %d", screen.SelectedIndex))
	screen.ScreenBasics.Setup(ss, win)
	render.NewTextDrawer(ss, logical.V(0, 0)).DrawText("      Press Keys 1 to 4", logical.V(0, 0), win)
	screen.TextDrawer.DrawText(screen.Players[screen.SelectedIndex].Name, logical.V(3, 7), win)
	screen.TextDrawer.DrawText("1. Examine Spells", logical.V(3, 5), win)
	screen.TextDrawer.DrawText("2. Select Spell", logical.V(3, 4), win)
	screen.TextDrawer.DrawText("3. Examine Board", logical.V(3, 3), win)
	screen.TextDrawer.DrawText("4. Continue With Game", logical.V(3, 2), win)
}

func (screen *TurnMenuScreen) Draw(win *pixelgl.Window) {
	c := captureNumKey(win)
	if c == 4 {
        fmt.Println("Set Continue")
		screen.Continue = true
	}
}

func (screen *TurnMenuScreen) NextScreen() GameScreen {
    fmt.Println("NextScreen")
	if screen.SelectedIndex < len(screen.Players)-1 {
        fmt.Println("NextScreen return screen")
		screen.SelectedIndex++
		screen.Continue = false
		return screen
	}
    fmt.Println("NextScreen return Initial")
	return nil
}

func (screen *TurnMenuScreen) Finished() bool {
	return screen.Continue
}
