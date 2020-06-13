package main

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
)

type Character struct {
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

type CharacterTypes map[string]Character

func LoadCharacters(fn string) CharacterTypes {
	yamlFile, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	cl := make([]Character, 0)
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
