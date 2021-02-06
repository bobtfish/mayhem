package character

import (
	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/spells"
)

type DisbelieveSpell struct {
	spells.ASpell
}

func (s DisbelieveSpell) DoCast(illusion bool, target logical.Vec, grid *grid.GameGrid, owner grid.GameObject) (bool, *fx.Fx) {
	if illusion {
		panic("DisbelieveSpell cannot be illusion")
	}
	character := grid.GetGameObject(target).(*Character)
	if character.IsIllusion {
		grid.GetGameObjectStack(target).RemoveTopObject()
		anim := fx.Disbelieve()
		grid.PlaceGameObject(target, anim)
		return true, anim
	}
	return false, nil
}

func (s DisbelieveSpell) CanCast(target grid.GameObject) bool {
	_, isCharacter := target.(*Character)
	return isCharacter
}

func (s DisbelieveSpell) NeedsLineOfSight() bool {
	return false
}
