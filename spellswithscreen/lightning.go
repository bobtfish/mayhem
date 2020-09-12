package spellswithscreen

import (
	"fmt"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
	"github.com/bobtfish/mayhem/spells"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type ScreenSpell struct {
	spells.ASpell
}

func (s ScreenSpell) TakeOverScreen(grid *grid.GameGrid, cleanupFunc func(), nextScreen screeniface.GameScreen, target logical.Vec) screeniface.GameScreen {
	return &LightningSpellScreen{
		Grid:        grid,
		NextScreen:  nextScreen,
		CleanupFunc: cleanupFunc,
		Target:      target,
	}
}

func (s ScreenSpell) DoCast(illusion bool, target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
	return false, nil
}

func (s ScreenSpell) CastFx() *fx.Fx {
	return nil
}

type LightningSpellScreen struct {
	Grid        *grid.GameGrid
	NextScreen  screeniface.GameScreen
	CleanupFunc func()
	Target      logical.Vec
}

func (screen *LightningSpellScreen) DrawBoard(ss pixel.Picture, win *pixelgl.Window) *pixel.Batch {
	sd := render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())
	return screen.Grid.DrawBatch(sd)
}

func (screen *LightningSpellScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	fmt.Printf("FOO\n")
}

func (screen *LightningSpellScreen) Step(ss pixel.Picture, win *pixelgl.Window) screeniface.GameScreen {
	batch := screen.DrawBoard(ss, win)
	batch.Draw(win)
	if 0 == 1 {
		screen.CleanupFunc()
		return screen.NextScreen
	}
	return screen
}

func init() {
	spells.CreateSpell(ScreenSpell{
		ASpell: spells.ASpell{ // Uses disbelive animation if it kills a thing. No corpse
			Name:          "Lightning",
			CastingChance: 100,
			CastRange:     4,
		},
		/*MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
			a, isAttackable := grid.GetGameObject(target).(movable.Attackable)
			if !isAttackable {
				return false, nil
			}
			if rand.Intn(9)+3 > a.GetDefence() {
				fmt.Printf("Killed by lightning\n")
				return true, nil
			}
			return true, nil
		},*/
	})
	spells.CreateSpell(ScreenSpell{
		ASpell: spells.ASpell{ // as above, just less strong
			Name:          "Magic Bolt",
			CastingChance: 100,
			CastRange:     6,
		},
		/*MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
			return false, nil
		},*/
	})
}
