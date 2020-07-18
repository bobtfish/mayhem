package player

import (
	"fmt"
	"image/color"

	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/spells"
)

func NewPlayer() Player {
	return Player{
		Defence:     3,
		Combat:      3,
		Manoeuvre:   4,
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
	IsAnimated    bool
	SpriteIdx     int
	Defence       int
	Combat        int
	Manoeuvre     int

	Flying bool // If the player has magic wings
}

// GameObject interface
func (p *Player) AnimationTick(odd bool) {
	if odd {
		return
	}
	if p.IsAnimated {
		p.SpriteIdx++
		if p.SpriteIdx == 4 {
			p.SpriteIdx = 0
		}
	}
}

func (p *Player) IsEmpty() bool {
	return false
}

func (p *Player) GetSpriteSheetCoordinates() logical.Vec {
	return p.CharacterIcon.Add(logical.V(p.SpriteIdx, 0))
}

func (p *Player) GetColor() color.Color {
	return p.Color
}

func (p *Player) Describe() string {
	return p.Name
}

// GameObjectStackable interface

func (p *Player) RemoveMe() bool {
	return false
}

// End GameObjectStackable interface

// Movable interface BEGIN

func (p *Player) GetMovement() int {
	if p.Flying {
		return 6 // FIXME - is this right?
	}
	return 1
}

func (p *Player) IsFlying() bool {
	return p.Flying
}

func (p *Player) CheckBelongsTo(player *Player) bool {
	return player == p
}

func (p *Player) SetBoardPosition(v logical.Vec) {
	p.BoardPosition = v
}

func (p *Player) GetBoardPosition() logical.Vec {
	return p.BoardPosition
}

func (p *Player) BreakEngagement() bool {
	if rand.IntN(9) >= p.Manoeuvre {
		return true
	}
	return false
}

// Movable interface END

// Attackable interface BEGIN

func (p *Player) GetDefence() int {
	return p.Defence
}

func (p *Player) Engageable() bool {
	return true
}

// SetBoardPosition is in GameObject interface also

// Attackable interface END

// Attackerable interface BEGIN

func (p *Player) GetCombat() int {
	return p.Combat
}

// Attackerable interface END

func (p *Player) CastSpell(target logical.Vec, grid *grid.GameGrid) bool {
	fmt.Printf("IN Player spell cast\n")
	i := p.ChosenSpell
	spell := p.Spells[i]
	if !spell.IsReuseable() {
		spells := make([]spells.Spell, 0)
		spells = append(p.Spells[:i], p.Spells[i+1:]...)
		p.Spells = spells
	}
	ret := spell.DoesCastWork(p.LawRating)
	if ret {
		fmt.Printf("Player spell %T cast on %T\n", spell, target)
		spell.Cast(target, grid, p)
	}
	p.ChosenSpell = -1
	return ret
}

// Player Spells (spells which only affect a player)

type PlayerSpell struct {
	spells.ASpell
	MutateFunc func(*Player)
}

func (s PlayerSpell) CanCast(target grid.GameObject) bool {
	_, ok := target.(*Player)
	if ok {
		return true
	}
	return false
}

func (s PlayerSpell) Cast(target logical.Vec, grid *grid.GameGrid, castor grid.GameObject) {
	tile := grid.GetGameObject(target)
	player := tile.(*Player)
	s.MutateFunc(player)
	// May have just become not animated
	player.SpriteIdx = 0
}

func init() {
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Magic Knife",
			LawRating:     1,
			CastingChance: 90,
		},
		MutateFunc: func(p *Player) {
			p.CharacterIcon = logical.V(4, 22)
			p.IsAnimated = true
		},
	})
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Magic Armour",
			LawRating:     1,
			CastingChance: 40,
		},
		MutateFunc: func(p *Player) {
			p.CharacterIcon = logical.V(1, 20)
			p.IsAnimated = false
		},
	})
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Magic Shield",
			LawRating:     1,
			CastingChance: 80,
		},
		MutateFunc: func(p *Player) {
			p.CharacterIcon = logical.V(0, 20)
			p.IsAnimated = false
		},
	})
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Magic Sword",
			LawRating:     1,  // FIXME
			CastingChance: 80, // FIXME
		},
		MutateFunc: func(p *Player) {
			p.CharacterIcon = logical.V(0, 21)
			p.IsAnimated = true
		},
	})
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Magic Bow",
			LawRating:     1,
			CastingChance: 50,
		},
		MutateFunc: func(p *Player) {
			p.CharacterIcon = logical.V(0, 22)
			p.IsAnimated = true
		},
	})
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Shadow Form",
			CastingChance: 80,
		},
		MutateFunc: func(p *Player) {
			// FIXME
		},
	})
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Magic Wings",
			CastingChance: 60,
		},
		MutateFunc: func(p *Player) {
			p.CharacterIcon = logical.V(4, 21)
			p.IsAnimated = true
		},
	})
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Law-1",
			CastingChance: 100,
			LawRating:     2,
			NoCastFx:      true,
		},
		MutateFunc: func(p *Player) {
			p.LawRating++
		},
	})
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Law-2",
			CastingChance: 60,
			LawRating:     4,
			NoCastFx:      true,
		},
		MutateFunc: func(p *Player) {
			p.LawRating = p.LawRating + 2
		},
	})
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Chaos-1",
			CastingChance: 80,
			LawRating:     -2,
			NoCastFx:      true,
		},
		MutateFunc: func(p *Player) {
			p.LawRating--
		},
	})
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Chaos-2",
			CastingChance: 60,
			LawRating:     -4,
			NoCastFx:      true,
		},
		MutateFunc: func(p *Player) {
			p.LawRating = p.LawRating - 2
		},
	})
}
