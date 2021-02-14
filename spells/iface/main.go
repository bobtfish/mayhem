package iface

import (
	"github.com/bobtfish/mayhem/rand"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"

	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

type Spell interface {
	GetName() string
	GetDescriptionArray(int) []string
	GetLawRating() int
	GetCastingChance(int) int
	GetCastRange() int
	CanCast(grid.GameObject) bool
	CanCastAsIllusion() bool
	DoCast(bool, logical.Vec, *grid.GameGrid, grid.GameObject) (bool, *fx.Fx)
	CastQuantity() int
	CastSucceeds(bool, int) bool
	IsReuseable() bool
	CastFx() *fx.Fx
	NeedsLineOfSight() bool

	TakeOverScreen(screeniface.GameCtx, func(), screeniface.GameScreen, int, logical.Vec) screeniface.GameScreen
}

func LawRatingSymbol(s Spell) string {
	if s == nil {
		panic("nil spell")
	}
	if s.GetLawRating() == 0 {
		return "-"
	}
	if s.GetLawRating() < 0 {
		return "*"
	}
	return "^"
}

var AllSpells []Spell

func CreateSpell(s Spell) {
	if AllSpells == nil {
		AllSpells = make([]Spell, 1) // Deliberately leave room for disbelieve spell as number 0
	}
	AllSpells = append(AllSpells, s)
}

func ChooseSpell() Spell {
	idx := rand.Intn(len(AllSpells)-1) + 1
	return AllSpells[idx]
}

func ChooseSpells() []Spell {
	spells := make([]Spell, 14)
	spells[0] = AllSpells[0]
	for i := 1; i < 14; i++ {
		spells[i] = ChooseSpell()
	}

	for i := 0; i < len(AllSpells); i++ {
		if AllSpells[i].GetName() == "Elf" {
			spells[1] = AllSpells[i]
		}
	}
	return spells
}
