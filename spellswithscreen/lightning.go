package spellswithscreen

import (
	"time"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
	screens "github.com/bobtfish/mayhem/screen"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
	"github.com/bobtfish/mayhem/spells"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type ScreenSpell struct {
	spells.ASpell
}

func (s ScreenSpell) TakeOverScreen(grid *grid.GameGrid, cleanupFunc func(), nextScreen screeniface.GameScreen, source, target logical.Vec) screeniface.GameScreen {
	four := logical.V(4, 4)
	mtarget := target.Multiply(four)
	msource := source.Multiply(four)
	anim := mtarget.Subtract(msource).Path()
	for i, s := range anim {
		anim[i] = msource.Add(s)
	}
	anim = append(anim, mtarget)
	return &LightningSpellScreen{
		Grid:        grid,
		NextScreen:  nextScreen,
		CleanupFunc: cleanupFunc,
		Source:      source,
		Target:      target,
		Anim:        anim,
		AnimCount:   -1,
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
	CleanupFunc func() // This is a closure that removes the spell from the player after casting, called when leaving
	Source      logical.Vec
	Target      logical.Vec
	Anim        []logical.Vec
	AnimCount   int
}

func (screen *LightningSpellScreen) DrawBoard(ss pixel.Picture, win *pixelgl.Window) *pixel.Batch {
	sd := render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())
	return screen.Grid.DrawBatch(sd)
}

func (screen *LightningSpellScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	//fmt.Printf("FOO\n")
}

// 25 Y, 0 X

func (screen *LightningSpellScreen) Step(ss pixel.Picture, win *pixelgl.Window) screeniface.GameScreen {
	batch := screen.DrawBoard(ss, win)
	screen.AnimCount++
	//fmt.Printf("Source is X%dY%d Target is X%dY%d\n", screen.Source.X, screen.Source.Y, screen.Target.X, screen.Target.Y)
	//fmt.Printf("AnimCount is %d Path is %d, %d\n", screen.AnimCount, screen.Anim[screen.AnimCount].X, screen.Anim[screen.AnimCount].Y)

	color := render.GetColor(255, 255, 255)
	sd := render.NewSpriteQuarterDrawer(ss).WithOffset(render.GameBoardV())
	for i := 0; i < screen.AnimCount; i++ {
		winPos := screen.Anim[i]
		//fmt.Printf("Draw at %d %d\n", winPos.X, winPos.Y)
		sd.DrawSpriteColor(logical.V(7, 25), winPos, color, batch)
	}

	batch.Draw(win)
	if screen.AnimCount+1 == len(screen.Anim) {
		// Do something?
		screen.CleanupFunc()
		return screen.NextScreen
	}
	return &screens.Pause{
		NextScreen: screen,
		For:        time.Microsecond,
	}
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
