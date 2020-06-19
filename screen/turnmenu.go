package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

const (
	ChoseNothing       = iota
	ChoseExamineSpells = iota
	ChoseSelectSpell   = iota
	ChoseContinue      = iota
)

type SpellListScreen struct {
	ScreenBasics
	MainMenu *TurnMenuScreen
	Player   *Player
	finished bool
}

func (screen *SpellListScreen) Setup(ss pixel.Picture, win *pixelgl.Window) {
	screen.ScreenBasics.Setup(ss, win)
	render.NewTextDrawer(ss, logical.V(0, 0)).DrawText("Press 0 to return to main menu", logical.V(0, 0), win)
}

func (screen *SpellListScreen) Draw(win *pixelgl.Window) {
	screen.TextDrawer.DrawText(fmt.Sprintf("%s's spells", screen.Player.Name), logical.V(0, 9), win)
	for i := 0; i < len(screen.Player.Spells); i++ {
		screen.TextDrawer.DrawText(fmt.Sprintf("%s%s%s", intToChar(i), screen.Player.Spells[i].LawRatingSymbol(), screen.Player.Spells[i].Name), logical.V(0, 8-i), win)
	}
	c := captureNumKey(win)
	if c == 0 {
		screen.finished = true
	}
}

func (screen *SpellListScreen) NextScreen() GameScreen {
	return screen.MainMenu
}

func (screen *SpellListScreen) Finished() bool {
	return screen.finished
}

type ExamineSpellsScreen struct {
	SpellListScreen
}

func (screen *ExamineSpellsScreen) Setup(ss pixel.Picture, win *pixelgl.Window) {
	screen.SpellListScreen.Setup(ss, win)
}

type SelectSpellScreen struct {
	SpellListScreen
}

func (screen *SelectSpellScreen) Setup(ss pixel.Picture, win *pixelgl.Window) {
	screen.SpellListScreen.Setup(ss, win)
}

type TurnMenuScreen struct {
	Players     []*Player
	PlayerIndex int
	ScreenBasics
	ChosenOption int
}

func (screen *TurnMenuScreen) Setup(ss pixel.Picture, win *pixelgl.Window) {
	fmt.Println(fmt.Sprintf("index %d", screen.PlayerIndex))
	screen.ScreenBasics.Setup(ss, win)
	render.NewTextDrawer(ss, logical.V(0, 0)).DrawText("      Press Keys 1 to 4", logical.V(0, 0), win)
	screen.TextDrawer.DrawText(screen.Players[screen.PlayerIndex].Name, logical.V(3, 7), win)
	screen.TextDrawer.DrawText("1. Examine Spells", logical.V(3, 5), win)
	screen.TextDrawer.DrawText("2. Select Spell", logical.V(3, 4), win)
	screen.TextDrawer.DrawText("3. Examine Board", logical.V(3, 3), win)
	screen.TextDrawer.DrawText("4. Continue With Game", logical.V(3, 2), win)
	screen.ChosenOption = ChoseNothing
}

func (screen *TurnMenuScreen) Draw(win *pixelgl.Window) {
	c := captureNumKey(win)
	if c == 1 {
		screen.ChosenOption = ChoseExamineSpells
	}
	if c == 2 {
		screen.ChosenOption = ChoseSelectSpell
	}
	if c == 4 {
		fmt.Println("Set Continue")
		screen.ChosenOption = ChoseContinue
	}
}

func (screen *TurnMenuScreen) NextScreen() GameScreen {
	fmt.Println("NextScreen")
	if screen.ChosenOption == ChoseContinue {
		if screen.PlayerIndex < len(screen.Players)-1 {
			fmt.Println("NextScreen return screen")
			screen.PlayerIndex++
			screen.ChosenOption = ChoseNothing
			return screen
		}
		fmt.Println("NextScreen return Initial")
		return nil
	}
	if screen.ChosenOption == ChoseSelectSpell {
		newS := &SelectSpellScreen{}
		newS.MainMenu = screen
		newS.Player = screen.Players[screen.PlayerIndex]
		return newS
	}
	if screen.ChosenOption == ChoseExamineSpells {
		newS := &ExamineSpellsScreen{}
		newS.MainMenu = screen
		newS.Player = screen.Players[screen.PlayerIndex]
		return newS
	}
	return nil
}

func (screen *TurnMenuScreen) Finished() bool {
	return screen.ChosenOption != ChoseNothing
}
