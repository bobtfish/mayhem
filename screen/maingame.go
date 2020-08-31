package screen

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/player"
)

type StartMainGame struct {
	Players []*player.Player
	Grid    *grid.GameGrid
}

// Initialize stuff for the main game
func (screen *StartMainGame) Enter(ss pixel.Picture, win *pixelgl.Window) {
	screen.Grid = grid.MakeGameGrid(GridWidth, GridHeight)
	for i, pos := range player.GetStartPositions(len(screen.Players)) {
		screen.Grid.PlaceGameObject(pos, screen.Players[i])
	}
}

func (screen *StartMainGame) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	return &TurnMenuScreen{
		Players:   screen.Players,
		Grid:      screen.Grid,
		LawRating: 2,
	}
}
