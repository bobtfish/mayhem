package grid

import (
	"image/color"

	"github.com/bobtfish/mayhem/logical"
)

type GameObject interface {
	AnimationTick(bool)
	GetSpriteSheetCoordinates() logical.Vec
	GetColor() color.Color
	IsEmpty() bool
	Describe() string
	SetBoardPosition(logical.Vec)
}

/* Object stack */

type GameObjectStackable interface {
	GameObject
	RemoveMe() bool
}

type GameObjectStack []GameObjectStackable

func (s *GameObjectStack) TopObject() GameObjectStackable {
	return (*s)[0]
}

func (s *GameObjectStack) SetBoardPosition(v logical.Vec) {
	(*s)[0].SetBoardPosition(v)
}

func (s *GameObjectStack) AnimationTick(odd bool) {
	s.TopObject().AnimationTick(odd)
	if s.TopObject().RemoveMe() {
		s.RemoveTopObject()
	}
}

func (s *GameObjectStack) GetSpriteSheetCoordinates() logical.Vec {
	return (*s)[0].GetSpriteSheetCoordinates()
}

func (s *GameObjectStack) GetColor() color.Color {
	return (*s)[0].GetColor()
}

func (s *GameObjectStack) Describe() string {
	return (*s)[0].Describe()
}

func (s *GameObjectStack) IsEmpty() bool {
	return len(*s) == 1
}

func (s *GameObjectStack) PlaceObject(o GameObjectStackable) {
	(*s) = append([]GameObjectStackable{o}, (*s)...)
}

func (s *GameObjectStack) RemoveTopObject() GameObjectStackable {
	if s.IsEmpty() {
		return nil
	}
	topObject := (*s)[0]
	(*s) = (*s)[1:]
	return topObject
}

func NewGameObjectStack() *GameObjectStack {
	os := make(GameObjectStack, 1)
	os[0] = EMPTY_OBJECT
	return &os
}
