package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

type PlayersScreen struct {
	ScreenBasics
	WizardCount        int
	ComputerDifficulty int

	origOffset logical.Vec

	Players       []Player
	CurrentPlayer BuildingPlayer
}

type Player struct {
	Name          string
	HumanPlayer   bool
	CharacterIcon logical.Vec
}

type BuildingPlayer struct {
	Player
	NameFinished bool
	AIFinished   bool
	IconChosen   bool
	ColorChosen  bool
}

func (screen *PlayersScreen) Setup(ss pixel.Picture, win *pixelgl.Window) {
	screen.ScreenBasics.Setup(ss, win)
	if screen.Players == nil {
		screen.Players = make([]Player, 0)
	}
	screen.origOffset = screen.SpriteDrawer.WinConverter.Offset
}

func (screen *PlayersScreen) Draw(win *pixelgl.Window) {
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

		screen.TextDrawer.DrawText(fmt.Sprintf("PLAYER %d", playerCount), logical.V(0, 9), win)
		screen.TextDrawer.DrawText("Enter name (12 letters max.)", logical.V(0, 8), win)
		screen.TextDrawer.DrawText("            ", logical.V(0, 7), win)
		screen.TextDrawer.DrawText(screen.CurrentPlayer.Name, logical.V(0, 7), win)

		if win.JustPressed(pixelgl.KeyEnter) && len(screen.CurrentPlayer.Name) > 0 {
			screen.CurrentPlayer.NameFinished = true
		}
	} else {
		if !screen.CurrentPlayer.AIFinished {
			screen.TextDrawer.DrawText("Computer controlled?", logical.V(0, 5), win)
			if win.JustPressed(pixelgl.KeyY) {
				screen.CurrentPlayer.AIFinished = true
				screen.TextDrawer.DrawText("Y", logical.V(21, 5), win)
			}
			if win.JustPressed(pixelgl.KeyN) {
				screen.CurrentPlayer.AIFinished = true
				screen.CurrentPlayer.HumanPlayer = true
				screen.TextDrawer.DrawText("N", logical.V(21, 5), win)
			}
		} else {
			if !screen.CurrentPlayer.IconChosen {
				screen.TextDrawer.DrawText("Which character?", logical.V(0, 4), win)
				screen.TextDrawer.DrawText("1  2  3  4  5  6  7  8", logical.V(0, 3), win)
				for x := 0; x < 8; x++ {
					screen.SpriteDrawer.WinConverter.Offset = logical.V(render.CHAR_PIXELS/4+render.CHAR_PIXELS/2*x*3, render.CHAR_PIXELS*2-render.CHAR_PIXELS/2)
					screen.SpriteDrawer.DrawSprite(logical.V(x, 23), logical.V(1, 3), win)
				}
				c := captureNumKey(win)
				if c >= 1 && c <= 8 {
					screen.CurrentPlayer.CharacterIcon = logical.V(c-1, 23)
					screen.CurrentPlayer.IconChosen = true
					screen.TextDrawer.DrawText(fmt.Sprintf("%d", c), logical.V(17, 4), win)
					screen.SpriteDrawer.WinConverter.Offset = screen.origOffset
					screen.SpriteDrawer.WinConverter.Offset.X = screen.SpriteDrawer.WinConverter.Offset.X + render.CHAR_PIXELS/4
					screen.SpriteDrawer.DrawSprite(screen.CurrentPlayer.CharacterIcon, logical.V(9, 4), win)
				}
			} else {
				if !screen.CurrentPlayer.ColorChosen {
					screen.TextDrawer.DrawText("Which color?", logical.V(0, 2), win)
					screen.TextDrawer.DrawText("1  2  3  4  5  6  7  8", logical.V(0, 1), win)
					for x := 0; x < 8; x++ {
						screen.SpriteDrawer.WinConverter.Offset = logical.V(render.CHAR_PIXELS/4+render.CHAR_PIXELS/2*x*3, render.CHAR_PIXELS*2-render.CHAR_PIXELS/2)
						screen.SpriteDrawer.DrawSprite(screen.CurrentPlayer.CharacterIcon, logical.V(1, 1), win)
					}
					c := captureNumKey(win)
					if c >= 1 && c <= 8 {
						// FIXME do something with the choice here
						screen.CurrentPlayer.ColorChosen = true
						screen.TextDrawer.DrawText(fmt.Sprintf("%d", c), logical.V(13, 2), win)
						screen.SpriteDrawer.WinConverter.Offset = screen.origOffset
						screen.SpriteDrawer.WinConverter.Offset.X = screen.SpriteDrawer.WinConverter.Offset.X + render.CHAR_PIXELS/4
						screen.SpriteDrawer.DrawSprite(screen.CurrentPlayer.CharacterIcon, logical.V(7, 2), win)
					}
				}
			}
		}
	}
}

func (screen *PlayersScreen) NextScreen() GameScreen {
	screen.Players = append(screen.Players, screen.CurrentPlayer.Player)
	screen.CurrentPlayer = BuildingPlayer{}
	if len(screen.Players) == screen.WizardCount {
		return NewMainGameScreen(screen.Players)
	}
	return screen
}

func (screen *PlayersScreen) Finished() bool {
	if screen.CurrentPlayer.ColorChosen == true {
		return true
	}
	return false
}
