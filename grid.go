package main

import (
	"github.com/faiface/pixel"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

const BLANK_SPRITE_X = 8
const BLANK_SPRITE_Y = 26

type GameGrid [][]*Character

func MakeGameGrid(v logical.Vec) *GameGrid {
	gg := make(GameGrid, v.Y)
	for i := 0; i < v.Y; i++ {
		gg[i] = make([]*Character, v.X)
	}
	return &gg
}

func (grid *GameGrid) PlaceCharacter(v logical.Vec, c *Character) {
	(*grid)[v.Y][v.X] = c
}

func (grid *GameGrid) GetCharacter(v logical.Vec) *Character {
	return (*grid)[v.Y][v.X]
}

func (grid *GameGrid) DrawBatch(sd *render.SpriteDrawer) *pixel.Batch {
	batch := sd.GetNewBatch()
	maxy := len(*grid)
	maxx := len((*grid)[0])
	for x := 0; x < maxx; x++ {
		for y := 0; y < maxy; y++ {
			c := grid.GetCharacter(logical.V(x, y))
			ssX := BLANK_SPRITE_X
			ssY := BLANK_SPRITE_Y
			if c != nil {
				ssX, ssY = c.GetSpriteSheetCoordinates()
			}
			sd.DrawSprite(logical.V(ssX, ssY), logical.V(x, y+1), batch)
		}
	}
	return batch
}

func (grid *GameGrid) AnimationTick() {
	maxy := len(*grid)
	maxx := len((*grid)[0])
	for x := 0; x < maxx; x++ {
		for y := 0; y < maxy; y++ {
			c := grid.GetCharacter(logical.V(x, y))
			if c != nil {
				c.AnimationTick()
			}
		}
	}
}
