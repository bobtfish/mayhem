package logical

import "github.com/faiface/pixel"

type VecConverter struct {
	Offset     Vec
	Multiplier int
}

func (c VecConverter) ToPixelVec(v Vec) pixel.Vec {
	return pixel.V(float64(v.X*c.Multiplier+c.Offset.X), float64(v.Y*c.Multiplier+c.Offset.Y))
}
