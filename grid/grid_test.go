package grid

import (
    "testing"

    "github.com/bobtfish/mayhem/logical"
)

func TestGrid(t *testing.T) {
    grid := MakeGameGrid(logical.V(3, 3))
    if grid.Width() != 3 {
        t.Errorf("3,3 grid width not 3")
    }
    if grid.Height() != 3 {
        t.Errorf("3,3 grid height not 3")
    }
    r := grid.AsRect()
    if r.Width() != 3 {
        t.Errorf("3,3 grid AsRect width not 3")
    }
    if r.Height() != 3 {
        t.Errorf("3,3 grid Height width not 3")
    }
}
