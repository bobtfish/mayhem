package spells

type Spell struct {
	Name string
}

var AllSpells []Spell

func init() {
	AllSpells = []Spell{
		Spell{Name: "Spell1"},
		Spell{Name: "Spell2"},
		Spell{Name: "Spell3"},
	}
}
