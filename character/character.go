package character

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"image/color"
	"math/rand"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/render"
	"github.com/bobtfish/mayhem/spells"
)

// Abstract character that can be created
type CharacterType struct {
	Name               string  `yaml:"name"`
	Combat             int     `yaml:"combat"`
	RangedCombat       int     `yaml:"ranged_combat"`
	AttackRange        int     `yaml:"range"`
	RangedAttackIsFire bool    `yaml:"breathes_fire"`
	Defence            int     `yaml:"defence"`
	Movement           int     `yaml:"movement"`
	Flying             bool    `yaml:"flying"`
	MagicalResistance  int     `yaml:"magical_resistance"`
	Manoeuvre          int     `yaml:"manoeuvre"`
	Unknown            int     `yaml:"unknown"`
	LawChaos           int     `yaml:"law_chaos"`
	Strength           int     `yaml:"strength"`
	Sprites            [][]int `yaml:"sprites"`
	DeadSprite         []int   `yaml:"deadsprite"`
	ColorR             int     `yaml:"color_r"`
	ColorG             int     `yaml:"color_g"`
	ColorB             int     `yaml:"color_b"`
	Undead             bool    `yaml:"undead"`
	CastRange          int     `yaml:"cast_range"`
	CanBeIllusion      bool    `yaml:"can_be_illusion"`
}

func LoadCharacterTemplates() {
	cl := make([]CharacterType, 0)
	err := yaml.Unmarshal([]byte(character_yaml), &cl)
	if err != nil {
		panic(err)
	}
	for _, v := range cl {
		if v.Sprites == nil {
			v.Sprites = make([][]int, 0)
		}
		if v.DeadSprite == nil {
			v.DeadSprite = make([]int, 2)
			v.DeadSprite[0] = 0
			v.DeadSprite[1] = 0
		}
		if v.CastRange == 0 {
			v.CastRange = 1
		}
		//fmt.Printf("Create %s range %d\n", v.Name, castRange)
		spells.CreateSpell(CharacterSpell{
			Name:               v.Name,
			LawRating:          v.LawChaos,
			CastingChance:      100, // FIXME
			Sprite:             logical.V(v.Sprites[0][0], v.Sprites[0][1]),
			Color:              render.GetColor(v.ColorR, v.ColorG, v.ColorB),
			Movement:           v.Movement,
			Flying:             v.Flying,
			Undead:             v.Undead,
			CastRange:          v.CastRange,
			Defence:            v.Defence,
			DeadSprite:         logical.V(v.DeadSprite[0], v.DeadSprite[1]),
			Combat:             v.Combat,
			Manoeuvre:          v.Manoeuvre,
			RangedCombat:       v.RangedCombat,
			AttackRange:        v.AttackRange,
			RangedAttackIsFire: v.RangedAttackIsFire,
			CanBeIllusion:      v.CanBeIllusion,
		})
	}

	// We know that the spells array is initialised now, add the disbelieve spell
	// This is done here as character depends on spells, and so we can't have spell depend on character
	spells.AllSpells[0] = DisbelieveSpell{spells.ASpell{
		Name:          "Disbelieve",
		LawRating:     0,
		Reuseable:     true,
		CastingChance: 100,
		CastRange:     20,
	}}
}

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

// This is the spell to create a character
type CharacterSpell struct {
	Name               string
	LawRating          int
	CastingChance      int
	CastRange          int
	Sprite             logical.Vec
	DeadSprite         logical.Vec
	Color              color.Color
	Movement           int
	Flying             bool
	Undead             bool
	Defence            int
	Combat             int
	RangedCombat       int
	AttackRange        int
	RangedAttackIsFire bool
	Manoeuvre          int
	CanBeIllusion      bool
}

// Spell interface begin
func (s CharacterSpell) GetName() string {
	return s.Name
}
func (s CharacterSpell) GetLawRating() int {
	return s.LawRating
}
func (s CharacterSpell) GetCastingChance(playerLawRating int) int {
	// FIXME do something with playerLawRating
	return s.CastingChance
}
func (s CharacterSpell) GetCastRange() int {
	return s.CastRange
}

func (s CharacterSpell) CanCast(target grid.GameObject) bool {
	if target.IsEmpty() {
		return true
	}
	return false
}

func (s CharacterSpell) Cast(illusion bool, playerLawRating int, target logical.Vec, grid *grid.GameGrid, castor grid.GameObject) (bool, *fx.Fx) {
	if rand.Intn(100) <= s.GetCastingChance(playerLawRating) {
		grid.PlaceGameObject(target, s.CreateCharacter(illusion, castor))
		return true, nil
	}
	return false, nil
}

func (s CharacterSpell) IsReuseable() bool {
	return false
}

func (s CharacterSpell) CastFx() *fx.Fx {
	return fx.FxSpellCast()
}

func (s CharacterSpell) CanCastAsIllusion() bool {
	return s.CanBeIllusion
}

// Spell interface end

func (s CharacterSpell) CreateCharacter(isIllusion bool, castor grid.GameObject) *Character {
	return &Character{
		Name:               s.Name,
		Sprite:             s.Sprite,
		Color:              s.Color,
		Movement:           s.Movement,
		Flying:             s.Flying,
		Undead:             s.Undead,
		Defence:            s.Defence,
		DeadSprite:         s.DeadSprite,
		Combat:             s.Combat,
		Manoeuvre:          s.Manoeuvre,
		RangedCombat:       s.RangedCombat,
		RangedAttackIsFire: s.RangedAttackIsFire,
		AttackRange:        s.AttackRange,
		IsIllusion:         isIllusion,

		// FIXME - ugh this is gross - would it be better done up a level?
		BelongsTo: castor.(*player.Player),
	}
}

// This is the actual character that gets created
type Character struct {
	Name               string
	Sprite             logical.Vec
	Color              color.Color
	Movement           int
	Flying             bool
	Undead             bool
	Defence            int
	Combat             int
	RangedCombat       int
	AttackRange        int
	RangedAttackIsFire bool
	Manoeuvre          int
	DeadSprite         logical.Vec
	IsDead             bool
	IsIllusion         bool

	BelongsTo *player.Player
	// Remember to add any fields you add here to the Clone method

	SpriteIdx     int
	BoardPosition logical.Vec
}

func (c *Character) Clone() *Character {
	return &Character{
		Name:               c.Name,
		Sprite:             c.Sprite,
		Color:              c.Color,
		Movement:           c.Movement,
		Flying:             c.Flying,
		Undead:             c.Undead,
		Defence:            c.Defence,
		Combat:             c.Combat,
		AttackRange:        c.AttackRange,
		RangedCombat:       c.RangedCombat,
		RangedAttackIsFire: c.RangedAttackIsFire,
		Manoeuvre:          c.Manoeuvre,
		DeadSprite:         c.DeadSprite,
		IsDead:             c.IsDead,
		BelongsTo:          c.BelongsTo,
		IsIllusion:         c.IsIllusion,
	}
}

// GameObject + GameObjectStackable interface BEGIN
func (c *Character) AnimationTick(odd bool) {
	if odd {
		return
	}
	c.SpriteIdx++
	if c.SpriteIdx == 4 {
		c.SpriteIdx = 0
	}
	return
}

func (c *Character) IsEmpty() bool {
	if c.IsDead {
		return true
	}
	return false
}

func (c *Character) GetSpriteSheetCoordinates() logical.Vec {
	if c.IsDead {
		return c.DeadSprite
	}
	return c.Sprite.Add(logical.V(c.SpriteIdx, 0))
}

func (c *Character) GetColor() color.Color {
	return c.Color
}

func (c *Character) Describe() string {
	return fmt.Sprintf("%s (%s)", c.Name, c.BelongsTo.Name)
}

func (c *Character) SetBoardPosition(v logical.Vec) {
	c.BoardPosition = v
}

// GameObject interface END

func (c *Character) RemoveMe() bool {
	return false // FIXME - what about if destroyed
}

// GameObjectStackable interface END

// Movable interface BEGIN

func (c *Character) GetMovement() int {
	if c.IsDead {
		return 0
	}
	return c.Movement
}

func (c *Character) IsFlying() bool {
	return c.Flying
}

func (c *Character) CheckBelongsTo(player *player.Player) bool {
	return player == c.BelongsTo
}

func (c *Character) BreakEngagement() bool {
	if rand.Intn(9) >= c.Manoeuvre {
		return true
	}
	return false
}

// SetBoardPosition is in GameObject interface also

func (c *Character) GetBoardPosition() logical.Vec {
	return c.BoardPosition
}

// Movable interface END

// Attackable interface BEGIN

func (c *Character) GetDefence() int {
	return c.Defence
}

func (c *Character) Engageable() bool {
	if c.Movement > 0 {
		return true
	}
	return false
}

// SetBoardPosition is in GameObject interface also

func (c *Character) IsUndead() bool {
	return c.Undead
}

// Attackable interface END

// Attackerable interface BEGIN

func (c *Character) GetCombat() int {
	return c.Combat
}

func (c *Character) GetRangedCombat() int {
	return c.RangedCombat
}

func (c *Character) GetAttackRange() int {
	return c.AttackRange
}

func (c *Character) GetAttackFx() *fx.Fx {
	if c.RangedAttackIsFire {
		return fx.FxFire()
	}
	return fx.FxRemoteAttack()
}

func (c *Character) CanAttackUndead() bool {
	return c.Undead
}

// Attackerable interface END

// Corpseable interface BEGIN

func (c *Character) CanMakeCorpse() bool {
	if c.Undead {
		return false
	}
	if c.DeadSprite.Equals(logical.ZeroVec()) {
		return false
	}
	return true
}

func (c *Character) MakeCorpse() {
	c.IsDead = true
}

// Corpsable interface END
