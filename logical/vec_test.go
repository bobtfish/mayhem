package logical

import "testing"

func TestVec(t *testing.T) {
	v := V(0, 0)
	if v.X != 0 {
		t.Errorf("Zero vec not zero X")
	}
	if v.Y != 0 {
		t.Errorf("Zero vec not zero Y")
	}
}

func TestVecEquals(t *testing.T) {
	if !V(0, 0).Equals(ZeroVec()) {
		t.Errorf("V(0, 0) not equal zero Vec")
	}
	if V(1, 0).Equals(ZeroVec()) {
		t.Errorf("V(1, 0) is equal zero Vec")
	}
	if V(0, 1).Equals(ZeroVec()) {
		t.Errorf("V(0, 1) is equal zero Vec")
	}
}

func TestVecAdd(t *testing.T) {
	v := V(1, 2).Add(V(1, 2))
	if v.X != 2 {
		t.Errorf("v(1, 2) + v(1, 2) x != 2")
	}
	if v.Y != 4 {
		t.Errorf("v(1, 2) + v(1, 2) y != 4")
	}
}

func TestVecSubtract(t *testing.T) {
	v := V(4, 8).Subtract(V(1, 2))
	if v.X != 3 {
		t.Errorf("v(4, 8) - v(1, 2) x != 3 is %d", v.X)
	}
	if v.Y != 6 {
		t.Errorf("v(4, 8) + v(1, 2) y != 6 is %d", v.Y)
	}
}

func TestVecMultiply(t *testing.T) {
	v := V(2, 3).Multiply(V(2, 2))
	if v.X != 4 {
		t.Errorf("v(2, 3) * v(2, 2) x != 4 is %d", v.X)
	}
	if v.Y != 6 {
		t.Errorf("v(2, 3) * v(2, 2) y != 6 is %d", v.Y)
	}
}

func TestVecIdentity(t *testing.T) {
	v := V(2, 3).Multiply(IdentityVec())
	if v.X != 2 {
		t.Errorf("v(2, 3) * IdentityVec() x != 2 is %d", v.X)
	}
	if v.Y != 3 {
		t.Errorf("v(2, 3) * IdentityVec() y != 3 is %d", v.Y)
	}
}

func TestVecZeroMultiply(t *testing.T) {
	v := V(2, 3).Multiply(ZeroVec())
	if v.X != 0 {
		t.Errorf("v(2, 3) * ZeroVec() x != 0 is %d", v.X)
	}
	if v.Y != 0 {
		t.Errorf("v(2, 3) * ZeroVec() y != 0 is %d", v.Y)
	}
}

func TestVecZeroAdd(t *testing.T) {
	v := V(2, 3).Add(ZeroVec())
	if v.X != 2 {
		t.Errorf("v(2, 3) + ZeroVec() x != 2 is %d", v.X)
	}
	if v.Y != 3 {
		t.Errorf("v(2, 3) + ZeroVec() y != 3 is %d", v.Y)
	}
}

func TestToPixelVec(t *testing.T) {
	v := V(12, 18).ToPixelVec()
	if v.X != 12.0 {
		t.Errorf("v(12, 18).ToPixelVec().X != 12.0 is %f", v.X)
	}
	if v.Y != 18.0 {
		t.Errorf("v(12, 18).ToPixelVec().X != 18.0 is %f", v.Y)
	}
}

func TestToPixelRectZeroOne(t *testing.T) {
	r := ZeroVec().ToPixelRect(IdentityVec())
	if r.Min.X != 0.0 {
		t.Errorf("r.Min.X != 0.0")
	}
	if r.Min.Y != 0.0 {
		t.Errorf("r.Min.Y != 0.0")
	}
	if r.Max.X != 1.0 {
		t.Errorf("r.Min.X != 1.0")
	}
	if r.Max.Y != 1.0 {
		t.Errorf("r.Min.Y != 1.0")
	}
}

func TestToPixelRectZeroTen(t *testing.T) {
	r := ZeroVec().ToPixelRect(V(10, 10))
	if r.Min.X != 0.0 {
		t.Errorf("r.Min.X != 0.0")
	}
	if r.Min.Y != 0.0 {
		t.Errorf("r.Min.Y != 0.0")
	}
	if r.Max.X != 10.0 {
		t.Errorf("r.Min.X != 10.0")
	}
	if r.Max.Y != 10.0 {
		t.Errorf("r.Min.Y != 10.0")
	}
}

func TestToPixelRectIdentityTen(t *testing.T) {
	r := IdentityVec().ToPixelRect(V(10, 10))
	if r.Min.X != 10.0 {
		t.Errorf("r.Min.X != 10.0")
	}
	if r.Min.Y != 10.0 {
		t.Errorf("r.Min.Y != 10.0")
	}
	if r.Max.X != 20.0 {
		t.Errorf("r.Min.X != 20.0")
	}
	if r.Max.Y != 20.0 {
		t.Errorf("r.Min.Y != 20.0")
	}
}

func TestToPixelRectOffsetZeroOne(t *testing.T) {
	r := ZeroVec().ToPixelRectOffset(IdentityVec(), V(5, 5))
	if r.Min.X != 5.0 {
		t.Errorf("r.Min.X != 5.0")
	}
	if r.Min.Y != 5.0 {
		t.Errorf("r.Min.Y != 5.0")
	}
	if r.Max.X != 6.0 {
		t.Errorf("r.Min.X != 6.0")
	}
	if r.Max.Y != 6.0 {
		t.Errorf("r.Min.Y != 6.0")
	}
}

func TestToPixelRectOffsetZeroTen(t *testing.T) {
	r := ZeroVec().ToPixelRectOffset(V(10, 10), V(5, 5))
	if r.Min.X != 5.0 {
		t.Errorf("r.Min.X != 5.0")
	}
	if r.Min.Y != 5.0 {
		t.Errorf("r.Min.Y != 5.0")
	}
	if r.Max.X != 15.0 {
		t.Errorf("r.Min.X != 15.0")
	}
	if r.Max.Y != 15.0 {
		t.Errorf("r.Min.Y != 15.0")
	}
}

func TestToPixelRectOffsetIdentityTen(t *testing.T) {
	r := IdentityVec().ToPixelRectOffset(V(10, 10), V(5, 5))
	if r.Min.X != 15.0 {
		t.Errorf("r.Min.X != 15.0")
	}
	if r.Min.Y != 15.0 {
		t.Errorf("r.Min.Y != 15.0")
	}
	if r.Max.X != 25.0 {
		t.Errorf("r.Min.X != 25.0")
	}
	if r.Max.Y != 25.0 {
		t.Errorf("r.Min.Y != 25.0")
	}
}
