package logical

type Vec struct {
	X int
	Y int
}

func V(x, y int) Vec {
	return Vec{X: x, Y: y}
}
