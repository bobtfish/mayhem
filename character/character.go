package character

import (
	"fmt"
	"image/color"

	"gopkg.in/yaml.v2"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/rand"
	"github.com/bobtfish/mayhem/render"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
	"github.com/bobtfish/mayhem/spells"
	spelliface "github.com/bobtfish/mayhem/spells/iface"
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
	MagicResistance    int     `yaml:"magical_resistance"`
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
	BaseCastingChance  int     `yaml:"base_casting_chance"`
	Mountable          bool    `yaml:"mountable"`
	Quantity           int     `yaml:"quantity"`
}

func LoadCharacterTemplates() {
	cl := make([]CharacterType, 0)
	err := yaml.Unmarshal([]byte(characterYaml), &cl)
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
		if v.BaseCastingChance == 0 {
			fmt.Printf("Character spell %s has no 'base_casting_chance', setting chance to 100\n", v.Name)
			v.BaseCastingChance = 100
		}
		spelliface.CreateSpell(CharacterSpell{
			Name:               v.Name,
			LawRating:          v.LawChaos,
			CastingChance:      v.BaseCastingChance,
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
			Mountable:          v.Mountable,
			MagicResistance:    v.MagicResistance,
			Quantity:           v.Quantity,
		})
	}

	// We know that the spells array is initialised now, add the disbelieve spell
	// This is done here as character depends on spells, and so we can't have spell depend on character
	spelliface.AllSpells[0] = DisbelieveSpell{spells.ASpell{
		Name:          "Disbelieve",
		LawRating:     0,
		Reuseable:     true,
		CastingChance: 100,
		CastRange:     20,
	}}
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
	Mountable          bool
	MagicResistance    int
	Quantity           int
}

// Spell interface begin
func (s CharacterSpell) TakeOverScreen(ctx screeniface.GameCtx, cleanupFunc func(), nextScreen screeniface.GameScreen, playeerIdx int, target logical.Vec) screeniface.GameScreen {
	return nil
}
func (s CharacterSpell) GetName() string {
	return s.Name
}
func (s CharacterSpell) GetLawRating() int {
	return s.LawRating
}
func (s CharacterSpell) GetCastingChance(lawRating int) int {
	// FIXME do something with lawRating
	return s.CastingChance
}
func (s CharacterSpell) GetCastRange() int {
	return s.CastRange
}

func (s CharacterSpell) CanCast(target grid.GameObject) bool {
	return target.IsEmpty()
}

func (s CharacterSpell) CastQuantity() int {
	if s.Quantity == 0 {
		return 1
	}
	return s.Quantity
}

func (s CharacterSpell) DoCast(illusion bool, target logical.Vec, grid *grid.GameGrid, castor grid.GameObject) (bool, *fx.Fx) {
	grid.PlaceGameObject(target, s.CreateCharacter(illusion, castor))
	return true, nil
}

// FIXME duplicate code with spells/main.go
func (s CharacterSpell) CastSucceeds(illusion bool, playerLawRating int) bool {
	if illusion && !s.CanCastAsIllusion() {
		panic(fmt.Sprintf("Spell %s (type %T) cannot be illusion, but was cast as one anyway", s.Name, s))
	}
	if illusion || rand.Intn(100) <= s.GetCastingChance(playerLawRating) {
		return true
	}
	return false
}

func (s CharacterSpell) IsReuseable() bool {
	return false
}

func (s CharacterSpell) CastFx() *fx.Fx {
	return fx.SpellCast()
}

func (s CharacterSpell) CanCastAsIllusion() bool {
	return s.CanBeIllusion
}

func (s CharacterSpell) NeedsLineOfSight() bool {
	return true
}

func (s CharacterSpell) GetDescriptionArray(lawRating int) []string {
	desc := make([]string, 0)
	desc = append(desc, fmt.Sprintf("Combat=%d", s.Combat))
	desc = append(desc, fmt.Sprintf("Ranged Combat=%d Range=%d", s.RangedCombat, s.AttackRange))
	desc = append(desc, fmt.Sprintf("Defence=%d", s.Defence))
	desc = append(desc, fmt.Sprintf("Movement Allowance=%d", s.Movement))
	desc = append(desc, fmt.Sprintf("Manoeuver Rating=%d", s.Manoeuvre))
	desc = append(desc, fmt.Sprintf("Magic Resistance=%d", s.MagicResistance))
	desc = append(desc, fmt.Sprintf("Casting Chance=%d%%", s.GetCastingChance(lawRating)))
	return desc
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
		Mountable:          s.Mountable,
		MagicResistance:    s.MagicResistance,

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
	Mountable          bool
	MagicResistance    int

	BelongsTo      *player.Player
	CarryingPlayer bool
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
		Mountable:          c.Mountable,
		MagicResistance:    c.MagicResistance,
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
}

func (c *Character) IsEmpty() bool {
	return c.IsDead
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

func (c *Character) Describe() (string, string) {
	if c.BelongsTo != nil {
		carry := ""
		if c.CarryingPlayer {
			carry = "#"
		}
		return fmt.Sprintf("%s%s", c.Name, carry), c.BelongsTo.Name
	}
	if c.IsDead {
		return fmt.Sprintf("%s (Dead)", c.Name), ""
	}
	return c.Name, ""
}

func (c *Character) SetBoardPosition(v logical.Vec) {
	c.BoardPosition = v
	if c.CarryingPlayer {
		c.BelongsTo.BoardPosition = v
	}
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
	return rand.Intn(9) >= c.Manoeuvre
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

func (c *Character) GetMagicResistance() int {
	return c.MagicResistance
}

func (c *Character) Engageable() bool {
	return c.Movement > 0
}

// SetBoardPosition is in GameObject interface also

func (c *Character) IsUndead() bool {
	return c.Undead
}

func (c *Character) IsMount() bool {
	return c.Mountable
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
		return fx.Fire()
	}
	return fx.RemoteAttack()
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
	c.BelongsTo = nil
}

// Corpsable interface END

func (c *Character) Mount() {
	if !c.IsMount() {
		panic("Tried to mount player on character which is not a mount")
	}
	fmt.Printf("Mount the player\n")
	c.CarryingPlayer = true
}
