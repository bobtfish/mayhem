package fx

import (
	"fmt"
	"image/color"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

// Special effects

type Fx struct {
	SpriteVec   logical.Vec
	SpriteCount int
	SpriteIdx   int
	Color       color.Color
	RepeatIdx   int
	RepeatCount int
}

// GameObject interface START
func (c *Fx) AnimationTick(odd bool) {
	c.SpriteIdx++
	if c.SpriteIdx == c.SpriteCount && c.RepeatIdx < c.RepeatCount {
		c.RepeatIdx++
		c.SpriteIdx = 0
	}
}

func (c *Fx) RemoveMe() bool {
	if c.SpriteIdx == c.SpriteCount {
		return true
	}
	return false
}

func (c *Fx) IsEmpty() bool {
	return false
}

func (c *Fx) GetSpriteSheetCoordinates() logical.Vec {
	return logical.V(c.SpriteVec.X+c.SpriteIdx, c.SpriteVec.Y)
}

func (f *Fx) GetColor() color.Color {
	return f.Color
}

func (f *Fx) Describe() string {
	return fmt.Sprintf("Fx:%T", *f)
}

func (f *Fx) SetBoardPosition(v logical.Vec) {}

// GameObject interface END

// Fx Constructors

func FxSpellCast() *Fx {
	return &Fx{
		SpriteVec:   logical.V(0, 28),
		SpriteCount: 8,
		Color:       render.GetColor(0, 255, 255),
	}
}

func FxRemoteAttack() *Fx {
	return &Fx{
		SpriteVec:   logical.V(0, 27),
		SpriteCount: 8,
		Color:       render.GetColor(255, 255, 255),
	}
}

func FxFire() *Fx {
	return &Fx{
		SpriteVec:   logical.V(0, 26),
		SpriteCount: 8,
		Color:       render.GetColor(255, 255, 255),
	}
}

func FxDisbelieve() *Fx {
	return &Fx{
		SpriteVec:   logical.V(0, 25),
		SpriteCount: 7,
		Color:       render.GetColor(255, 255, 255),
	}
}

func FxAttack() *Fx {
	return &Fx{
		SpriteVec:   logical.V(0, 24),
		SpriteCount: 4,
		Color:       render.GetColor(255, 255, 255),
		RepeatCount: 3,
	}
}
