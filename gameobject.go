package main

import (
	"github.com/bobtfish/mayhem/logical"
)

type GameObject interface {
	AnimationTick() bool
	GetSpriteSheetCoordinates() logical.Vec
}

/* Object stack */

type GameObjectStack []GameObject

func (s *GameObjectStack) AnimationTick() bool {
	return (*s)[0].AnimationTick()
}

func (s *GameObjectStack) GetSpriteSheetCoordinates() logical.Vec {
	return (*s)[0].GetSpriteSheetCoordinates()
}

func (s *GameObjectStack) PlaceObject(o GameObject) {
	(*s) = append([]GameObject{o}, (*s)...)
}

func NewGameObjectStack() *GameObjectStack {
	os := make(GameObjectStack, 1)
	os[0] = EMPTY_OBJECT
	return &os
}

/* Empty game object (bottom level tile */

const BLANK_SPRITE_X = 8
const BLANK_SPRITE_Y = 26

type EmptyObject struct {
	SpriteCoordinates logical.Vec
}

var EMPTY_OBJECT = EmptyObject{
	SpriteCoordinates: logical.V(BLANK_SPRITE_X, BLANK_SPRITE_Y),
}

func (e EmptyObject) AnimationTick() bool {
	return false
}

func (e EmptyObject) GetSpriteSheetCoordinates() logical.Vec {
	return e.SpriteCoordinates
}
