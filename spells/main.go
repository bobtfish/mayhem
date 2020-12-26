package spells

import (
	"fmt"
	"image/color"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/rand"
	"github.com/bobtfish/mayhem/render"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

type Spell interface {
	GetName() string
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
	TakeOverScreen(*grid.GameGrid, func(), screeniface.GameScreen, logical.Vec, logical.Vec) screeniface.GameScreen
}

type ASpell struct {
	Name                string
	LawRating           int
	Reuseable           bool
	CastingChance       int
	CastRange           int
	NoCastFx            bool
	Tries               int
	NoLineOfSightNeeded bool
}

func (s ASpell) TakeOverScreen(grid *grid.GameGrid, cleanup func(), nextScreen screeniface.GameScreen, source, target logical.Vec) screeniface.GameScreen {
	return nil
}
func (s ASpell) CastFx() *fx.Fx {
	if s.NoCastFx {
		return nil
	}
	return fx.SpellCast()
}
func (s ASpell) GetName() string {
	return s.Name
}
func (s ASpell) GetLawRating() int {
	return s.LawRating
}
func (s ASpell) IsReuseable() bool {
	return s.Reuseable
}
func (s ASpell) CastQuantity() int {
	if s.Tries == 0 {
		return 1
	}
	return s.Tries
}
func (s ASpell) GetCastingChance(gameLawRating int) int {
	// Neutral spells or neutral games give base chance.
	// Spells never get harder, so a lawful game + chaos spell or chaotic game + lawful spell both give base chance
	if s.LawRating == 0 || gameLawRating == 0 || (s.LawRating > 0 && gameLawRating < 0) || (s.LawRating < 0 && gameLawRating > 0) {
		return s.CastingChance
	}
	absLaw := gameLawRating
	if absLaw < 0 {
		absLaw = -absLaw
	}
	cc := s.CastingChance + (absLaw/4)*10 // Every 4 points makes it 10% easier
	if cc > 100 {
		return 100
	}
	return cc
}

func (s ASpell) GetCastRange() int {
	return s.CastRange
}

func (s ASpell) CanCast(target grid.GameObject) bool {
	return true
}
func (s ASpell) CanCastAsIllusion() bool {
	return false
}
func (s ASpell) CastSucceeds(illusion bool, gameLawRating int) bool {
	if illusion && !s.CanCastAsIllusion() {
		panic(fmt.Sprintf("Spell %s (type %T) cannot be illusion, but was cast as one anyway", s.Name, s))
	}
	if illusion || rand.Intn(100) <= s.GetCastingChance(gameLawRating) {
		return true
	}
	return false
}
func (s ASpell) NeedsLineOfSight() bool {
	return !s.NoLineOfSightNeeded
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

func CastingChanceColor(chance int) color.Color {
	switch {
	case chance >= 100:
		return render.GetColor(255, 255, 255)
	case chance >= 80:
		return render.GetColor(255, 255, 0)
	case chance >= 60:
		return render.GetColor(0, 255, 255)
	case chance >= 40:
		return render.GetColor(0, 255, 0)
	default:
		return render.GetColor(255, 0, 255)
	}
}

func ChooseSpells() []Spell {
	spells := make([]Spell, 14)
	spells[0] = AllSpells[0]
	for i := 1; i < 14; i++ {
		idx := rand.Intn(len(AllSpells)-1) + 1
		spells[i] = AllSpells[idx]
	}
	/*
		for i := 0; i < len(AllSpells); i++ {
			if AllSpells[i].GetName() == "Magic Bolt" {
				spells[1] = AllSpells[i]
			}
		}*/
	return spells
}

var AllSpells []Spell

func CreateSpell(s Spell) {
	if AllSpells == nil {
		AllSpells = make([]Spell, 1) // Deliberately leave room for disbelieve spell as number 0
	}
	AllSpells = append(AllSpells, s)
}
