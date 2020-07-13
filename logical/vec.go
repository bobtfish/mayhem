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

// FIXME this is not how the original game considers distance
//       as V(1,1) is only 1 square away
func (v Vec) Distance(w Vec) int {
	x := v.Subtract(w)
	return math.Abs(x.X) + math.Abs(x.Y)
}

func (v Vec) ToPixelVec() pixel.Vec {
	return pixel.V(float64(v.X), float64(v.Y))
}

func (v Vec) ToPixelRect(scale Vec, offsets ...Vec) pixel.Rect {
	min := v.Multiply(scale)
	for _, offset := range offsets {
		min = min.Add(offset)
	}
	return pixel.Rect{
		Min: min.ToPixelVec(),
		Max: min.Add(scale).ToPixelVec(),
	}
}
