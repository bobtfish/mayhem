package main

import (
	"github.com/bobtfish/mayhem/logical"
)

type PlayerList []Player

type Player interface {
	GameObjectStackable
}

type HumanPlayer struct {
	Name       string
	BaseSprite logical.Vec
	Spells     []Spell
	NextSpell  Spell
}

func NewHumanPlayer(name string, spriteId int) *HumanPlayer {
	if spriteId < 0 || spriteId > 7 {
		panic("spriteId wrong")
	}
	spells := make([]Spell, 1)
	spells[0] = &DoNothingSpell{}
	return &HumanPlayer{
		Name:       name,
		BaseSprite: logical.V(spriteId, 23),
		Spells:     spells,
	}
}

// GameObject interface
func (p *HumanPlayer) AnimationTick() {}

func (p *HumanPlayer) GetSpriteSheetCoordinates() logical.Vec {
	return p.BaseSprite
}

// GameObjectStackable interface
func (h *HumanPlayer) RemoveMe() bool {
	return false
}
