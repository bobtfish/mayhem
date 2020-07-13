package spells

import (
	"image/color"
	"math/rand"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/render"
)

type Spell interface {
	GetName() string
	GetLawRating() int
	GetCastingChance(int) int
	GetRange() int
	DoesCastWork(int) bool
	Cast(grid.GameObject)
	IsReuseable() bool
	CastFx() *fx.Fx
}

type ASpell struct {
	Name          string
	LawRating     int
	Reuseable     bool
	CastingChance int
	Range         int
	NoCastFx      bool
}

func (s ASpell) CastFx() *fx.Fx {
	if s.NoCastFx {
		return nil
	}
	return fx.FxSpellCast()
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
func (s ASpell) GetCastingChance(playerLawRating int) int {
	// FIXME - adjust casting chance based on law rating
	// of the player and of the spell
	return s.CastingChance
}
func (s ASpell) GetRange() int {
	return s.Range
}

func (s ASpell) DoesCastWork(playerLawRating int) bool {
	// FIXME
	return true
	if rand.Intn(100) <= s.GetCastingChance(playerLawRating) {
		return true
	}
	return false
}
func (s ASpell) CanCast(target grid.GameObject) bool {
	return true
}

type DisbelieveSpell struct {
	ASpell
}

func (s DisbelieveSpell) Cast(target grid.GameObject) {
}

type CreatureSpell struct {
	ASpell
}

type OtherSpell struct {
	ASpell
}

func (s OtherSpell) Cast(target grid.GameObject) {
}

func LawRatingSymbol(s Spell) string {
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
		idx := rand.Intn(len(AllSpells)-2) + 1
		spells[i] = AllSpells[idx]
	}
	return spells
}

var AllSpells []Spell

func CreateSpell(s Spell) {
	if AllSpells == nil {
		AllSpells = make([]Spell, 1)
		AllSpells[0] = DisbelieveSpell{ASpell{
			Name:          "Disbelieve",
			LawRating:     0,
			Reuseable:     true,
			CastingChance: 100,
			Range:         20,
		}}
	}
	AllSpells = append(AllSpells, s)
}

func init() {
	CreateSpell(OtherSpell{ASpell{
		Name:          "Raise Dead",
		LawRating:     -1,
		CastingChance: 60,
		Range:         4,
	}})
	CreateSpell(OtherSpell{ASpell{
		Name:          "Lightning",
		CastingChance: 100,
		Range:         4,
	}})
	CreateSpell(OtherSpell{ASpell{
		Name:          "Magic Bolt",
		CastingChance: 100,
		Range:         6,
	}})
	CreateSpell(OtherSpell{ASpell{
		Name:          "Vengence",
		CastingChance: 80,
		Range:         20,
	}})
	CreateSpell(OtherSpell{ASpell{
		Name:          "Decree",
		CastingChance: 80,
		Range:         20,
		LawRating:     1,
	}})
	CreateSpell(OtherSpell{ASpell{
		Name:          "Dark Power",
		CastingChance: 50,
		Range:         20,
		LawRating:     -2,
	}})
	CreateSpell(OtherSpell{ASpell{
		Name:          "Justice",
		CastingChance: 50,
		Range:         20,
		LawRating:     2,
	}})
	CreateSpell(OtherSpell{ASpell{
		Name:          "Subversion",
		CastingChance: 100,
		Range:         7,
	}})
}
