package grid

import (
	"image/color"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

type GameObject interface {
	AnimationTick()
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

func (s *GameObjectStack) TopObject() GameObject {
	return (*s)[0]
}

func (s *GameObjectStack) SetBoardPosition(v logical.Vec) {
	(*s)[0].SetBoardPosition(v)
}

func (s *GameObjectStack) AnimationTick() {
	(*s)[0].AnimationTick()
	if (*s)[0].RemoveMe() {
		(*s) = (*s)[1:]
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

func NewGameObjectStack() *GameObjectStack {
	os := make(GameObjectStack, 1)
	os[0] = EMPTY_OBJECT
	return &os
}

/* Empty game object (bottom level tile) */

const BLANK_SPRITE_X = 8
const BLANK_SPRITE_Y = 26

type EmptyObject struct {
	SpriteCoordinates logical.Vec
}

var EMPTY_OBJECT = EmptyObject{
	SpriteCoordinates: logical.V(BLANK_SPRITE_X, BLANK_SPRITE_Y),
}

func (e EmptyObject) AnimationTick() {}

func (e EmptyObject) RemoveMe() bool {
	return false
}

func (e EmptyObject) GetColor() color.Color {
	return render.GetColor(0, 0, 0)
}

func (e EmptyObject) Describe() string {
	return "                                  "
}

func (e EmptyObject) IsEmpty() bool {
	return true
}

func (e EmptyObject) GetSpriteSheetCoordinates() logical.Vec {
	return e.SpriteCoordinates
}

// Ignore this
func (e EmptyObject) SetBoardPosition(v logical.Vec) {}
