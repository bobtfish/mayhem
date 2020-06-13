package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type GameGrid [][]*Character

func MakeGameGrid(x int, y int) *GameGrid {
	gg := make(GameGrid, y)
	for i := 0; i < y; i++ {
		gg[i] = make([]*Character, x)
	}
	return &gg
}

func (grid *GameGrid) PlaceCharacter(x, y int, c *Character) {
	(*grid)[y][x] = c
}

func (grid *GameGrid) GetCharacter(x, y int) *Character {
	return (*grid)[y][x]
}

func drawCharacter(char *Character, x, y int, win *pixelgl.Window) {
	imd := imdraw.New(nil)
	width := 1
	if char == nil {
		width = 0
		imd.Color = pixel.RGB(0, 0, 0)
	} else {
		imd.Color = pickColour()
	}
	imd.Push(pixel.V(float64(x), float64(y)))
	imd.Push(pixel.V(float64(x+CHAR_PIXELS), float64(y+CHAR_PIXELS)))
	imd.Rectangle(float64(width))
	imd.Draw(win)

	if char != nil {
		char.GetText(x, y).Draw(win, pixel.IM)
	}
}

func (grid *GameGrid) Draw(win *pixelgl.Window) {
	xof := CHAR_PIXELS / 2
	yof := WIN_Y - (CHAR_PIXELS*GRID_Y + CHAR_PIXELS/2)
	maxy := len(*grid)
	maxx := len((*grid)[0])
	for x := 0; x < maxx; x++ {
		for y := 0; y < maxy; y++ {
			char := grid.GetCharacter(x, y)
			drawCharacter(char, x*CHAR_PIXELS+xof, y*CHAR_PIXELS+yof, win)
		}
	}
}
