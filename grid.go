package main

import (
	"github.com/faiface/pixel"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

type GameGrid [][]*GameObjectStack

func MakeGameGrid(v logical.Vec) *GameGrid {
	gg := make(GameGrid, v.Y)
	for i := 0; i < v.Y; i++ {
		gg[i] = make([]*GameObjectStack, v.X)
		for j := 0; j < v.X; j++ {
			gg[i][j] = NewGameObjectStack()
		}
	}
	return &gg
}

func (grid *GameGrid) PlaceGameObject(v logical.Vec, c GameObjectStackable) {
	(*grid)[v.Y][v.X].PlaceObject(c)
}

func (grid *GameGrid) GetGameObject(v logical.Vec) GameObject {
	return (*grid)[v.Y][v.X]
}

func (grid *GameGrid) DrawBatch(sd render.SpriteDrawer) *pixel.Batch {
	batch := sd.GetNewBatch()
	maxy := len(*grid)
	maxx := len((*grid)[0])
	for x := 0; x < maxx; x++ {
		for y := 0; y < maxy; y++ {
			c := grid.GetGameObject(logical.V(x, y))
			v := c.GetSpriteSheetCoordinates()
			sd.DrawSprite(v, logical.V(x, y), batch)
		}
	}
	return batch
}

func (grid *GameGrid) AnimationTick() {
	maxy := len(*grid)
	maxx := len((*grid)[0])
	for x := 0; x < maxx; x++ {
		for y := 0; y < maxy; y++ {
			c := grid.GetGameObject(logical.V(x, y))
			if c != nil {
				c.AnimationTick()
			}
		}
	}
}
