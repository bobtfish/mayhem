package screen

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type MainGameScreen struct {
	Players       []Player
	CurrentScreen GameScreen
}

func NewMainGameScreen(players []Player) *MainGameScreen {
	var refPlayers = make([]*Player, len(players))
	for i := 0; i < len(players); i++ {
		refPlayers[i] = &players[i]
	}
	main := MainGameScreen{
		Players: players,
		CurrentScreen: &TurnMenuScreen{
			Players: refPlayers,
		},
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
