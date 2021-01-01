package spellswithscreen

import (
	"fmt"
	"time"

	"github.com/bobtfish/mayhem/rand"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/game"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/movable"
	"github.com/bobtfish/mayhem/render"
	screens "github.com/bobtfish/mayhem/screen"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
	"github.com/bobtfish/mayhem/spells"
	spelliface "github.com/bobtfish/mayhem/spells/iface"
	"github.com/faiface/pixel"
)

// FIXME this is duplicate
func DrawBoard(ctx screeniface.GameCtx) *pixel.Batch {
	return ctx.GetGrid().DrawBatch(render.NewSpriteDrawer(ctx.GetSpriteSheet()).WithOffset(render.GameBoardV()))
}

type ScreenSpell struct {
	spells.ASpell
	TakeOverFunc func(screeniface.GameCtx, func(), screeniface.GameScreen, int, logical.Vec) screeniface.GameScreen
}

func (s ScreenSpell) TakeOverScreen(ctx screeniface.GameCtx, cleanupFunc func(), nextScreen screeniface.GameScreen, playerIdx int, target logical.Vec) screeniface.GameScreen {
	return s.TakeOverFunc(ctx, cleanupFunc, nextScreen, playerIdx, target)
}

func (s ScreenSpell) DoCast(illusion bool, target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
	return false, nil
}

func (s ScreenSpell) CastFx() *fx.Fx {
	return nil
}

type LightningSpellScreen struct {
	NextScreen  screeniface.GameScreen
	CleanupFunc func() // This is a closure that removes the spell from the player after casting, called when leaving
	Target      logical.Vec
	Anim        []logical.Vec
	AnimCount   int
	Lightning   bool
}

func (screen *LightningSpellScreen) Enter(ctx screeniface.GameCtx) {
	//fmt.Printf("FOO\n")
}

// 25 Y, 0 X

func (screen *LightningSpellScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	batch := DrawBoard(ctx)
	ss := ctx.GetSpriteSheet()
	grid := ctx.GetGrid()
	win := ctx.GetWindow()
	screen.AnimCount++
	//fmt.Printf("Source is X%dY%d Target is X%dY%d\n", screen.Source.X, screen.Source.Y, screen.Target.X, screen.Target.Y)
	//fmt.Printf("AnimCount is %d Path is %d, %d\n", screen.AnimCount, screen.Anim[screen.AnimCount].X, screen.Anim[screen.AnimCount].Y)

	color := render.GetColor(255, 255, 255)
	sd := render.NewSpriteQuarterDrawer(ss).WithOffset(render.GameBoardV())
	startAt := 0
	if !screen.Lightning { // Magic bolt - lightning is a solid animation, bolt isn't
		if screen.AnimCount > 1 {
			startAt = screen.AnimCount - 1
		}
	}
	for i := startAt; i < screen.AnimCount; i++ {
		winPos := screen.Anim[i]
		//fmt.Printf("Draw at %d %d\n", winPos.X, winPos.Y)
		sd.DrawSpriteColor(logical.V(7, 25), winPos, color, batch)
	}

	batch.Draw(win)
	if screen.AnimCount+1 == len(screen.Anim) {
		ob := grid.GetGameObject(screen.Target)
		a, isAttackable := ob.(movable.Attackable)
		if isAttackable {
			chance := rand.Intn(5) // Defence is 1-5 for a player
			if screen.Lightning {
				chance += 2
			}
			fmt.Printf("Chance %d > Defence %d\n", chance, a.GetDefence())
			if chance > a.GetDefence() {
				died := grid.GetGameObjectStack(screen.Target).RemoveTopObject()
				if screens.KillIfPlayer(died, grid) {
					players := ctx.(*game.Window).GetPlayers()
					if screens.WeHaveAWinner(players) {
						return &screens.WinnerScreen{
							Players: players,
						}
					}
				}
			}
		}
		screen.CleanupFunc()
		return screen.NextScreen
	}
	return &screens.Pause{
		NextScreen: screen,
		For:        time.Microsecond,
	}
}

func init() {
	//ctx screeniface.GameCtx, cleanupFunc func(), nextScreen screeniface.GameScreen, playeerIdx int, target logical.Vec) screeniface.GameScreen
	lightningTakeOver := func(isLightning bool) func(ctx screeniface.GameCtx, cleanupFunc func(), nextScreen screeniface.GameScreen, playerIdx int, target logical.Vec) screeniface.GameScreen {
		return func(ctx screeniface.GameCtx, cleanupFunc func(), nextScreen screeniface.GameScreen, playerIdx int, target logical.Vec) screeniface.GameScreen {
			four := logical.V(4, 4)
			mtarget := target.Multiply(four)
			source := ctx.(*game.Window).GetPlayers()[playerIdx].BoardPosition
			msource := source.Multiply(four)
			anim := mtarget.Subtract(msource).Path()
			for i, s := range anim {
				anim[i] = msource.Add(s)
			}
			anim = append(anim, mtarget)
			return &LightningSpellScreen{
				Lightning:   isLightning,
				NextScreen:  nextScreen,
				CleanupFunc: cleanupFunc,
				Target:      target,
				Anim:        anim,
				AnimCount:   -1,
			}
		}
	}
	spelliface.CreateSpell(ScreenSpell{
		ASpell: spells.ASpell{ // Uses disbelive animation if it kills a thing. No corpse
			Name:          "Lightning",
			CastingChance: 100,
			CastRange:     4,
		},
		TakeOverFunc: lightningTakeOver(true),
	})
	spelliface.CreateSpell(ScreenSpell{
		ASpell: spells.ASpell{ // as above, just less strong
			Name:          "Magic Bolt",
			CastingChance: 100,
			CastRange:     6,
		},
		TakeOverFunc: lightningTakeOver(false),
	})
}
