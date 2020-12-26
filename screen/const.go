package screen

import (
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

const GridWidth = 15
const GridHeight = 10

const CursorSpell = 0
const CursorBox = 1
const CursorFly = 2
const CursorRangedAttack = 3

func cursorSprite(index int) logical.Vec {
	// All the cursors in the original sprite sheet are
	// are drawn with black on transparent so we need to
	// use the inverse version.
	return logical.V(4+index+render.InverseVideoOffset, 24)
}
