package grid

import (
	"image/color"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

/* Empty game object (bottom level tile) */

const BlankSpriteX = 8
const BlankSpriteY = 26

type EmptyObject struct {
	SpriteCoordinates logical.Vec
}

var AnEmptyObject = EmptyObject{
	SpriteCoordinates: logical.V(BlankSpriteX, BlankSpriteY),
}

func (e EmptyObject) AnimationTick(odd bool) {}

func (e EmptyObject) RemoveMe() bool {
	return false
}

func (e EmptyObject) GetColor() color.Color {
	return render.GetColor(0, 0, 0)
}

func (e EmptyObject) Describe() (string, string) {
	return "", ""
}

func (e EmptyObject) IsEmpty() bool {
	return true
}

func (e EmptyObject) GetSpriteSheetCoordinates() logical.Vec {
	return e.SpriteCoordinates
}

// Ignore this
func (e EmptyObject) SetBoardPosition(v logical.Vec) {}
