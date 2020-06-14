package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"image/color"
	"io/ioutil"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

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

type Character struct {
	CharacterType
	SpriteIdx int
}

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

func (c *Character) GetSprite(ss pixel.Picture) *pixel.Sprite {
	if len(c.Sprites) == 0 {
		return nil
	}
	spriteLocation := c.Sprites[c.SpriteIdx]
	x := spriteLocation[0]
	y := spriteLocation[1]
	//	fmt.Printf("Character %s has %d sprites, in sprite sheet sprite 0 is x %d, y %d topx %d topy %d\n", c.Name, len(c.Sprites), x*SPRITE_SIZE, y*SPRITE_SIZE, x*SPRITE_SIZE+SPRITE_SIZE, y*SPRITE_SIZE+SPRITE_SIZE)
	return pixel.NewSprite(ss, pixel.R(float64(x*SPRITE_SIZE), float64(y*SPRITE_SIZE), float64(x*SPRITE_SIZE+SPRITE_SIZE), float64(y*SPRITE_SIZE+SPRITE_SIZE)))
}

func (c *Character) GetColorMask() color.Color {
	if c.ColorR == 0 && c.ColorG == 0 && c.ColorB == 0 {
		return pixel.RGB(1, 1, 1)
	}
	return pixel.RGB(float64(c.ColorR)/255, float64(c.ColorG)/255, float64(c.ColorB)/255)
}

func (c *Character) GetText(x, y int) *text.Text {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(float64(x+2), float64(y+2+CHAR_PIXELS-16)), atlas)
	fmt.Fprintln(basicTxt, c.Name)
	fmt.Fprintf(basicTxt, "R%d C%d RC%d D%d\n", c.Range, c.Combat, c.RangedCombat, c.Defence)
	fmt.Fprintf(basicTxt, "M%d MR%d\n", c.Movement, c.MagicalResistance)
	fmt.Fprintf(basicTxt, "LC%d ST%d\n", c.LawChaos, c.Strength)
	return basicTxt
}

type CharacterTypes map[string]CharacterType

func LoadCharacterTemplates(fn string) CharacterTypes {
	yamlFile, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	cl := make([]CharacterType, 0)
	err = yaml.Unmarshal(yamlFile, &cl)
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
