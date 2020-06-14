package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const CHAR_PIXELS = 64
const WIN_Y = 768

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

func drawCharacter(char *Character, x, y int, win pixel.Target, ss pixel.Picture) {
	imd := imdraw.New(nil)
	if char == nil {
		//  return
		imd.Color = pixel.RGB(0, 0, 0)
		imd.Push(pixel.V(float64(x), float64(y)))
		imd.Push(pixel.V(float64(x+CHAR_PIXELS), float64(y+CHAR_PIXELS)))
		imd.Rectangle(0)
		imd.Draw(win)
		return
	}
	sprite := char.GetSprite(ss)
	if sprite == nil {
		imd.Color = pickColour()
		imd.Push(pixel.V(float64(x), float64(y)))
		imd.Push(pixel.V(float64(x+CHAR_PIXELS), float64(y+CHAR_PIXELS)))
		imd.Rectangle(1)
		imd.Draw(win)
		char.GetText(x, y).Draw(win, pixel.IM)
	} else {
		mat := pixel.IM
		v := pixel.V(float64(x), float64(y))
		mat = mat.Moved(v)
		mat = mat.ScaledXY(v, pixel.V(4, 4))
		mat = mat.Moved(pixel.V(31, 31))
		sprite.DrawColorMask(win, mat, char.GetColorMask())
	}
}

func (grid *GameGrid) Draw(win *pixelgl.Window, ss pixel.Picture) {
	xof := CHAR_PIXELS / 2
	yof := WIN_Y - (CHAR_PIXELS*GRID_Y + CHAR_PIXELS/2)
	maxy := len(*grid)
	maxx := len((*grid)[0])

	batch := pixel.NewBatch(&pixel.TrianglesData{}, ss)
	batch.Clear()
	for x := 0; x < maxx; x++ {
		for y := 0; y < maxy; y++ {
			char := grid.GetCharacter(x, y)
			drawCharacter(char, x*CHAR_PIXELS+xof, y*CHAR_PIXELS+yof, batch, ss)
		}
	}
	batch.Draw(win)
}

func (grid *GameGrid) AnimationTick() {
	maxy := len(*grid)
	maxx := len((*grid)[0])
	for x := 0; x < maxx; x++ {
		for y := 0; y < maxy; y++ {
			c := grid.GetCharacter(x, y)
			if c != nil {
				c.AnimationTick()
			}
		}
	}
}
