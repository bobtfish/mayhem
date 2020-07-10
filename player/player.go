package player

import (
	"image/color"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/spells"
)

func NewPlayer() Player {
	return Player{
		ChosenSpell: -1,
		Spells:      spells.ChooseSpells(),
	}
}

type Player struct {
	Name          string
	Spells        []spells.Spell
	HumanPlayer   bool
	CharacterIcon logical.Vec
	ChosenSpell   int
	Color         color.Color
	LawRating     int
	BoardPosition logical.Vec
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

func (p *Player) SetBoardPosition(v logical.Vec) {
	p.BoardPosition = v
}

// GameObjectStackable interface

func (p *Player) RemoveMe() bool {
	return false
}

func (p *Player) CastSpell() bool {
	i := p.ChosenSpell
	spell := p.Spells[i]
	if !spell.IsReuseable() {
		spells := make([]spells.Spell, 0)
		spells = append(p.Spells[:i], p.Spells[i+1:]...)
		p.Spells = spells
	}
	ret := spell.Cast(p.LawRating)
	p.ChosenSpell = -1
	return ret
}
