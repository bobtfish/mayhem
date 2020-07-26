package character

import (
	"math/rand"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/spells"
)

type DisbelieveSpell struct {
	spells.ASpell
}

func (s DisbelieveSpell) Cast(illusion bool, playerLawRating int, target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
	if illusion {
		panic("DisbelieveSpell cannot be illusion")
	}
	character := grid.GetGameObject(target).(*Character)
	if character.IsIllusion {
		grid.GetGameObjectStack(target).RemoveTopObject()
		anim := fx.FxDisbelieve()
		grid.PlaceGameObject(target, anim)
		return true, anim
	}
	return false, nil
}

func (s DisbelieveSpell) CanCast(target grid.GameObject) bool {
	_, isCharacter := target.(*Character)
	if isCharacter {
		return true
	}
	return false
}

func init() {
	spells.CreateSpell(spells.OtherSpell{
		ASpell: spells.ASpell{
			Name:          "Raise Dead",
			LawRating:     -1,
			CastingChance: 60,
			CastRange:     4,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) bool {
			p := owner.(*player.Player)
			ob := grid.GetGameObject(target)
			char, isChar := ob.(*Character)
			if !isChar || !char.IsDead {
				return false
			}
			char.IsDead = false
			char.Undead = true
			char.BelongsTo = p
			return true
		},
	})
	spells.CreateSpell(spells.OtherSpell{
		ASpell: spells.ASpell{
			Name:          "Subversion",
			CastingChance: 100,
			CastRange:     7,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) bool {
			p := owner.(*player.Player)
			ob := grid.GetGameObject(target)
			char, isChar := ob.(*Character)
			if !isChar || char.IsDead {
				return false
			}
			if rand.Intn(9)+1 > char.MagicResistance {
				char.BelongsTo = p
				return true
			}
			return false
		},
	})
	/*spells.CreateSpell(spells.OtherSpell{
		ASpell: spells.ASpell{ // 1 chance only, makes creatures belonging to player explode
			Name:          "Vengence",
			CastingChance: 80,
			CastRange:     20,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) bool {
			return ExplodeCreatures(target, grid)
		},
	})
	spells.CreateSpell(spells.OtherSpell{
		ASpell: spells.ASpell{ // 1 chance only, doesn't kill player - makes their creatures explode maybe?
			Name:          "Decree",
			CastingChance: 80,
			CastRange:     20,
			LawRating:     1,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) bool {
			return ExplodeCreatures(target, grid)
		},
	})
	spells.CreateSpell(spells.OtherSpell{
		ASpell: spells.ASpell{ // 3 tries, doesn't kill player - makes their creatures explode
			Name:          "Dark Power",
			CastingChance: 50,
			CastRange:     20,
			LawRating:     -2,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) bool {
			return ExplodeCreatures(target, grid)
		},
	})
    */
	spells.CreateSpell(spells.OtherSpell{
		ASpell: spells.ASpell{ // 3 tries, doesn't kill player - makes their creatures explode
			Name:          "Justice",
			CastingChance: 50,
			CastRange:     20,
			LawRating:     2,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) bool {
			return ExplodeCreatures(target, grid)
		},
	})
}
