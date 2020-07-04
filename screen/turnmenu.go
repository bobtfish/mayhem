package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
	"github.com/bobtfish/mayhem/spells"
)

const (
	ChoseNothing       = iota
	ChoseExamineSpells = iota
	ChoseSelectSpell   = iota
	ChoseContinue      = iota
)

type ExamineOneSpellScreen struct {
	Spell        *spells.Spell
	ReturnScreen *ExamineSpellsScreen
	finished     bool
}

func (screen *ExamineOneSpellScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	render.NewTextDrawer(ss).DrawText("Press any key to continue", logical.V(0, 0), win)
	td := TextDrawer(ss)
	td.DrawText(screen.Spell.Name, logical.V(0, 9), win)
	td.DrawText("FIXME add stuff per spell", logical.V(0, 7), win)
}

func (screen *ExamineOneSpellScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	if win.Typed() != "" {
		return screen.ReturnScreen
	}
	return screen
}

// Shared SpellListScreen is common functionality
type SpellListScreen struct {
	MainMenu *TurnMenuScreen
	Player   *Player
	finished bool
}

func (screen *SpellListScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	render.NewTextDrawer(ss).DrawText("Press 0 to return to main menu", logical.V(0, 0), win)
}

func (screen *SpellListScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	td := TextDrawer(ss)
	td.DrawText(fmt.Sprintf("%s's spells", screen.Player.Name), logical.V(0, 9), win)
	for i := 0; i < len(screen.Player.Spells); i++ {
		td.DrawText(fmt.Sprintf("%s%s%s", intToChar(i), screen.Player.Spells[i].LawRatingSymbol(), screen.Player.Spells[i].Name), logical.V(0, 8-i), win)
	}
	c := captureNumKey(win)
	if c == 0 {
		return screen.MainMenu
	}
	return screen
}

// End SpellListScreen

// Begin Examine Spells list screen
type ExamineSpellsScreen struct {
	SpellListScreen
	SpellToExamine *spells.Spell
}

//func (screen *ExamineSpellsScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
//	screen.SpellListScreen.Enter(ss, win)
//}

func (screen *ExamineSpellsScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	screen.SpellListScreen.Step(ss, win)
	i := captureSpellKey(win)
	if i >= 0 && i < len(screen.Player.Spells) {
		fmt.Println("Examine a spell")
		return &ExamineOneSpellScreen{
			Spell:        &screen.Player.Spells[i],
			ReturnScreen: screen,
		}
	}
	return screen
}

// End

// Being Select a Spell screen
type SelectSpellScreen struct {
	SpellListScreen
}

//func (screen *SelectSpellScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
//	screen.SpellListScreen.Enter(ss, win)
//}

func (screen *SelectSpellScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	screen.SpellListScreen.Step(ss, win)
	i := captureSpellKey(win)
	if i >= 0 && i < len(screen.Player.Spells) {
		screen.Player.ChosenSpell = i
		fmt.Println("Chose a spell")
		return screen.MainMenu
	}
	return screen
}

// End

// Begin main turn menu screen
type TurnMenuScreen struct {
	Players      []*Player
	PlayerIndex  int
	ChosenOption int
}

func (screen *TurnMenuScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	fmt.Println(fmt.Sprintf("index %d", screen.PlayerIndex))
	render.NewTextDrawer(ss).DrawText("      Press Keys 1 to 4", logical.V(0, 0), win)
	td := TextDrawer(ss)
	td.DrawText(screen.Players[screen.PlayerIndex].Name, logical.V(3, 7), win)
	td.DrawText("1. Examine Spells", logical.V(3, 5), win)
	td.DrawText("2. Select Spell", logical.V(3, 4), win)
	td.DrawText("3. Examine Board", logical.V(3, 3), win)
	td.DrawText("4. Continue With Game", logical.V(3, 2), win)
	screen.ChosenOption = ChoseNothing
}

func (screen *TurnMenuScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	c := captureNumKey(win)
	if c == 1 {
		newS := &ExamineSpellsScreen{}
		newS.MainMenu = screen
		newS.Player = screen.Players[screen.PlayerIndex]
		return newS
	}
	if c == 2 {
		newS := &SelectSpellScreen{}
		newS.MainMenu = screen
		newS.Player = screen.Players[screen.PlayerIndex]
		return newS
	}
	if c == 4 {
		fmt.Println("Set Continue")
		screen.PlayerIndex++
		screen.ChosenOption = ChoseNothing
		return screen
	}
	return screen
}
