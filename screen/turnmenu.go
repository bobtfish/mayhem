package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/spells"
)

type ExamineOneSpellScreen struct {
	Spell        spells.Spell
	ReturnScreen GameScreen
}

func (screen *ExamineOneSpellScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	textBottom("Press any key to continue", ss, win)
	td := TextDrawer(ss)
	td.DrawText(screen.Spell.GetName(), logical.V(0, 9), win)
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
	Player   *player.Player
}

func (screen *SpellListScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	textBottom("Press 0 to return to main menu", ss, win)
	td := TextDrawer(ss)
	td.DrawText(fmt.Sprintf("%s's spells", screen.Player.Name), logical.V(0, 9), win)
	for i := 0; i < len(screen.Player.Spells); i++ {
		mod := i % 2
		if mod == 1 {
			mod = 15
		}
		spell := screen.Player.Spells[i]
		td.DrawTextColor(
			fmt.Sprintf("%s%s%s", intToChar(i), spells.LawRatingSymbol(spell), spell.GetName()),
			logical.V(mod, 8-(i/2)),
			spells.CastingChanceColor(spell.GetCastingChance(screen.Player.LawRating)),
			win,
		)
	}
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
	c := captureNumKey(win)
	if c == 0 {
		return screen.MainMenu
	}
	i := captureSpellKey(win)
	if i >= 0 && i < len(screen.Player.Spells) {
		return &ExamineOneSpellScreen{
			Spell:        screen.Player.Spells[i],
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
	c := captureNumKey(win)
	if c == 0 {
		return screen.MainMenu
	}
	i := captureSpellKey(win)
	if i >= 0 && i < len(screen.Player.Spells) {
		screen.Player.ChosenSpell = i
		return screen.MainMenu
	}
	return screen
}

// End

// Begin main turn menu screen
type TurnMenuScreen struct {
	Players   []*player.Player
	PlayerIdx int
	Grid      *grid.GameGrid
}

func (screen *TurnMenuScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	fmt.Println(fmt.Sprintf("index %d", screen.PlayerIdx))
	textBottom("      Press Keys 1 to 4", ss, win)
	td := TextDrawer(ss)
	td.DrawText(screen.Players[screen.PlayerIdx].Name, logical.V(3, 7), win)
	td.DrawText("1. Examine Spells", logical.V(3, 5), win)
	td.DrawText("2. Select Spell", logical.V(3, 4), win)
	td.DrawText("3. Examine Board", logical.V(3, 3), win)
	td.DrawText("4. Continue With Game", logical.V(3, 2), win)
}

func (screen *TurnMenuScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	c := captureNumKey(win)
	if c == 1 {
		return &ExamineSpellsScreen{
			SpellListScreen: SpellListScreen{
				MainMenu: screen,
				Player:   screen.Players[screen.PlayerIdx],
			},
		}
	}
	if c == 2 {
		return &SelectSpellScreen{
			SpellListScreen: SpellListScreen{
				MainMenu: screen,
				Player:   screen.Players[screen.PlayerIdx],
			},
		}
	}
	if c == 3 {
		return &ExamineBoardScreen{
			MainMenu: screen,
			WithBoard: &WithBoard{
				Grid:    screen.Grid,
				Players: screen.Players,
			},
		}
	}
	if c == 4 {
		if len(screen.Players) == screen.PlayerIdx+1 {
			return &DisplaySpellCastScreen{
				WithBoard: &WithBoard{
					Grid:    screen.Grid,
					Players: screen.Players,
				},
			}
		}
		return &TurnMenuScreen{
			Players:   screen.Players,
			PlayerIdx: screen.PlayerIdx + 1,
			Grid:      screen.Grid,
		}
	}
	return screen
}
