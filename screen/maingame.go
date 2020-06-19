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

func (screen *MainGameScreen) Setup(ss pixel.Picture, win *pixelgl.Window) {
	screen.CurrentScreen.Setup(ss, win)
}

func (screen *MainGameScreen) Draw(win *pixelgl.Window) {
	screen.CurrentScreen.Draw(win)
}

func (screen *MainGameScreen) NextScreen() GameScreen {
	screen.CurrentScreen = screen.CurrentScreen.NextScreen()
	return screen
}

func (screen *MainGameScreen) Finished() bool {
	return screen.CurrentScreen.Finished()
}
