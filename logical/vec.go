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

func (v Vec) IsDiagonal() bool {
	a := v.Abs()
	if a.X > 0 && a.Y > 0 {
		return true
	}
	return false
}

func (v Vec) Path() []Vec {
	var Xstep, Ystep, Xcurrent, Ycurrent float64
	Xsign, Ysign := 1, 1
	if v.X < 0 {
		Xsign = -1
	}
	if v.Y < 0 {
		Ysign = -1
	}
	w := v.Abs()
	if w.X == w.Y {
		Xstep = 1.0
		Ystep = 1.0
	} else {
		if w.X > w.Y {
			Xstep = 1.0
			Ystep = float64(w.Y) / float64(w.X)
		} else {
			Ystep = 1.0
			Xstep = float64(w.X) / float64(w.Y)
		}
	}
	Xcurrent, Ycurrent = Xstep, Ystep
	path := make([]Vec, 0)
	for Xcurrent < float64(w.X) || Ycurrent < float64(w.Y) {
		path = append(path, V(int(Xcurrent)*Xsign, int(Ycurrent)*Ysign))
		Xcurrent = Xcurrent + Xstep
		Ycurrent = Ycurrent + Ystep
	}
	return path
}
