package logical

type Rect struct {
	Vec
}

func R(x, y int) Rect {
	return Rect{Vec: V(x, y)}
}

func (r Rect) Contains(v Vec) bool {
	if v.X < 0 || v.Y < 0 {
		return false
	}
	if v.X > r.Vec.X || v.Y > r.Vec.Y {
		return false
	}
	return true
}

func (r Rect) Adjacents(v Vec) []Vec {
	vecs := make([]Vec, 0)
	if r.Contains(v) {
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				if x == 0 && y == 0 {
					continue
				}
				newV := V(v.X+x, v.Y+y)
				if r.Contains(newV) {
					vecs = append(vecs, newV)
				}
			}
		}
	}
	return vecs
}

func (r Rect) Width() int {
	return r.Vec.X
}

func (r Rect) Height() int {
	return r.Vec.Y
}

func (r Rect) Clamp(v Vec) Vec {
	if v.X > r.Vec.X {
		v.X = r.Vec.X
	}
	if v.X < 0 {
		v.X = 0
	}
	if v.Y > r.Vec.Y {
		v.Y = r.Vec.Y
	}
	if v.Y < 0 {
		v.Y = 0
	}
	return v
}
