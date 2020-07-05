package grid

import (
	"time"

	"github.com/faiface/pixel"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

type GameGrid [][]*GameObjectStack

func MakeGameGrid(width, height int) *GameGrid {
	gg := make(GameGrid, height)
	for i := 0; i < height; i++ {
		gg[i] = make([]*GameObjectStack, width)
		for j := 0; j < width; j++ {
			gg[i][j] = NewGameObjectStack()
		}
	}
	grid := &gg

	go func() {
		Qsecond := time.Tick(time.Second / 4)
		for true == true {
			select {
			case <-Qsecond:
				grid.AnimationTick()
			}
		}
	}()

	return grid
}

func (grid *GameGrid) PlaceGameObject(v logical.Vec, c GameObjectStackable) {
	(*grid)[v.Y][v.X].PlaceObject(c)
}

func (grid *GameGrid) GetGameObject(v logical.Vec) GameObject {
	return (*grid)[v.Y][v.X]
}

func (grid *GameGrid) DrawBatch(sd render.SpriteDrawer) *pixel.Batch {
	batch := sd.GetNewBatch()
	for x := 0; x <= grid.MaxX(); x++ {
		for y := 0; y <= grid.MaxY(); y++ {
			c := grid.GetGameObject(logical.V(x, y))
			v := c.GetSpriteSheetCoordinates()
			sd.DrawSprite(v, logical.V(x, y), batch)
		}
	}
	return batch
}

func (grid *GameGrid) Height() int {
	return len(*grid)
}

func (grid *GameGrid) MaxY() int {
	return len(*grid) - 1
}

func (grid *GameGrid) Width() int {
	return len((*grid)[0])
}

func (grid *GameGrid) MaxX() int {
	return len((*grid)[0]) - 1
}

func (grid *GameGrid) AnimationTick() {
	for x := 0; x < grid.MaxX(); x++ {
		for y := 0; y < grid.MaxY(); y++ {
			c := grid.GetGameObject(logical.V(x, y))
			if c != nil {
				c.AnimationTick()
			}
		}
	}
}

func (grid *GameGrid) AsRect() logical.Rect {
	return logical.R(grid.MaxX(), grid.MaxY())
}
