package logical

import "github.com/faiface/pixel"

type Vec struct {
	X int
	Y int
}

func ZeroVec() Vec     { return Vec{0, 0} }
func IdentityVec() Vec { return Vec{1, 1} }

func V(x, y int) Vec {
	return Vec{X: x, Y: y}
}

func (v Vec) Equals(w Vec) bool {
	return v.X == w.X && v.Y == w.Y
}

func (v Vec) Add(w Vec) Vec {
	return Vec{v.X + w.X, v.Y + w.Y}
}

func (v Vec) Subtract(w Vec) Vec {
	return Vec{v.X - w.X, v.Y - w.Y}
}

func (v Vec) Multiply(w Vec) Vec {
	return Vec{v.X * w.X, v.Y * w.Y}
}

func (v Vec) ToPixelVec() pixel.Vec {
	return pixel.V(float64(v.X), float64(v.Y))
}

func (v Vec) ToPixelRect(scale Vec) pixel.Rect {
    return pixel.Rect{
        Min: v.Multiply(scale).ToPixelVec(),
        Max: v.Multiply(scale).Add(scale).ToPixelVec(),
    }
}

func (v Vec) ToPixelRectOffset(scale Vec, offset Vec) pixel.Rect {
    return pixel.Rect{
        Min: v.Multiply(scale).Add(offset).ToPixelVec(),
        Max: v.Multiply(scale).Add(scale).Add(offset).ToPixelVec(),
    }
}
