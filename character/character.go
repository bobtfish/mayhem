package character

import (
	"gopkg.in/yaml.v2"
	"image/color"
	"math/rand"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
	"github.com/bobtfish/mayhem/spells"
)

// Abstract character that can be created
type CharacterType struct {
	Name              string  `yaml:"name"`
	Combat            int     `yaml:"combat"`
	RangedCombat      int     `yaml:"ranged_combat"`
	AttackRange       int     `yaml:"range"`
	Defence           int     `yaml:"defence"`
	Movement          int     `yaml:"movement"`
	MagicalResistance int     `yaml:"magical_resistance"`
	Manoeuvre         int     `yaml:"manoeuvre"`
	Unknown           int     `yaml:"unknown"`
	LawChaos          int     `yaml:"law_chaos"`
	Strength          int     `yaml:"strength"`
	Sprites           [][]int `yaml:"sprites"`
	DeadSprite        []int   `yaml:"deadsprite"`
	ColorR            int     `yaml:"color_r"`
	ColorG            int     `yaml:"color_g"`
	ColorB            int     `yaml:"color_b"`
}

// Individual character instance
/*type Character struct {
	CharacterSpell
	SpriteIdx     int
	BoardPosition logical.Vec
}*/

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
			v.DeadSprite = make([]int, 0)
		}
		castRange := 1
		//fmt.Printf("Create %s range %d\n", v.Name, castRange)
		spells.CreateSpell(CharacterSpell{
			Name:          v.Name,
			LawRating:     v.LawChaos,
			CastingChance: 100, // FIXME
			Range:         castRange,
			Sprite:        logical.V(v.Sprites[0][0], v.Sprites[0][1]),
			Color:         render.GetColor(v.ColorR, v.ColorG, v.ColorB),
		})
	}
}

// This is the spell to create a character
type CharacterSpell struct {
	Name          string
	LawRating     int
	CastingChance int
	Range         int
	Sprite        logical.Vec
	Color         color.Color
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
func (s CharacterSpell) GetRange() int {
	return s.Range
}

// FIXME this is duplicate code
func (s CharacterSpell) DoesCastWork(playerLawRating int) bool {
	// FIXME
	return true
	if rand.Intn(100) <= s.GetCastingChance(playerLawRating) {
		return true
	}
	return false
}

func (s CharacterSpell) CanCast(target grid.GameObject) bool {
	if target.IsEmpty() {
		return true
	}
	return false
}

func (s CharacterSpell) Cast(target logical.Vec, grid *grid.GameGrid) {
	grid.PlaceGameObject(target, s.CreateCharacter())
}

func (s CharacterSpell) IsReuseable() bool {
	return false
}

func (s CharacterSpell) CastFx() *fx.Fx {
	return fx.FxSpellCast()
}

// Spell interface end

func (s CharacterSpell) CreateCharacter() *Character {
	return &Character{
		Name:   s.Name,
		Sprite: s.Sprite,
		Color:  s.Color,
	}
}

// This is the actual character that gets created
type Character struct {
	Name   string
	Sprite logical.Vec
	Color  color.Color
	// Remember to add any fields you add here to the Clone method

	SpriteIdx     int
	BoardPosition logical.Vec
}

// GameObject interface BEGIN
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

func (c *Character) Clone() *Character {
	return &Character{
		Name:   c.Name,
		Sprite: c.Sprite,
		Color:  c.Color,
	}
}

func (c *Character) RemoveMe() bool {
	return false // FIXME - what about if destroyed
}

func (c *Character) IsEmpty() bool {
	return false
}

func (c *Character) GetSpriteSheetCoordinates() logical.Vec {
	return c.Sprite.Add(logical.V(c.SpriteIdx, 0))
}

func (c *Character) GetColor() color.Color {
	return c.Color
	//    return render.GetColor(c.ColorR, c.ColorG, c.ColorB)
}

func (c *Character) Describe() string {
	return c.Name
}

func (c *Character) SetBoardPosition(v logical.Vec) {
	c.BoardPosition = v
}

/*
func (ct CharacterTypes) NewCharacter(typeName string) *Character {
	c := ct[typeName]
	ch := &Character{CharacterType: c}

	spriteCount := len(ch.Sprites)
	if spriteCount > 1 {
		ch.SpriteIdx = rand.Intn(spriteCount - 1)
	}

	return ch
}

// GameObject interface BEGIN
func (c *Character) AnimationTick() {
	if c.Sprites == nil {
		return
	}
	spriteCount := len(c.Sprites)
	if spriteCount == 0 {
		return
	}
	c.SpriteIdx++
	if c.SpriteIdx == spriteCount {
		c.SpriteIdx = 0
	}
	return
}

func (c *Character) RemoveMe() bool {
	return false
}

func (c *Character) IsEmpty() bool {
	return false
}

func (c *Character) GetSpriteSheetCoordinates() logical.Vec {
	return logical.V(c.Sprites[c.SpriteIdx][0], c.Sprites[c.SpriteIdx][1])
}

func (c *Character) GetColor() color.Color {
	return render.GetColor(c.ColorR, c.ColorG, c.ColorB)
}

func (c *Character) Describe() string {
	return c.Name
}

func (c *Character) SetBoardPosition(v logical.Vec) {
	c.BoardPosition = v
}

// GameObject interface END

func (c *Character) GetColorMask() color.Color {
	if c.ColorR == 0 && c.ColorG == 0 && c.ColorB == 0 {
		return pixel.RGB(1, 1, 1)
	}
	return pixel.RGB(float64(c.ColorR)/255, float64(c.ColorG)/255, float64(c.ColorB)/255)
} */
