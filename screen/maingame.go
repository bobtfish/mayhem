package screen

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/spells"
)

type StartMainGame struct {
	Players []*player.Player
	Grid    *grid.GameGrid
}

// Initialize stuff for the main game
func (screen *StartMainGame) Enter(ss pixel.Picture, win *pixelgl.Window) {
	for i := 0; i < len(screen.Players); i++ {
		fmt.Printf("%v", screen.Players[i])
		screen.Players[i].Spells = spells.AllSpells
	}
	screen.Grid = grid.MakeGameGrid(GRID_WIDTH, GRID_HEIGHT)
	for i, pos := range player.GetStartPositions(len(screen.Players)) {
		screen.Grid.PlaceGameObject(pos, screen.Players[i])
	}
}

func (screen *StartMainGame) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	return &TurnMenuScreen{
		Players: screen.Players,
		Grid:    screen.Grid,
	}
}
