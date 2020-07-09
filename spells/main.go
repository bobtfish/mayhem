package spells

import (
	"image/color"
	"math/rand"

	"github.com/bobtfish/mayhem/render"
)

type Spell interface {
	GetName() string
	GetLawRating() int
	GetCastingChance(int) int
	GetRange() int
}

type ASpell struct {
	Name          string
	LawRating     int
	Reuseable     bool
	CastingChance int
	Range         int
}

func (s ASpell) GetName() string {
	return s.Name
}
func (s ASpell) GetLawRating() int {
	return s.LawRating
}
func (s ASpell) GetCastingChance(playerLawRating int) int {
	// FIXME - adjust casting chance based on law rating
	// of the player and of the spell
	return s.CastingChance
}
func (s ASpell) GetRange() int {
	return s.Range
}

type DisbelieveSpell struct {
	ASpell
}

type CreatureSpell struct {
	ASpell
}

type OtherSpell struct {
	ASpell
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

func init() {
	AllSpells = []Spell{
		DisbelieveSpell{ASpell{
			Name:          "Disbelieve",
			LawRating:     0,
			Reuseable:     true,
			CastingChance: 100,
			Range:         20,
		}},
		OtherSpell{ASpell{
			Name:          "Raise Dead",
			LawRating:     -1,
			CastingChance: 60,
			Range:         4,
		}},
		OtherSpell{ASpell{
			Name:          "Magic Knife",
			LawRating:     1,
			CastingChance: 90,
		}},
		OtherSpell{ASpell{
			Name:          "Magic Armour",
			LawRating:     1,
			CastingChance: 40,
		}},
		OtherSpell{ASpell{
			Name:          "Magic Shield",
			LawRating:     1,
			CastingChance: 80,
		}},
		OtherSpell{ASpell{
			Name:          "Magic Bow",
			LawRating:     1,
			CastingChance: 50,
		}},
		OtherSpell{ASpell{
			Name:          "Shadow Form",
			CastingChance: 80,
		}},
		OtherSpell{ASpell{
			Name:          "Magic Wings",
			CastingChance: 60,
		}},
		OtherSpell{ASpell{
			Name:          "Lightning",
			CastingChance: 100,
			Range:         4,
		}},
		OtherSpell{ASpell{
			Name:          "Magic Bolt",
			CastingChance: 100,
			Range:         6,
		}},
		OtherSpell{ASpell{
			Name:          "Law-1",
			CastingChance: 100,
			LawRating:     2,
		}},
		OtherSpell{ASpell{
			Name:          "Law-2",
			CastingChance: 60,
			LawRating:     4,
		}},
		OtherSpell{ASpell{
			Name:          "Chaos-1",
			CastingChance: 80,
			LawRating:     -2,
		}},
		OtherSpell{ASpell{
			Name:          "Chaos-2",
			CastingChance: 60,
			LawRating:     -4,
		}},
		OtherSpell{ASpell{
			Name:          "Vengence",
			CastingChance: 80,
			Range:         20,
		}},
		OtherSpell{ASpell{
			Name:          "Decree",
			CastingChance: 80,
			Range:         20,
			LawRating:     1,
		}},
		OtherSpell{ASpell{
			Name:          "Dark Power",
			CastingChance: 50,
			Range:         20,
			LawRating:     -2,
		}},
		OtherSpell{ASpell{
			Name:          "Justice",
			CastingChance: 50,
			Range:         20,
			LawRating:     2,
		}},
		OtherSpell{ASpell{
			Name:          "Subversion",
			CastingChance: 100,
			Range:         7,
		}},
	}
}
