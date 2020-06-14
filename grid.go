package main

import (
	"github.com/faiface/pixel"

	"github.com/bobtfish/mayhem/render"
)

const BLANK_SPRITE_X = 8
const BLANK_SPRITE_Y = 26

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

func (grid *GameGrid) DrawBatch(sd *render.SpriteDrawer) *pixel.Batch {
	batch := sd.GetNewBatch()
	maxy := len(*grid)
	maxx := len((*grid)[0])
	for x := 0; x < maxx; x++ {
		for y := 0; y < maxy; y++ {
			c := grid.GetCharacter(x, y)
			ssX := BLANK_SPRITE_X
			ssY := BLANK_SPRITE_Y
			if c != nil {
				ssX, ssY = c.GetSpriteSheetCoordinates()
			}
			sd.DrawSprite(ssX, ssY, x, y+1, batch)
		}
	}
	return batch
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
