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

func TestVecAdd(t *testing.T) {
    v := V(1, 2).Add(V(1, 2))
    if v.X != 2 {
        t.Errorf("v(1, 2) + v(1, 2) x != 2")
    }
    if v.Y != 4 {
        t.Errorf("v(1, 2) + v(1, 2) x != 4")
    }
}

func TestVecSubtract(t *testing.T) {
    v := V(4, 8).Subtract(V(1, 2))
    if v.X != 3 {
        t.Errorf("v(4, 8) - v(1, 2) x != 3 is %d", v.X)
    }
    if v.Y != 6 {
        t.Errorf("v(4, 8) + v(1, 2) x != 6 is %d", v.Y)
    }
}
