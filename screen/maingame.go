package screen

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/spells"
)

type MainGameScreen struct {
	Players       []Player
	CurrentScreen GameScreen
}

func NewMainGameScreen(players []Player) *MainGameScreen {
	for i := 0; i < len(players); i++ {
		players[i].Spells = spells.AllSpells
	}
	var refPlayers = make([]*Player, len(players))
	turnmenu := &TurnMenuScreen{
		Players: refPlayers,
	}
	main := MainGameScreen{
		Players:       players,
		CurrentScreen: turnmenu,
	}
	for i := 0; i < len(players); i++ {
		turnmenu.Players[i] = &main.Players[i]
	}
	return &main
}

func (screen *MainGameScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	screen.CurrentScreen.Enter(ss, win)
}

func (screen *MainGameScreen) Step(win *pixelgl.Window) GameScreen {
	return screen.CurrentScreen.Step(win)
}
