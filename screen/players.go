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
	ScreenBasics
	WizardCount        int
	ComputerDifficulty int

	Players       []Player
	CurrentPlayer BuildingPlayer
}

type Player struct {
	Name          string
	HumanPlayer   bool
	CharacterIcon logical.Vec

	Spells      []spells.Spell
	ChosenSpell int
}

type BuildingPlayer struct {
	Player
	NameFinished bool
	AIFinished   bool
	IconChosen   bool
	ColorChosen  bool
}

func (screen *PlayersScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	screen.ScreenBasics.Enter(ss, win)
	if screen.Players == nil {
		screen.Players = make([]Player, 0)
	}
}

func (screen *PlayersScreen) Step(win *pixelgl.Window) GameScreen {
	td := TextDrawer(screen.SpriteSheet)
	sd := SpriteDrawer(screen.SpriteSheet)
	if !screen.CurrentPlayer.NameFinished {
		playerCount := len(screen.Players) + 1

		if win.JustPressed(pixelgl.KeyEnter) && len(screen.CurrentPlayer.Name) >= 0 {
			screen.CurrentPlayer.NameFinished = true
		} else {
			if win.JustPressed(pixelgl.KeyBackspace) && len(screen.CurrentPlayer.Name) >= 0 {
				length := len(screen.CurrentPlayer.Name) - 1
				screen.CurrentPlayer.Name = screen.CurrentPlayer.Name[:length]
			} else {
				screen.CurrentPlayer.Name = fmt.Sprintf("%s%s", screen.CurrentPlayer.Name, win.Typed())
			}
		}

		if len(screen.CurrentPlayer.Name) > 12 {
			screen.CurrentPlayer.Name = screen.CurrentPlayer.Name[:12]
		}

		td.DrawText(fmt.Sprintf("PLAYER %d", playerCount), logical.V(0, 9), win)
		td.DrawText("Enter name (12 letters max.)", logical.V(0, 8), win)
		td.DrawText("            ", logical.V(0, 7), win)
		td.DrawText(screen.CurrentPlayer.Name, logical.V(0, 7), win)

		if win.JustPressed(pixelgl.KeyEnter) && len(screen.CurrentPlayer.Name) > 0 {
			screen.CurrentPlayer.NameFinished = true
		}
	} else {
		if !screen.CurrentPlayer.AIFinished {
			td.DrawText("Computer controlled?", logical.V(0, 5), win)
			if win.JustPressed(pixelgl.KeyY) {
				screen.CurrentPlayer.AIFinished = true
				td.DrawText("Y", logical.V(21, 5), win)
			}
			if win.JustPressed(pixelgl.KeyN) {
				screen.CurrentPlayer.AIFinished = true
				screen.CurrentPlayer.HumanPlayer = true
				td.DrawText("N", logical.V(21, 5), win)
			}
		} else {
			if !screen.CurrentPlayer.IconChosen {
				td.DrawText("Which character?", logical.V(0, 4), win)
				td.DrawText("1  2  3  4  5  6  7  8", logical.V(0, 3), win)
				for x := 0; x < 8; x++ {
					offset := logical.V(render.CHAR_PIXELS/4+render.CHAR_PIXELS/2*x*3, render.CHAR_PIXELS*2-render.CHAR_PIXELS/2)
					sd.WithOffset(offset).DrawSprite(logical.V(x, 23), logical.V(1, 3), win)
				}
				c := captureNumKey(win)
				if c >= 1 && c <= 8 {
					screen.CurrentPlayer.CharacterIcon = logical.V(c-1, 23)
					screen.CurrentPlayer.IconChosen = true
					td.DrawText(fmt.Sprintf("%d", c), logical.V(17, 4), win)
					offset := sd.WinOffsetV.Add(logical.V(render.CHAR_PIXELS/4, 0))
					sd.WithOffset(offset).DrawSprite(screen.CurrentPlayer.CharacterIcon, logical.V(9, 4), win)
				}
			} else {
				if !screen.CurrentPlayer.ColorChosen {
					td.DrawText("Which color?", logical.V(0, 2), win)
					td.DrawText("1  2  3  4  5  6  7  8", logical.V(0, 1), win)
					for x := 0; x < 8; x++ {
						offset := logical.V(render.CHAR_PIXELS/4+render.CHAR_PIXELS/2*x*3, render.CHAR_PIXELS*2-render.CHAR_PIXELS/2)
						sd.WithOffset(offset).DrawSprite(screen.CurrentPlayer.CharacterIcon, logical.V(1, 1), win)
					}
					c := captureNumKey(win)
					if c >= 1 && c <= 8 {
						// FIXME do something with the choice here
						screen.CurrentPlayer.ColorChosen = true
						td.DrawText(fmt.Sprintf("%d", c), logical.V(13, 2), win)
						offset := sd.WinOffsetV.Add(logical.V(render.CHAR_PIXELS/4, 0))
						sd.WithOffset(offset).DrawSprite(screen.CurrentPlayer.CharacterIcon, logical.V(7, 2), win)
					}
				}
			}
		}
	}
	if screen.CurrentPlayer.ColorChosen == true {
		screen.Players = append(screen.Players, screen.CurrentPlayer.Player)
		screen.CurrentPlayer = BuildingPlayer{}
		screen.CurrentPlayer.ChosenSpell = -1
		if len(screen.Players) == screen.WizardCount {
			return NewMainGameScreen(screen.Players)
		}
	}
	return screen
}
