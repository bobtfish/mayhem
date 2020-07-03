package logical

type Vec struct {
	X int
	Y int
}

func V(x, y int) Vec {
	return Vec{X: x, Y: y}
}

func (v Vec) Add(w Vec) Vec {
    return Vec{v.X+w.X, v.Y+w.Y}
}

func (v Vec) Subtract(w Vec) Vec {
    return Vec{v.X-w.X, v.Y-w.Y}
}
