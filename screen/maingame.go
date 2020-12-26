package screen

import (
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/player"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

type StartMainGame struct {
	Players []*player.Player
	Grid    *grid.GameGrid
}

// Initialize stuff for the main game
func (screen *StartMainGame) Enter(ctx screeniface.GameCtx) {
	screen.Grid = grid.MakeGameGrid(GridWidth, GridHeight)
	for i, pos := range player.GetStartPositions(len(screen.Players)) {
		screen.Grid.PlaceGameObject(pos, screen.Players[i])
	}
}

func (screen *StartMainGame) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	return &TurnMenuScreen{
		Players: screen.Players,
		Grid:    screen.Grid,
	}
}
