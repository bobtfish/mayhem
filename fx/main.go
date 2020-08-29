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
func (f *Fx) AnimationTick(odd bool) {
	f.SpriteIdx++
	if f.SpriteIdx == f.SpriteCount && f.RepeatIdx < f.RepeatCount {
		f.RepeatIdx++
		f.SpriteIdx = 0
	}
}

func (f *Fx) RemoveMe() bool {
	return f.SpriteIdx == f.SpriteCount
}

func (f *Fx) IsEmpty() bool {
	return false
}

func (f *Fx) GetSpriteSheetCoordinates() logical.Vec {
	return logical.V(f.SpriteVec.X+f.SpriteIdx, f.SpriteVec.Y)
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
