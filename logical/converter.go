package logical

import "github.com/faiface/pixel"

type VecConverter struct {
	Offset      Vec
	Multiplier  Vec
}

func NewVecConverter(offset Vec, multipler Vec) VecConverter {
	return VecConverter{offset, multipler}
}

func (c VecConverter) ToPixelVec(v Vec) pixel.Vec {
	return pixel.V(float64(v.X*c.Multiplier.X+c.Offset.X), float64(v.Y*c.Multiplier.Y+c.Offset.Y))
}
