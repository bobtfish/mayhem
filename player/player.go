package player

import (
	"fmt"
	"image/color"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/rand"
	spelliface "github.com/bobtfish/mayhem/spells/iface"
)

func NewPlayer() Player {
	return Player{
		Defence:         rand.Intn(4) + 1, // 1-5
		Combat:          rand.Intn(4) + 1, // 1-5
		Manoeuvre:       rand.Intn(4) + 3, // 3-7
		MagicResistance: rand.Intn(2) + 6, // 6-8
		Movement:        1,
		ChosenSpell:     -1,
		Spells:          spelliface.ChooseSpells(),
		Alive:           true,
	}
}

type Player struct {
	Name            string
	Spells          []spelliface.Spell
	HumanPlayer     bool
	CharacterIcon   logical.Vec
	ChosenSpell     int
	CastIllusion    bool
	Color           color.Color
	BoardPosition   logical.Vec
	IsAnimated      bool
	SpriteIdx       int
	Defence         int
	Combat          int
	RangedCombat    int
	AttackRange     int
	Manoeuvre       int
	MagicResistance int
	Movement        int

	Flying         bool // If the player has magic wings
	HasMagicWeapon bool // If the player can attack undead

	Alive bool
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
	fmt.Printf("Get player movement, is %d\n", p.Movement)
	return p.Movement
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
	return rand.Intn(9) >= p.Manoeuvre
}

// Movable interface END

// Attackable interface BEGIN

func (p *Player) GetDefence() int {
	return p.Defence
}

func (p *Player) GetMagicResistance() int {
	return p.MagicResistance
}

func (p *Player) Engageable() bool {
	return true
}

// SetBoardPosition is in GameObject interface also

func (p *Player) IsUndead() bool {
	return false
}

func (p *Player) IsMount() bool {
	return false
}

// Attackable interface END

// Attackerable interface BEGIN

func (p *Player) GetCombat() int {
	return p.Combat
}

func (p *Player) GetRangedCombat() int {
	return p.RangedCombat
}

func (p *Player) GetAttackRange() int {
	return p.AttackRange
}

func (p *Player) GetAttackFx() *fx.Fx {
	return fx.RemoteAttack()
}

func (p *Player) CanAttackUndead() bool {
	return p.HasMagicWeapon
}

// Attackerable interface END
