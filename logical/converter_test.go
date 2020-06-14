package logical

import "testing"

func TestNullConverter(t *testing.T) {
	c := OffsetVecConverter{
		Offset:    V(0, 0),
		Multipler: 1,
	}
	v := c.ToPixelVec(V(0, 0))
	if v.X != 0.0 {
		t.Errorf("Zero vec not zero X")
	}
	if v.Y != 0.0 {
		t.Errorf("Zero vec not zero Y")
	}
	v = c.ToPixelVec(V(1, 1))
	if v.X != 1.0 {
		t.Errorf("1 Vec not 1 X")
	}
	if v.Y != 1.0 {
		t.Errorf("1 vec not 1 Y")
	}
}

func TestDoubleConverter(t *testing.T) {
	c := OffsetVecConverter{
		Offset:    V(0, 0),
		Multipler: 2,
	}
	v := c.ToPixelVec(V(0, 0))
	if v.X != 0.0 {
		t.Errorf("Zero vec not zero X")
	}
	if v.Y != 0.0 {
		t.Errorf("Zero vec not zero Y")
	}
	v = c.ToPixelVec(V(1, 1))
	if v.X != 2.0 {
		t.Errorf("1 Vec not 2 X")
	}
	if v.Y != 2.0 {
		t.Errorf("1 vec not 2 Y")
	}
}

func TestOffsetConverter(t *testing.T) {
	c := OffsetVecConverter{
		Offset:    V(10, 20),
		Multipler: 1,
	}
	v := c.ToPixelVec(V(0, 0))
	if v.X != 10.0 {
		t.Errorf("Zero vec not zero X")
	}
	if v.Y != 20.0 {
		t.Errorf("Zero vec not zero Y")
	}
	v = c.ToPixelVec(V(1, 1))
	if v.X != 11.0 {
		t.Errorf("1 Vec not 1 X")
	}
	if v.Y != 21.0 {
		t.Errorf("1 vec not 1 Y")
	}
}
