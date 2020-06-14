package logical

type Vec {
    X int
    Y int
}

func Vec(x, y int) Vec {
    return Vec{X: x, Y: y}
}

