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
	r := ZeroVec().ToPixelRect(IdentityVec(), V(5, 5))
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
	r := ZeroVec().ToPixelRect(V(10, 10), V(5, 5))
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
	r := IdentityVec().ToPixelRect(V(10, 10), V(5, 5))
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

func TestAbs(t *testing.T) {
	v := V(1, 1).Abs()
	if v.X != 1 {
		t.Errorf("v(1, 1) abs X != 1 is %d", v.X)
	}
	if v.Y != 1 {
		t.Errorf("v(1, 1) abs Y != 1 is %d", v.Y)
	}
	v = V(-2, -3).Abs()
	if v.X != 2 {
		t.Errorf("v(-2, -3) abs X != 2 is %d", v.X)
	}
	if v.Y != 3 {
		t.Errorf("v(-2, -3) abs Y != 3 is %d", v.Y)
	}
}

func TestSmallestSquare(t *testing.T) {
	l := V(1, 3).smallestSquare()
	if l != 1 {
		t.Errorf("v(1, 3) smallestSquare l != 1 is %d", l)
	}
	l = V(-1, 3).smallestSquare()
	if l != 1 {
		t.Errorf("v(-1, 3) smallestSquare l != 1 is %d", l)
	}
	l = V(1, -3).smallestSquare()
	if l != 1 {
		t.Errorf("v(1, -3) smallestSquare l != 1 is %d", l)
	}
	l = V(-26, 3).smallestSquare()
	if l != 3 {
		t.Errorf("v(26, 3) smallestSquare X != 26 is %d", l)
	}
}

func TestSquareDistance(t *testing.T) {
	table := map[int]int{
		0: 0,
		1: 1,
		2: 3,
		3: 4,
		4: 6,
	}
	for l, expD := range table {
		d := squareDistance(l)
		if d != expD {
			t.Errorf("squareDistance(%d) != %d is %d", l, expD, d)
		}
	}
}

func TestVecDistance(t *testing.T) {
	d := ZeroVec().Distance(ZeroVec())
	if d != 0 {
		t.Errorf("v(0, 0) to v(0, 0) should be distance 0, not %d", d)
	}
	d = ZeroVec().Distance(IdentityVec())
	if d != 1 {
		t.Errorf("v(0, 0) to v(1, 1) should be distance 1, not %d", d)
	}
	d = ZeroVec().Distance(V(2, 2))
	if d != 3 {
		t.Errorf("v(0, 0) to v(2, 2) should be distance 3, not %d", d)
	}
	d = ZeroVec().Distance(V(3, 3))
	if d != 4 {
		t.Errorf("v(0, 0) to v(3, 3) should be distance 4, not %d", d)
	}
	d = ZeroVec().Distance(V(4, 4))
	if d != 6 {
		t.Errorf("v(0, 0) to v(4, 4) should be distance 6, not %d", d)
	}
	d = ZeroVec().Distance(V(10, 4))
	if d != 12 {
		t.Errorf("v(0, 0) to v(10, 4) should be distance 12, not %d", d)
	}
}
