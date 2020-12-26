package screen

import (
	"github.com/bobtfish/mayhem/game"
	"github.com/bobtfish/mayhem/player"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

type StartMainGame struct{}

// Initialize stuff for the main game
func (screen *StartMainGame) Enter(ctx screeniface.GameCtx) {
	grid := ctx.GetGrid()
	players := ctx.(*game.Window).GetPlayers()
	for i, pos := range player.GetStartPositions(len(players)) {
		grid.PlaceGameObject(pos, players[i])
	}
}

func (screen *StartMainGame) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	return &TurnMenuScreen{}
}
