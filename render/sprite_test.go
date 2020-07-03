package render

import (
	"testing"

	"github.com/bobtfish/mayhem/logical"
)

func TestSpriteDrawerZeroWinZeroSprite(t *testing.T) {
	sd := NewSpriteDrawer(nil, logical.ZeroVec())
	r := sd.GetPixelRect(logical.ZeroVec())
	if r.Min.X != 0.0 {
		t.Errorf("r.Min.X != 0 is %f", r.Min.X)
	}
	if r.Min.Y != 0.0 {
		t.Errorf("r.Min.Y != 0 is %f", r.Min.Y)
	}
	if r.Max.X != SPRITE_SIZE {
		t.Errorf("r.Max.X != SPRITE_SIZE is %f", r.Max.X)
	}
	if r.Max.Y != SPRITE_SIZE {
		t.Errorf("r.Max.Y != SPRITE_SIZE is %f", r.Max.Y)
	}
}

func TestSpriteDrawerOffsetWinZeroSprite(t *testing.T) {
	sd := NewSpriteDrawer(nil, logical.V(20, 20))
	r := sd.GetPixelRect(logical.ZeroVec())
	if r.Min.X != 0.0 {
		t.Errorf("r.Min.X != 0 is %f", r.Min.X)
	}
	if r.Min.Y != 0.0 {
		t.Errorf("r.Min.Y != 0 is %f", r.Min.Y)
	}
	if r.Max.X != SPRITE_SIZE {
		t.Errorf("r.Max.X != SPRITE_SIZE is %f", r.Max.X)
	}
	if r.Max.Y != SPRITE_SIZE {
		t.Errorf("r.Max.Y != SPRITE_SIZE is %f", r.Max.Y)
	}
}

func TestSpriteDrawerZeroWinOneOneSprite(t *testing.T) {
	sd := NewSpriteDrawer(nil, logical.ZeroVec())
	r := sd.GetPixelRect(logical.V(1, 1))
	if r.Min.X != SPRITE_SIZE {
		t.Errorf("r.Min.X != SPRITE_SIZE is %f", r.Min.X)
	}
	if r.Min.Y != SPRITE_SIZE {
		t.Errorf("r.Min.Y != SPRITE_SIZE is %f", r.Min.Y)
	}
	if r.Max.X != SPRITE_SIZE*2 {
		t.Errorf("r.Max.X != SPRITE_SIZE*2 is %f", r.Max.X)
	}
	if r.Max.Y != SPRITE_SIZE*2 {
		t.Errorf("r.Max.Y != SPRITE_SIZE*2 is %f", r.Max.Y)
	}
}
