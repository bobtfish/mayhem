package spells

import (
	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
)

type LawSpell struct {
	ASpell
}

func (s LawSpell) DoCast(illusion bool, target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
	return true, nil
}

func init() {
	CreateSpell(LawSpell{ASpell{
		Name:          "Law-1",
		CastingChance: 100,
		LawRating:     2,
		NoCastFx:      true,
	}})
	CreateSpell(LawSpell{ASpell{
		Name:          "Law-2",
		CastingChance: 60,
		LawRating:     4,
		NoCastFx:      true,
	}})
	CreateSpell(LawSpell{ASpell{
		Name:          "Chaos-1",
		CastingChance: 80,
		LawRating:     -2,
		NoCastFx:      true,
	}})
	CreateSpell(LawSpell{ASpell{
		Name:          "Chaos-2",
		CastingChance: 60,
		LawRating:     -4,
		NoCastFx:      true,
	}})
}
