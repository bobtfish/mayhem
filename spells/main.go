package spells

import (
	"fmt"
	"image/color"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/rand"
	"github.com/bobtfish/mayhem/render"
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
}

type ASpell struct {
	Name          string
	LawRating     int
	Reuseable     bool
	CastingChance int
	CastRange     int
	NoCastFx      bool
	Tries         int
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
func (s ASpell) GetCastingChance(playerLawRating int) int {
	// FIXME - adjust casting chance based on law rating
	// of the player and of the spell
	return s.CastingChance
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
func (s ASpell) CastSucceeds(illusion bool, playerLawRating int) bool {
	if illusion && !s.CanCastAsIllusion() {
		panic(fmt.Sprintf("Spell %s (type %T) cannot be illusion, but was cast as one anyway", s.Name, s))
	}
	if illusion || rand.Intn(100) <= s.GetCastingChance(playerLawRating) {
		return true
	}
	return false
}

type CreatureSpell struct {
	ASpell
}

type OtherSpell struct {
	ASpell
	MutateFunc func(logical.Vec, *grid.GameGrid, grid.GameObject) (bool, *fx.Fx)
}

func (s OtherSpell) DoCast(illusion bool, target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
	return s.MutateFunc(target, grid, owner)
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
		AllSpells = make([]Spell, 1) // Deliberately leave room for disbelieve spell as number 0
	}
	AllSpells = append(AllSpells, s)
}

func init() {
	/*CreateSpell(OtherSpell{
		ASpell: ASpell{ // Uses disbelive animation if it kills a thing. No corpse
			Name:          "Lightning",
			CastingChance: 100,
			CastRange:     4,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) bool {
			a, isAttackable := grid.GetGameObject(target).(movable.Attackable)
			if !isAttackable {
				return false
			}
			if rand.Intn(9)+3 > a.GetDefence() {
				fmt.Printf("Killed by lightning\n")
				return true
			}
			return true
		},
	})
	CreateSpell(OtherSpell{
		ASpell: ASpell{ // as above, just less strong
			Name:          "Magic Bolt",
			CastingChance: 100,
			CastRange:     6,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) bool {
			return false
		},
	})*/
}
