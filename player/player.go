package player

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/spells"
)

func NewPlayer() Player {
	return Player{
		Defence:         rand.Intn(4) + 1, // 1-5
		Combat:          rand.Intn(4) + 1, // 1-5
		Manoeuvre:       rand.Intn(4) + 3, // 3-7
		MagicResistance: rand.Intn(2) + 6, // 6-8
		Movement:        1,
		ChosenSpell:     -1,
		Spells:          spells.ChooseSpells(),
		Alive:           true,
	}
}

type Player struct {
	Name            string
	Spells          []spells.Spell
	HumanPlayer     bool
	CharacterIcon   logical.Vec
	ChosenSpell     int
	CastIllusion    bool
	Color           color.Color
	LawRating       int
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
	if rand.Intn(9) >= p.Manoeuvre {
		return true
	}
	return false
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
	return fx.FxRemoteAttack()
}

func (p *Player) CanAttackUndead() bool {
	return p.HasMagicWeapon
}

// Attackerable interface END

func (p *Player) CastSpell(target logical.Vec, grid *grid.GameGrid) (bool, *fx.Fx) {
	var anim *fx.Fx
	fmt.Printf("IN Player spell cast\n")
	i := p.ChosenSpell
	spell := p.Spells[i]
	if !spell.IsReuseable() {
		spells := make([]spells.Spell, 0)
		spells = append(p.Spells[:i], p.Spells[i+1:]...)
		p.Spells = spells
	}
	fmt.Printf("Player spell %T cast on %T\n", spell, target)
	ret, anim := spell.Cast(p.CastIllusion, p.LawRating, target, grid, p)
	p.ChosenSpell = -1
	return ret, anim
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

func (s PlayerSpell) Cast(illusion bool, playerLawRating int, target logical.Vec, grid *grid.GameGrid, castor grid.GameObject) (bool, *fx.Fx) {
	if illusion {
		panic("PlayerSpell cannot be illusion")
	}
	if rand.Intn(100) <= s.GetCastingChance(playerLawRating) {
		tile := grid.GetGameObject(target)
		player := tile.(*Player)
		s.MutateFunc(player)
		// May have just become not animated
		player.SpriteIdx = 0
		return true, nil
	}
	return false, nil
}

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
			p.Defence = p.Defence + 4
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
			p.Defence = p.Defence + 2
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
			p.Combat = p.Combat + 2
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
			p.Combat = p.Combat + 4
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
