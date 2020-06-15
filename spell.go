package main

type Spell interface {
	Name() string
}

type DoNothingSpell struct{}

func (s DoNothingSpell) Name() string {
	return "Do Nothing"
}
