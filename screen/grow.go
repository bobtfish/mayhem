package screen

import (
	"fmt"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/character"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
)

const GROW_CHANCE = 10
const VANISH_CHANCE = 2

var growable map[string]func()

func init() {
	growable = map[string]func(){
		"Gooey Blob": func() {},
		"Fire":       func() {},
	}
}

func doesItGrow() bool {
	return rand.Intn(100) <= GROW_CHANCE
}

func doesItVanish() bool {
	return rand.Intn(100) <= VANISH_CHANCE
}

type GrowScreen struct {
	*WithBoard
	Consider logical.Vec
}

func (screen *GrowScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
}

func (screen *GrowScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	batch.Draw(win)

	nextItem := screen.FindNextItem()
	if nextItem != nil {
	}

	return &Pause{
		Grid: screen.WithBoard.Grid,
		NextScreen: &TurnMenuScreen{
			Players: screen.WithBoard.Players,
			Grid:    screen.WithBoard.Grid,
		},
	}
}

func (screen *GrowScreen) FindNextItem() grid.GameObject {
	for screen.Consider.X < screen.WithBoard.Grid.MaxX() && screen.Consider.Y < screen.WithBoard.Grid.MaxY() {
		character, ok := screen.WithBoard.Grid.GetGameObject(screen.Consider).(*character.Character)
		if ok {
			for name, actor := range growable {
				if name == character.Name {
					fmt.Printf("v(%d, %d) has %s\n", screen.Consider.X, screen.Consider.Y, character.Name)
				}
			}
		}

		screen.Consider.X = screen.Consider.X + 1
		if screen.Consider.X == screen.WithBoard.Grid.MaxX() {
			screen.Consider.X = 0
			screen.Consider.Y = screen.Consider.Y + 1
		}
	}
	return nil
}
