package otherspells

import (
	"fmt"

	"github.com/bobtfish/mayhem/character"
	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/movable"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/rand"
	"github.com/bobtfish/mayhem/spells"
	spelliface "github.com/bobtfish/mayhem/spells/iface"
)

type OtherSpell struct {
	spells.ASpell
	MutateFunc func(logical.Vec, *grid.GameGrid, grid.GameObject) (bool, *fx.Fx)
}

func (s OtherSpell) DoCast(illusion bool, target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
	return s.MutateFunc(target, grid, owner)
}

func init() {
	spelliface.CreateSpell(OtherSpell{
		ASpell: spells.ASpell{
			Name:          "Raise Dead",
			LawRating:     -1,
			CastingChance: 60,
			CastRange:     4,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
			p := owner.(*player.Player)
			ob := grid.GetGameObject(target)
			char, isChar := ob.(*character.Character)
			if !isChar || !char.IsDead {
				return false, nil
			}
			char.IsDead = false
			char.Undead = true
			char.BelongsTo = p
			return true, nil
		},
	})
	spelliface.CreateSpell(OtherSpell{
		ASpell: spells.ASpell{
			Name:          "Subversion",
			CastingChance: 100,
			CastRange:     7,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
			p := owner.(*player.Player)
			ob := grid.GetGameObject(target)
			char, isChar := ob.(*character.Character)
			if !isChar || char.IsDead {
				return false, nil
			}
			if rand.Intn(9)+1 > char.MagicResistance {
				char.BelongsTo = p
				return true, nil
			}
			return false, nil
		},
	})
	spelliface.CreateSpell(OtherSpell{
		ASpell: spells.ASpell{ // 1 chance only, makes creatures belonging to player explode
			Name:                "Vengeance",
			CastingChance:       80,
			CastRange:           20,
			NoLineOfSightNeeded: true,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
			return ExplodeCreatures(target, grid)
		},
	})
	spelliface.CreateSpell(OtherSpell{
		ASpell: spells.ASpell{ // 1 chance only, doesn't kill player - makes their creatures explode maybe?
			Name:                "Decree",
			CastingChance:       80,
			CastRange:           20,
			LawRating:           1,
			NoLineOfSightNeeded: true,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
			return ExplodeCreatures(target, grid)
		},
	})
	spelliface.CreateSpell(OtherSpell{
		ASpell: spells.ASpell{ // 3 tries, doesn't kill player - makes their creatures explode
			Name:                "Dark Power",
			CastingChance:       50,
			CastRange:           20,
			LawRating:           -2,
			Tries:               3,
			NoLineOfSightNeeded: true,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
			return ExplodeCreatures(target, grid)
		},
	})
	spelliface.CreateSpell(OtherSpell{
		ASpell: spells.ASpell{ // 3 tries, doesn't kill player - makes their creatures explode
			Name:                "Justice",
			CastingChance:       50,
			CastRange:           20,
			LawRating:           2,
			Tries:               3,
			NoLineOfSightNeeded: true,
		},
		MutateFunc: func(target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
			return ExplodeCreatures(target, grid)
		},
	})
}

func ExplodeCreatures(target logical.Vec, grid *grid.GameGrid) (bool, *fx.Fx) {
	ob := grid.GetGameObject(target)
	a, isAttackable := ob.(movable.Attackable)
	if !isAttackable {
		return false, nil
	}
	chance := rand.Intn(9)
	fmt.Printf("Chance %d > Resistance %d\n", chance, a.GetMagicResistance())
	if chance > a.GetMagicResistance() {
		player, isPlayer := ob.(*player.Player)
		f := fx.Disbelieve()
		if isPlayer {
			// Loop through the board and explode every character belonging to this player
			for x := 0; x < grid.Width(); x++ {
				for y := 0; y < grid.Height(); y++ {
					vec := logical.V(x, y)
					if target.Equals(vec) {
						grid.PlaceGameObject(vec, f)
					} else {
						otherA, otherIsAttackable := grid.GetGameObject(vec).(movable.Attackable)
						if otherIsAttackable {
							if otherA.CheckBelongsTo(player) {
								grid.GetGameObjectStack(vec).RemoveTopObject()
								grid.PlaceGameObject(vec, fx.Disbelieve())
							}
						}
					}
				}
			}
		} else {
			// Just explode this character
			grid.GetGameObjectStack(target).RemoveTopObject()
			grid.PlaceGameObject(target, f)
		}
		return true, f
	}
	return false, nil
}
