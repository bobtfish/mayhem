package grid

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

type GameGrid [][]*GameObjectStack

func MakeGameGrid(v logical.Vec) *GameGrid {
	fmt.Printf("Make grid X%d Y%d\n", v.X, v.Y)
	gg := make(GameGrid, v.Y)
	for i := 0; i < v.Y; i++ {
		gg[i] = make([]*GameObjectStack, v.X)
		for j := 0; j < v.X; j++ {
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
	for x := 0; x < grid.Width(); x++ {
		for y := 0; y < grid.Height(); y++ {
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

func (grid *GameGrid) Width() int {
	return len((*grid)[0])
}

func (grid *GameGrid) AnimationTick() {
	for x := 0; x < grid.Width(); x++ {
		for y := 0; y < grid.Height(); y++ {
			fmt.Printf("Animation tick for grid X%d Y%d\n", x, y)
			c := grid.GetGameObject(logical.V(x, y))
			if c != nil {
				c.AnimationTick()
			}
		}
	}
}

func (grid *GameGrid) AsRect() logical.Rect {
	return logical.R(grid.Width(), grid.Height())
}
