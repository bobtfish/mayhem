package logical

import (
	"math"

	"github.com/faiface/pixel"
)

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

func (v Vec) Abs() Vec {
	return Vec{abs(v.X), abs(v.Y)}
}

func (v Vec) smallestSquare() int {
	w := v.Abs()
	if w.X < w.Y {
		return w.X
	}
	return w.Y
}

// Distance of diagonal - first square is 1, second is 3, third is 4 etc
func squareDistance(i int) int {
	if i == 0 {
		return 0
	}
	return i + int(math.Floor(float64(i)/2))
}

// D&D distance, because of course this is what the original game does...
func (v Vec) Distance(w Vec) int {
	x := v.Subtract(w).Abs()
	ss := x.smallestSquare()
	y := x.Subtract(V(ss, ss)) // 1 is has 0 X or 0 Y so just add X and Y 'remainder'
	return squareDistance(ss) + y.X + y.Y
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
