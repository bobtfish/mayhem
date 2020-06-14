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
