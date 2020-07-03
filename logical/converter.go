package logical

import "github.com/faiface/pixel"

type VecConverter struct {
	Offset      Vec
	XMultiplier int
	YMultiplier int
}

func NewVecConverter(offset Vec, mulX int, mulY int) VecConverter {
    return VecConverter{offset, mulX, mulY}
}

func (c VecConverter) ToPixelVec(v Vec) pixel.Vec {
	return pixel.V(float64(v.X*c.XMultiplier+c.Offset.X), float64(v.Y*c.YMultiplier+c.Offset.Y))
}
