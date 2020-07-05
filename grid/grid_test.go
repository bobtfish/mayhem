package grid

import (
	"testing"

	"github.com/bobtfish/mayhem/logical"
)

func TestGrid(t *testing.T) {
	grid := MakeGameGrid(3, 3)
	if grid.Width() != 3 {
		t.Errorf("3,3 grid width not 3 is %d", grid.Width())
	}
	if grid.Height() != 3 {
		t.Errorf("3,3 grid height not 3 is %d", grid.Height())
	}
	r := grid.AsRect()
	if r.Width() != 3 {
		t.Errorf("3,3 grid AsRect width not 3 is %d", r.Width())
	}
	if r.Height() != 3 {
		t.Errorf("3,3 grid Height width not 3 is %d", r.Height())
	}
}

func TestGridTwoByThree(t *testing.T) {
	grid := MakeGameGrid(2, 3)
	if grid.Width() != 2 {
		t.Errorf("2,3 grid width not 2 is %d", grid.Width())
	}
	if grid.Height() != 3 {
		t.Errorf("2,3 grid height not 3 is %d", grid.Height())
	}

	clamped := grid.AsRect().Clamp(logical.V(20, 20))
	if clamped.X != 1 {
		t.Errorf("clamp vec of width 2 grid != 1, is %d", clamped.X)
	}
	if clamped.Y != 2 {
		t.Errorf("clamp vec of height 3 grid != 2, is %d", clamped.Y)
	}
}
