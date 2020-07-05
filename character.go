package main

import (
	"gopkg.in/yaml.v2"
	"image/color"
	"math/rand"

	"github.com/faiface/pixel"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

// Characters on the board

// Abstract character that can be created
type CharacterTypes map[string]CharacterType
type CharacterType struct {
	Name              string  `yaml:"name"`
	Combat            int     `yaml:"combat"`
	RangedCombat      int     `yaml:"ranged_combat"`
	Range             int     `yaml:"range"`
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
type Character struct {
	CharacterType
	SpriteIdx int
}

func LoadCharacterTemplates() CharacterTypes {
	cl := make([]CharacterType, 0)
	err := yaml.Unmarshal([]byte(character_yaml), &cl)
	if err != nil {
		panic(err)
	}
	ct := make(CharacterTypes, 0)
	for _, v := range cl {
		if v.Sprites == nil {
			v.Sprites = make([][]int, 0)
		}
		if v.DeadSprite == nil {
			v.DeadSprite = make([]int, 0)
		}
		ct[v.Name] = v
	}
	return ct
}

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

// GameObject interface END

func (c *Character) GetColorMask() color.Color {
	if c.ColorR == 0 && c.ColorG == 0 && c.ColorB == 0 {
		return pixel.RGB(1, 1, 1)
	}
	return pixel.RGB(float64(c.ColorR)/255, float64(c.ColorG)/255, float64(c.ColorB)/255)
}
