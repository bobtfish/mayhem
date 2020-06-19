package spells

type Spell struct {
	Name      string
	LawRating int
}

func (s Spell) LawRatingSymbol() string {
	if s.LawRating == 0 {
		return "-"
	}
	if s.LawRating < 0 {
		return "*"
	}
	return "^"
}

var AllSpells []Spell

func init() {
	AllSpells = []Spell{
		Spell{
			Name:      "Spell1",
			LawRating: 0,
		},
		Spell{
			Name:      "Spell2",
			LawRating: 1,
		},
		Spell{
			Name:      "Spell3",
			LawRating: -1,
		},
	}
}
