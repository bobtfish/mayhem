package main

import (
	"github.com/faiface/pixel/pixelgl"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type CharacterType struct {
	Name              string `yaml:"name"`
	Combat            int    `yaml:"combat"`
	RangedCombat      int    `yaml:"ranged_combat"`
	Range             int    `yaml:"range"`
	Defence           int    `yaml:"defence"`
	Movement          int    `yaml:"movement"`
	MagicalResistance int    `yaml:"magical_resistance"`
	Manoeuvre         int    `yaml:"manoeuvre"`
	Unknown           int    `yaml:"unknown"`
	LawChaos          int    `yaml:"law_chaos"`
	Strength          int    `yaml:"strength"`
}

type Character struct {
	CharacterType
}

func (c *Character) DrawCharacter(win *pixelgl.Window) {
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
		ct[v.Name] = v
	}
	return ct
}

func (ct CharacterTypes) NewCharacter(typeName string) *Character {
	c := ct[typeName]
	return &Character{CharacterType: c}
}
