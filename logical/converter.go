package logical

import "github.com/faiface/pixel"

type VecConverter interface {
    ToPixelVec(Vec) pixel.Vec
}

type OffsetVecConverter struct {
    Offset Vec
    Multipler int
}

func (c OffsetVecConverter) ToPixelVec(v Vec) pixel.Vec {
    return pixel.V(float64(v.X*c.Multipler+c.Offset.X), float64(v.Y*c.Multipler+c.Offset.Y))
}

