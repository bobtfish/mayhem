package movable

import (
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
)

type Movable interface {
	GetMovement() int
	IsFlying() bool
	CheckBelongsTo(*player.Player) bool
	GetBoardPosition() logical.Vec
	SetBoardPosition(logical.Vec)
}

type Attackable interface {
	GetDefence() int
	CheckBelongsTo(*player.Player) bool
}

type Corpseable interface {
	CanMakeCorpse() bool
	MakeCorpse()
}
