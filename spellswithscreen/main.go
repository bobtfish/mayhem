package spellswithscreen

import (
	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
	"github.com/bobtfish/mayhem/spells"
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
