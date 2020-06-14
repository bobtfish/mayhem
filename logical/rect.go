package logical

type Rect struct {
	Vec
}

func Rect(x, y int) Rect {
	return Rect{Vec: Vec(x, y)}
}

func (r Rect) Contains(v Vec) bool {
	if v.X < 0 || v.Y < 0 {
		return false
	}
	if v.X > r.Vec.X || v.Y > r.Vec.Y {
		return false
	}
}

func (r Rect) Adjacents(v Vec) []Vec {
	vecs := make([]Vec, 0)
	for x := -1; x <= 1; x++ {
		for y := -1; x <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			newV := Vec(v.X+x, v.Y+y)
			if r.Contains(newV) {
				vecs = append(vecs, newV)
			}
		}
	}
	return adjacents
}
