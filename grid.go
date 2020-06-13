package main

import (
	"github.com/faiface/pixel/pixelgl"
)

type GameGrid [][]Character

func MakeGameGrid(x int, y int) *GameGrid {
	gg := make(GameGrid, y)
	for i := 0; i < y; i++ {
		gg[i] = make([]Character, x)
	}
	return &gg
}

func (grid *GameGrid) PlaceCharacter(x, y int, c Character) {
	(*grid)[x][y] = c
}

func (grid *GameGrid) Draw(xof, yof float64, win *pixelgl.Window) {
	maxy := len(*grid)
	maxx := len((*grid)[0])
	for x := 0; x < maxx; x++ {
		for y := 0; y < maxy; y++ {
			getRectangle(float64(x*CHAR_PIXELS)+xof, float64(y*CHAR_PIXELS)+yof, float64(CHAR_PIXELS)).Draw(win)
		}
	}
}
