package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
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
		atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(pixel.V(float64(x+2), float64(y+2)), atlas)
		fmt.Fprintln(basicTxt, char.Name)
		basicTxt.Draw(win, pixel.IM)
	}
}

func (grid *GameGrid) Draw(xof, yof int, win *pixelgl.Window) {
	maxy := len(*grid)
	maxx := len((*grid)[0])
	for x := 0; x < maxx; x++ {
		for y := 0; y < maxy; y++ {
			char := grid.GetCharacter(x, y)
			drawCharacter(char, x*CHAR_PIXELS+xof, y*CHAR_PIXELS+yof, win)
		}
	}
}
