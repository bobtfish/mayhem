package main

import (
	"github.com/bobtfish/mayhem/logical"
)

// Special effects

type Fx struct {
	SpriteVec   logical.Vec
	SpriteCount int
	SpriteIdx   int
}

// GameObject interface START
func (c *Fx) AnimationTick() {
	c.SpriteIdx++
}

func (c *Fx) RemoveMe() bool {
	if c.SpriteIdx == c.SpriteCount {
		return true
	}
	return false
}

func (c *Fx) GetSpriteSheetCoordinates() logical.Vec {
	return logical.V(c.SpriteVec.X+c.SpriteIdx, c.SpriteVec.Y)
}

// GameObject interface END

// Fx Constructors

func FxWarp() *Fx {
	return &Fx{
		SpriteVec:   logical.V(0, 28),
		SpriteCount: 8,
	}
}

func FxBlam() *Fx {
	return &Fx{
		SpriteVec:   logical.V(0, 27),
		SpriteCount: 8,
	}
}

func FxFire() *Fx {
	return &Fx{
		SpriteVec:   logical.V(0, 26),
		SpriteCount: 8,
	}
}

func FxBoom() *Fx {
	return &Fx{
		SpriteVec:   logical.V(0, 25),
		SpriteCount: 7,
	}
}

func FxPop() *Fx {
	return &Fx{
		SpriteVec:   logical.V(0, 24),
		SpriteCount: 4,
	}
}