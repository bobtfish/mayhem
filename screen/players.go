package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"

	"github.com/bobtfish/mayhem/spells"
)

type PlayersScreen struct {
	WizardCount        int
	ComputerDifficulty int

	Players       []Player
	CurrentPlayer Player
}

type Player struct {
	Name          string
	HumanPlayer   bool
	CharacterIcon logical.Vec

	Spells      []spells.Spell
	ChosenSpell int
}

type PlayerNameScreen struct {
	PlayersScreen
}

func (screen *PlayerNameScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	if screen.Players == nil {
		screen.Players = make([]Player, 0)
	}
	td := TextDrawer(ss)
	td.DrawText(fmt.Sprintf("PLAYER %d", screen.WizardCount), logical.V(0, 9), win)
	td.DrawText("Enter name (12 letters max.)", logical.V(0, 8), win)
}

func (screen *PlayerNameScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	if win.JustPressed(pixelgl.KeyEnter) && len(screen.CurrentPlayer.Name) >= 0 {
		return &PlayerAIScreen{PlayersScreen: screen.PlayersScreen}
	}
	if win.JustPressed(pixelgl.KeyBackspace) && len(screen.CurrentPlayer.Name) > 0 {
		length := len(screen.CurrentPlayer.Name) - 1
		screen.CurrentPlayer.Name = screen.CurrentPlayer.Name[:length]
	} else {
		screen.CurrentPlayer.Name = fmt.Sprintf("%s%s", screen.CurrentPlayer.Name, win.Typed())
	}
	if len(screen.CurrentPlayer.Name) > 12 {
		screen.CurrentPlayer.Name = screen.CurrentPlayer.Name[:12]
	}
	td := TextDrawer(ss)
	td.DrawText("            ", logical.V(0, 7), win)
	td.DrawText(screen.CurrentPlayer.Name, logical.V(0, 7), win)
	return screen
}

type PlayerAIScreen struct {
	PlayersScreen
}

func (screen *PlayerAIScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	TextDrawer(ss).DrawText("Computer controlled?", logical.V(0, 5), win)
}

func (screen *PlayerAIScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	td := TextDrawer(ss)
	if win.JustPressed(pixelgl.KeyY) || win.JustPressed(pixelgl.KeyN) {
		if win.JustPressed(pixelgl.KeyY) {
			td.DrawText("Y", logical.V(21, 5), win)
		}
		if win.JustPressed(pixelgl.KeyN) {
			screen.CurrentPlayer.HumanPlayer = true
			td.DrawText("N", logical.V(21, 5), win)
		}
		return &PlayerIconScreen{PlayersScreen: screen.PlayersScreen}
	}
	return screen
}

type PlayerIconScreen struct {
	PlayersScreen
}

func (screen *PlayerIconScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	td := TextDrawer(ss)
	sd := SpriteDrawer(ss)
	td.DrawText("Which character?", logical.V(0, 4), win)
	td.DrawText("1  2  3  4  5  6  7  8", logical.V(0, 3), win)
	for x := 0; x < 8; x++ {
		offset := logical.V(render.CHAR_PIXELS/4+render.CHAR_PIXELS/2*x*3, render.CHAR_PIXELS*2-render.CHAR_PIXELS/2)
		sd.WithOffset(offset).DrawSprite(logical.V(x, 23), logical.V(1, 3), win)
	}
}

func (screen *PlayerIconScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	c := captureNumKey(win)
	if c >= 1 && c <= 8 {
		screen.CurrentPlayer.CharacterIcon = logical.V(c-1, 23)
		TextDrawer(ss).DrawText(fmt.Sprintf("%d", c), logical.V(17, 4), win)
		sd := SpriteDrawer(ss)
		offset := sd.WinOffsetV.Add(logical.V(render.CHAR_PIXELS/4, 0))
		sd.WithOffset(offset).DrawSprite(screen.CurrentPlayer.CharacterIcon, logical.V(9, 4), win)
		return &PlayerColorScreen{PlayersScreen: screen.PlayersScreen}
	}
	return screen
}

type PlayerColorScreen struct {
	PlayersScreen
}

func (screen *PlayerColorScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	td := TextDrawer(ss)
	sd := SpriteDrawer(ss)
	td.DrawText("Which color?", logical.V(0, 2), win)
	td.DrawText("1  2  3  4  5  6  7  8", logical.V(0, 1), win)
	for x := 0; x < 8; x++ {
		offset := logical.V(render.CHAR_PIXELS/4+render.CHAR_PIXELS/2*x*3, render.CHAR_PIXELS*2-render.CHAR_PIXELS/2)
		sd.WithOffset(offset).DrawSprite(screen.CurrentPlayer.CharacterIcon, logical.V(1, 1), win)
	}
}
func (screen *PlayerColorScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	td := TextDrawer(ss)
	sd := SpriteDrawer(ss)
	c := captureNumKey(win)
	if c >= 1 && c <= 8 {
		// FIXME do something with the choice here
		td.DrawText(fmt.Sprintf("%d", c), logical.V(13, 2), win)
		offset := sd.WinOffsetV.Add(logical.V(render.CHAR_PIXELS/4, 0))
		sd.WithOffset(offset).DrawSprite(screen.CurrentPlayer.CharacterIcon, logical.V(7, 2), win)
		screen.Players = append(screen.Players, screen.CurrentPlayer)
		screen.CurrentPlayer = Player{}
		if len(screen.Players) == screen.WizardCount {
			// FIXME do something with ComputerDifficulty here
			return NewMainGameScreen(screen.Players)
		} else {
			return &PlayerNameScreen{PlayersScreen: screen.PlayersScreen}
		}
	}
	return screen
}
