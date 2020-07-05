package screen

import (
	"github.com/bobtfish/mayhem/logical"
)

const WIN_X = 1024
const WIN_Y = 768
const GRID_WIDTH = 15
const GRID_HEIGHT = 10

func cursorSprite() logical.Vec {
	return logical.V(4, 24)
}
