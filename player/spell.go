package player

// Player Spells (spells which only affect a player)

import (
	"fmt"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/spells"
)

type PlayerSpell struct {
	spells.ASpell
	MutateFunc func(*Player)
}

func (s PlayerSpell) CanCast(target grid.GameObject) bool {
	_, ok := target.(*Player)
	return ok
}

func (s PlayerSpell) DoCast(illusion bool, target logical.Vec, grid *grid.GameGrid, castor grid.GameObject) (bool, *fx.Fx) {
	if illusion {
		panic("PlayerSpells should never be illusions")
	}
	tile := grid.GetGameObject(target)
	player, isPlayer := tile.(*Player)
	if !isPlayer {
		// FIXME - bug here if a player is mounted!
		panic(fmt.Sprintf("Player spell '%s' cast on non player - should never happen", s.Name))
		/*rideable, isRideable := tile.(*movable.Rideable)
		if !isRideable
			panic(fmt.Sprintf("Player spell '%s' cast on non player - should never happen", s.Name))
		}
		player = rideable.GetRider();
		if player == nil {
			panic(fmt.Sprintf("Player spell '%s' cast on non player - should never happen", s.Name))
		}*/
	}
	s.MutateFunc(player)
	return true, nil
}

// Setup all the player spells.
func init() {
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Magic Armour",
			LawRating:     1,
			CastingChance: 40,
		},
		MutateFunc: func(p *Player) {
			p.CharacterIcon = logical.V(1, 20)
			p.IsAnimated = false
			p.SpriteIdx = 0
			p.Defence += 4
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
			p.SpriteIdx = 0
			p.Defence += 2
		},
	})
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Magic Knife",
			LawRating:     1,
			CastingChance: 90,
		},
		MutateFunc: func(p *Player) {
			p.CharacterIcon = logical.V(4, 22)
			p.IsAnimated = true
			p.Combat += 2
			p.HasMagicWeapon = true
		},
	})
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Magic Sword",
			LawRating:     1,
			CastingChance: 50,
		},
		MutateFunc: func(p *Player) {
			p.CharacterIcon = logical.V(0, 21)
			p.IsAnimated = true
			p.Combat += 4
			p.HasMagicWeapon = true
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
			p.AttackRange = 6
			p.RangedCombat = 3
		},
	})
	spells.CreateSpell(PlayerSpell{
		ASpell: spells.ASpell{
			Name:          "Shadow Form",
			CastingChance: 80,
		},
		MutateFunc: func(p *Player) {
			// FIXME - animation
			if p.Movement < 3 {
				p.Movement = 3
			}
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
			p.Flying = true
			p.Movement = 6
		},
	})
}
