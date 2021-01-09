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
		ticker := time.Tick(time.Second / 8)
		count := 0
		for true == true {
			select {
			case <-ticker:
				count++
				grid.AnimationTick(count%2 != 0)
			}
		}
	}()

	return grid
}

func (grid *GameGrid) PlaceGameObject(v logical.Vec, c GameObjectStackable) {
	c.SetBoardPosition(v)
	(*grid)[v.Y][v.X].PlaceObject(c)
}

func (grid *GameGrid) GetGameObjectStack(v logical.Vec) *GameObjectStack {
	return (*grid)[v.Y][v.X]
}

func (grid *GameGrid) GetGameObject(v logical.Vec) GameObject {
	return grid.GetGameObjectStack(v).TopObject()
}

func (grid *GameGrid) DrawBatch(sd render.SpriteDrawer) *pixel.Batch {
	batch := sd.GetNewBatch()
	for x := 0; x <= grid.MaxX(); x++ {
		for y := 0; y <= grid.MaxY(); y++ {
			c := grid.GetGameObjectStack(logical.V(x, y))
			v := c.GetSpriteSheetCoordinates()
			sd.DrawSpriteColor(v, logical.V(x, y), c.GetColor(), batch)
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

func (grid *GameGrid) AnimationTick(odd bool) {
	for x := 0; x <= grid.MaxX(); x++ {
		for y := 0; y <= grid.MaxY(); y++ {
			grid.GetGameObjectStack(logical.V(x, y)).AnimationTick(odd)
		}
	}
}

func (grid *GameGrid) AsRect() logical.Rect {
	return logical.R(grid.MaxX(), grid.MaxY())
}

func (grid *GameGrid) HaveLineOfSight(from, to logical.Vec) bool {
	for _, pathV := range to.Subtract(from).Path() {
		if !grid.GetGameObject(from.Add(pathV)).IsEmpty() {
			return false
		}
	}
	return true
}
