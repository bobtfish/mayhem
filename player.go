package main

import (
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/spells"
)

type PlayerList []Player

type Player interface {
	grid.GameObjectStackable
}

type HumanPlayer struct {
	Name       string
	BaseSprite logical.Vec
	Spells     []spells.Spell
	NextSpell  spells.Spell
}

func NewHumanPlayer(name string, spriteId int) *HumanPlayer {
	if spriteId < 0 || spriteId > 7 {
		panic("spriteId wrong")
	}
	sps := make([]spells.Spell, 1)
	sps[0] = spells.AllSpells[0]
	return &HumanPlayer{
		Name:       name,
		BaseSprite: logical.V(spriteId, 23),
		Spells:     sps,
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
