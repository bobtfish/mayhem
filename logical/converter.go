package logical

import "github.com/faiface/pixel"

type VecConverter struct {
	Offset     Vec
	Multiplier Vec
}

func NewVecConverter(offset Vec, multipler Vec) VecConverter {
	return VecConverter{offset, multipler}
}

func (c VecConverter) ToPixelVec(v Vec) pixel.Vec {
	return v.Multiply(c.Multiplier).Add(c.Offset).ToPixelVec()
}
