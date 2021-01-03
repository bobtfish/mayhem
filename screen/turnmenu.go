package screen

import (
	"fmt"

	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/game"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
	"github.com/bobtfish/mayhem/spells"
	spelliface "github.com/bobtfish/mayhem/spells/iface"
)

type ExamineOneSpellScreen struct {
	Spell        spelliface.Spell
	ReturnScreen screeniface.GameScreen
}

func (screen *ExamineOneSpellScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	textBottom("   Press any key to continue", ss, win)
	td := TextDrawer(ss)
	td.DrawText(screen.Spell.GetName(), logical.V(0, 9), win)
	td.DrawText("FIXME add stuff per spell", logical.V(0, 7), win)
}

func (screen *ExamineOneSpellScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
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

func (screen *SpellListScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	textBottom("Press 0 to return to main menu", ss, win)
	td := TextDrawer(ss)
	td.DrawText(fmt.Sprintf("%s's spells", screen.Player.Name), logical.V(0, 9), win)
	for i := 0; i < len(screen.Player.Spells); i++ {
		mod := i % 2
		if mod == 1 {
			mod = 14
		}
		spell := screen.Player.Spells[i]
		td.DrawTextColor(
			fmt.Sprintf("%s%s%s", intToChar(i), spelliface.LawRatingSymbol(spell), spell.GetName()),
			logical.V(mod, 8-(i/2)),
			spells.CastingChanceColor(spell.GetCastingChance(ctx.GetLawRating())),
			win,
		)
	}
}

// End SpellListScreen

// Begin Examine Spells list screen
type ExamineSpellsScreen struct {
	SpellListScreen
	SpellToExamine *spelliface.Spell
}

//func (screen *ExamineSpellsScreen) Enter(ctx screeniface.GameCtx) {
//	screen.SpellListScreen.Enter(ss, win)
//}

func (screen *ExamineSpellsScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
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

//func (screen *SelectSpellScreen) Enter(ctx screeniface.GameCtx) {
//	screen.SpellListScreen.Enter(ss, win)
//}

func (screen *SelectSpellScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	c := captureNumKey(win)
	if c == 0 {
		return screen.MainMenu
	}
	i := captureSpellKey(win)
	if i >= 0 && i < len(screen.Player.Spells) {
		screen.Player.ChosenSpell = i
		screen.Player.CastIllusion = false // Remember to reset this
		if screen.Player.Spells[i].CanCastAsIllusion() {
			return &IsIllusionScreen{
				SpellListScreen: screen.SpellListScreen,
			}
		}
		return screen.MainMenu
	}
	return screen
}

// End

type IsIllusionScreen struct {
	SpellListScreen
}

func (screen *IsIllusionScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	textBottom("Illusion? (Press Y or N)", ss, win)
	screen.Player.CastIllusion = false
}

func (screen *IsIllusionScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	if win.JustPressed(pixelgl.KeyY) {
		screen.Player.CastIllusion = true
	}
	if win.JustPressed(pixelgl.KeyY) || win.JustPressed(pixelgl.KeyN) {
		return screen.MainMenu
	}
	return screen
}

// Begin main turn menu screen
type TurnMenuScreen struct {
	PlayerIdx int
}

func (screen *TurnMenuScreen) Enter(ctx screeniface.GameCtx) {
	players := ctx.(*game.Window).GetPlayers()
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	fmt.Printf("index %d\n", screen.PlayerIdx)
	textBottom("      Press Keys 1 to 4", ss, win)
	td := TextDrawer(ss)
	td.DrawText(players[screen.PlayerIdx].Name, logical.V(3, 7), win)
	td.DrawText(lawRatingText(ctx.GetLawRating()), logical.V(3, 6), win)
	td.DrawText("1. Examine Spells", logical.V(3, 5), win)
	td.DrawText("2. Select Spell", logical.V(3, 4), win)
	td.DrawText("3. Examine Board", logical.V(3, 3), win)
	td.DrawText("4. Continue With Game", logical.V(3, 2), win)
}

func lawRatingText(r int) string {
	if r == 0 {
		return ""
	}
	if r > 0 {
		return fmt.Sprintf("(Law %s)", lawRatingSymbolText(r))
	}
	return fmt.Sprintf("(Chaos %s)", lawRatingSymbolText(r))
}

func lawRatingSymbolText(r int) string {
	ar := r / 4 // We display 1 symbol per 10% we changed spell chances
	if ar < 0 {
		ar = -ar
	}
	ra := make([]rune, ar)
	for i := 0; i < ar; i++ {
		if r > 0 {
			ra[i] = '^'
		} else {
			ra[i] = '*'
		}
	}
	return string(ra)
}

func (screen *TurnMenuScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	c := captureNumKey(win)
	players := ctx.(*game.Window).GetPlayers()
	if c == 1 {
		return &ExamineSpellsScreen{
			SpellListScreen: SpellListScreen{
				MainMenu: screen,
				Player:   players[screen.PlayerIdx],
			},
		}
	}
	if c == 2 {
		return &SelectSpellScreen{
			SpellListScreen: SpellListScreen{
				MainMenu: screen,
				Player:   players[screen.PlayerIdx],
			},
		}
	}
	if c == 3 {
		return &ExamineBoardScreen{
			MainMenu:   screen,
			WithCursor: WithCursor{CursorPosition: players[screen.PlayerIdx].BoardPosition},
		}
	}
	if c == 4 {
		if len(players) == screen.PlayerIdx+1 {
			return &DisplaySpellCastScreen{}
		}
		return &TurnMenuScreen{
			PlayerIdx: screen.PlayerIdx + 1,
		}
	}
	return screen
}
