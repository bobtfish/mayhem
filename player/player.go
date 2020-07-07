package player

import (
	"image/color"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/spells"
)

type Player struct {
	Name          string
	Spells        []spells.Spell
	NextSpell     spells.Spell
	HumanPlayer   bool
	CharacterIcon logical.Vec
	ChosenSpell   int
	Color         color.Color
	LawRating     int
}

// GameObject interface
func (p *Player) AnimationTick() {}

func (p *Player) IsEmpty() bool {
	return false
}

func (p *Player) GetSpriteSheetCoordinates() logical.Vec {
	return p.CharacterIcon
}

func (p *Player) GetColor() color.Color {
	return p.Color
}

func (p *Player) Describe() string {
	return p.Name
}

// GameObjectStackable interface

func (h *Player) RemoveMe() bool {
	return false
}
