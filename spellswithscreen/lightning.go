package spellswithscreen

import (
	"fmt"
	"math/rand"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/movable"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
	"github.com/bobtfish/mayhem/spells"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type ScreenSpell struct {
	spells.ASpell
	MutateFunc func(logical.Vec, *grid.GameGrid, grid.GameObject) (bool, *fx.Fx)
}

func (s ScreenSpell) TakeOverScreen(grid *grid.GameGrid, nextScreen screeniface.GameScreen) screeniface.GameScreen {
	return nil
}

func (s ScreenSpell) DoCast(illusion bool, target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
	return s.MutateFunc(target, grid, owner)
}

type LightningSpellScreen struct {
	Grid *grid.GameGrid
}

func (screen *LightningSpellScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
}

func (screen *LightningSpellScreen) Step(ss pixel.Picture, win *pixelgl.Window) screeniface.GameScreen {
	return screen
}

func init() {
	spells.CreateSpell(ScreenSpell{
		ASpell: spells.ASpell{ // Uses disbelive animation if it kills a thing. No corpse
			Name:          "Lightning",
			CastingChance: 100,
			CastRange:     4,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
			a, isAttackable := grid.GetGameObject(target).(movable.Attackable)
			if !isAttackable {
				return false, nil
			}
			if rand.Intn(9)+3 > a.GetDefence() {
				fmt.Printf("Killed by lightning\n")
				return true, nil
			}
			return true, nil
		},
	})
	spells.CreateSpell(ScreenSpell{
		ASpell: spells.ASpell{ // as above, just less strong
			Name:          "Magic Bolt",
			CastingChance: 100,
			CastRange:     6,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
			return false, nil
		},
	})
}
