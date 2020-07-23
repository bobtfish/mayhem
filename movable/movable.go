package movable

import (
	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
)

type Movable interface {
	GetMovement() int
	IsFlying() bool
	CheckBelongsTo(*player.Player) bool
	GetBoardPosition() logical.Vec
	SetBoardPosition(logical.Vec)
	BreakEngagement() bool
}

type Attackable interface {
	GetDefence() int
	CheckBelongsTo(*player.Player) bool
	Engageable() bool
	IsUndead() bool
	IsMount() bool
}

type Attackerable interface {
	GetCombat() int
	GetRangedCombat() int
	GetAttackRange() int
	GetAttackFx() *fx.Fx
	CanAttackUndead() bool
}

type Corpseable interface {
	CanMakeCorpse() bool
	MakeCorpse()
	GetBoardPosition() logical.Vec
}
