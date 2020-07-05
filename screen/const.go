package screen

import (
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

const WIN_X = 1024
const WIN_Y = 768
const GRID_WIDTH = 15
const GRID_HEIGHT = 10

const CURSOR_SPELL = 0
const CURSOR_BOX   = 1
const CURSOR_FLY   = 2
const CURSOR_BOOM  = 3

func cursorSprite(index int) logical.Vec {
	// All the cursors in the original sprite sheet are
	// are drawn with black on transparent so we need to
	// use the inverse version.
	return logical.V(4+index+render.INVERSE_VIDEO_OFFSET, 24)
}
