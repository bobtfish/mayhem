package logical

import "testing"

func TestRecZero(t *testing.T) {
	r := R(0, 0)
	if r.X != 0 {
		t.Errorf("Zero Rec not zero X")
	}
	if r.Y != 0 {
		t.Errorf("Zero Rec not zero Y")
	}
	if r.Contains(V(0, 0)) != true {
		t.Errorf("Zero Rec does not contain Zero Vec")
	}
	if r.Contains(V(-1, 0)) != false {
		t.Errorf("Zero Rec contains negative X vec")
	}
	if r.Contains(V(0, -1)) != false {
		t.Errorf("Zero Rec contains negative Y vec")
	}
	if r.Contains(V(1, 0)) != false {
		t.Errorf("Zero Rec contains X1 vec")
	}
	if r.Contains(V(0, 1)) != false {
		t.Errorf("Zero Rec contains Y1 vec")
	}
	if r.Contains(V(1, 1)) != false {
		t.Errorf("Zero Rec contains Y1 vec")
	}
	adj := r.Adjacents(V(0, 0))
	if len(adj) != 0 {
		t.Errorf("Zero Rec has adjacents")
	}
	if r.Width() != 0 {
		t.Errorf("width != 0")
	}
	if r.Height() != 0 {
		t.Errorf("Height != 0")
	}
}

func TestRecOne(t *testing.T) {
	r := R(1, 1)
	if r.Contains(V(0, 0)) != true {
		t.Errorf("One Rec does not contain Zero Vec")
	}
	if r.Contains(V(1, 1)) != true {
		t.Errorf("One Rec does not contain One Vec")
	}
	adj := r.Adjacents(V(0, 0))
	if len(adj) != 3 {
		t.Errorf("Zero vec in One rec does not have 3 adjacents")
	}
	if r.Width() != 1 {
		t.Errorf("width != 1")
	}
	if r.Height() != 1 {
		t.Errorf("Height != 1")
	}
}

func TestRecTwo(t *testing.T) {
	r := R(2, 2)
	adj := r.Adjacents(V(1, 1))
	if len(adj) != 8 {
		t.Errorf("One vec in Two rec does not have 8 adjacents")
	}
	adj = r.Adjacents(V(0, 1))
	if len(adj) != 5 {
		t.Errorf("0,1 vec in 2 rec does not have 5 adjacents")
	}
	if r.Width() != 2 {
		t.Errorf("width != 2")
	}
	if r.Height() != 2 {
		t.Errorf("Height != 2")
	}
}

func TestClamp(t *testing.T) {
	r := R(2, 2)
	if r.Clamp(V(1, 1)).X != 1 {
		t.Errorf("Clamp in R(2, 2) should not touch X of V(1, 1)")
	}
	if r.Clamp(V(1, 1)).Y != 1 {
		t.Errorf("Clamp in R(2, 2) should not touch Y of V(1, 1)")
	}
	if r.Clamp(V(3, 1)).X != 2 {
		t.Errorf("Clamp in R(2, 2) should clamp V(3, 1) to V(2, 1)")
	}
	if r.Clamp(V(1, 3)).Y != 2 {
		t.Errorf("Clamp in R(2, 2) should clamp V(1, 3) to V(1, 2)")
	}
	if r.Clamp(V(-1, -1)).X != 0 {
		t.Errorf("Clamp in R(2, 2) should clamp V(-1, -1) to V(0, 0) X")
	}
	if r.Clamp(V(-1, -1)).Y != 0 {
		t.Errorf("Clamp in R(2, 2) should clamp V(-1, -1) to V(0, 0) Y")
	}
}
